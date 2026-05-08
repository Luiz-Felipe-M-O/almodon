package web

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

var ErrNotExists = errors.New("template does not exists")

func Base() *template.Template {
	return template.Must(base.Clone())
}

func Text(name string) (string, error) {
	if text, in := texts[name]; in {
		return text, nil
	}

	return "", fmt.Errorf("name %s: %w", name, ErrNotExists)
}

func MustText(name string) string {
	text, err := Text(name)
	if err != nil {
		panic(err)
	}
	return text
}

func ParseChain(tmpl *template.Template, texts ...string) (*template.Template, error) {
	for _, text := range texts {
		var err error

		tmpl, err = tmpl.Parse(text)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

var exts = [...]string{".html", ".gohtml"}

const pages_dir = "dist/pages"

//go:embed dist/pages
var pages embed.FS

var (
	base  *template.Template
	texts map[string]string
)

func init() {
	texts = make(map[string]string)

	err := fs.WalkDir(pages, pages_dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() != "." && strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if !slices.Contains(exts[:], ext) {
			return nil
		}

		text, err := pages.ReadFile(path)
		if err != nil {
			return err
		}
		tmpl := string(text)

		name := strings.TrimSuffix(path[len(pages_dir)+1:], ext)

		if deindex, ok := strings.CutSuffix(name, "/index"); ok {
			texts[deindex] = tmpl
		}
		texts[name] = tmpl

		return nil
	})
	if err != nil {
		panic(err)
	}

	base = template.Must(template.New("almodon").Parse(MustText("index")))
}
