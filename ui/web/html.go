package web

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"time"
)

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

type Page struct {
	build  time.Time
	reader bytes.Reader
}

func MakePage(tmpl *template.Template, data any) (*Page, error) {
	var b bytes.Buffer
	if err := tmpl.Execute(&b, data); err != nil {
		return nil, err
	}

	return &Page{
		build:  time.Now(),
		reader: *bytes.NewReader(b.Bytes()),
	}, nil
}

func (p *Page) Build() time.Time {
	return p.build
}

func (p *Page) ReadSeeker() io.ReadSeeker {
	return &p.reader
}

func (p *Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, ".html", p.build, &p.reader)
}
