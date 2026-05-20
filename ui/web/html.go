package web

import (
	"bytes"
	"html/template"
	"net/http"
	"time"
)

func Page(tmpl *template.Template, data any) (http.Handler, error) {
	var b bytes.Buffer
	if err := tmpl.Execute(&b, data); err != nil {
		return nil, err
	}

	return &rendered{
		build:  time.Now(),
		reader: *bytes.NewReader(b.Bytes()),
	}, nil
}

type rendered struct {
	build  time.Time
	reader bytes.Reader
}

func (p *rendered) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, ".html", p.build, &p.reader)
}
