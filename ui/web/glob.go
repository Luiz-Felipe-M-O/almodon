package web

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

var ErrNotExists = errors.New("template does not exists")

type Glob struct {
	texts reader
	fnmap template.FuncMap
}

func NewGlobFS(fsys fs.FS, root string) (*Glob, error) {
	if fsys == nil {
		r, err := os.OpenRoot(root)
		if err != nil {
			return nil, err
		}

		fsys = r.FS()
		root = "."
	}

	r, err := new_reader(fsys, root)
	if err != nil {
		return nil, fmt.Errorf("parsing FS: %w", err)
	}

	return &Glob{texts: r}, nil
}

func NewGlob() (*Glob, error) {
	return NewGlobFS(pages, "dist/pages")
}

func NewGlobDyn() (*Glob, error) {
	return NewGlobFS(nil, "./ui/web/dist/pages")
}

func (g *Glob) Text(name string) string {
	text, err := g.texts.Read(name)
	if err != nil {
		return ""
	}

	return text
}

func (g *Glob) Template(name string) (*template.Template, error) {
	text, err := g.texts.Read(name)
	if err != nil {
		return nil, err
	}

	return template.New(name).Funcs(g.fnmap).Parse(text)
}

func (g *Glob) Clone() *Glob {
	return &Glob{
		texts: g.texts,
		fnmap: maps.Clone(g.fnmap),
	}
}

func (g *Glob) Parse(base string, names ...string) (*template.Template, error) {
	text, err := g.texts.Read(base)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(base).Funcs(g.fnmap).Parse(text)
	if err != nil {
		return nil, err
	}

	for _, name := range names {
		text, err := g.texts.Read(name)
		if err != nil {
			return nil, err
		}

		tmpl, err = tmpl.Parse(text)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

func (g *Glob) Funcs(fnmap template.FuncMap) *Glob {
	if g.fnmap == nil {
		g.fnmap = make(template.FuncMap, len(fnmap))
	}

	for name, fn := range fnmap {
		g.Func(name, fn)
	}

	return g
}

func (g *Glob) Func(name string, fn any) *Glob {
	if g.fnmap == nil {
		g.fnmap = make(template.FuncMap)
	}

	if fn == nil {
		delete(g.fnmap, name)
	} else {
		g.fnmap[name] = fn
	}

	return g
}

var exts = [...]string{".html", ".gohtml", ".tmpl"}

//go:embed dist/pages
var pages embed.FS

type reader struct {
	fsys  fs.FS
	index map[string]int
	files []file
	dyn   bool
}

type file struct {
	Path   string
	Text   string
	Mod    time.Time
	loaded bool

	mu sync.RWMutex
}

func new_reader(fsys fs.FS, root string) (reader, error) {
	r := reader{
		fsys:  fsys,
		index: make(map[string]int),
		dyn:   true,
	}

	if _, ok := fsys.(embed.FS); ok {
		r.dyn = false
	}

	root = filepath.Clean(root)

	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		switch {
		case err != nil:
			return err

		case strings.HasPrefix(d.Name(), "."):
			if d.Name() != root && d.IsDir() {
				return fs.SkipDir
			}

			return nil

		case d.IsDir():
			return nil
		}

		ext := filepath.Ext(path)
		if !slices.Contains(exts[:], ext) {
			return nil
		}

		return r.register(root, path)
	})

	return r, err
}

func (r *reader) Read(name string) (string, error) {
	index, in := r.index[name]
	if !in {
		return "", os.ErrNotExist
	}

	file := &r.files[index]

	file.mu.RLock()
	loaded := file.loaded
	file.mu.RUnlock()

	if !loaded || r.dyn {
		return r.load(file)
	}

	return file.Text, nil
}

func (r *reader) load(file *file) (string, error) {
	file.mu.Lock()
	defer file.mu.Unlock()

	f, err := r.fsys.Open(file.Path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return "", err
	}

	if !file.loaded || stat.ModTime().After(file.Mod) {
		var b strings.Builder
		b.Grow(int(stat.Size()))

		if _, err := io.Copy(&b, f); err != nil {
			return "", err
		}

		file.Text = b.String()
		file.Mod = stat.ModTime()
		file.loaded = true
	}

	return file.Text, nil
}

func (r *reader) register(root, path string) error {
	ext := filepath.Ext(path)

	name := filepath.ToSlash(path)
	name = filepath.Clean(name)
	name = strings.TrimSuffix(name, ext)
	name = strings.TrimPrefix(name, root)
	name = strings.TrimPrefix(name, "/")

	index := len(r.files)
	r.files = append(r.files, file{})

	file := &r.files[index]
	file.Path = path

	if i, in := r.index[name]; in {
		return fmt.Errorf("conflict: %q and %q share the same name %q", r.files[i].Path, path, name)
	}
	r.index[name] = index

	if name, ok := strings.CutSuffix(name, "index"); ok {
		name = strings.TrimSuffix(name, "/")
		if _, in := r.index[name]; !in {
			r.index[name] = index
		}
	}

	return nil
}
