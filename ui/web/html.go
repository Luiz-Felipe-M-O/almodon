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

type Content struct {
	Ext   string
	Build time.Time
	bytes []byte
}

func NewContent(ext string, b []byte) *Content {
	return &Content{Ext: ext, Build: time.Now(), bytes: b}
}

func (p *Content) ReadSeeker() io.ReadSeeker {
	return bytes.NewReader(p.bytes)
}

func (p *Content) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, p.Ext, p.Build, p.ReadSeeker())
}

func MakePage(tmpl *template.Template, data any) (*Content, error) {
	var b bytes.Buffer
	if err := tmpl.Execute(&b, data); err != nil {
		return nil, err
	}

	return NewContent(".html", b.Bytes()), nil
}
