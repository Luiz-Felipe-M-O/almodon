package web

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
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
		if d.Name() != "." && strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		text, err := pages.ReadFile(path)
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(strings.TrimPrefix(path, pages_dir+"/"), filepath.Ext(path))

		texts[name] = string(text)
		return nil
	})
	if err != nil {
		panic(err)
	}

	base = template.Must(template.New("almodon").Parse(MustText("base")))
}
