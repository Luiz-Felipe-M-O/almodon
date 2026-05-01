package doc

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/alan-b-lima/almodon/ui/web"
)

type Ref struct {
	Title string
	Docs  []*Doc

	mux *http.ServeMux
	buf bytes.Reader
}

func NewRef(title string, docs []*Doc) (*Ref, error) {
	ref := Ref{
		Title: title,
		Docs:  docs,
	}

	var buf bytes.Buffer
	if err := ref_tmpl.Execute(&buf, ref); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	for _, doc := range docs {
		mux.Handle("/"+doc.Path, doc)
	}
	mux.HandleFunc("/", ref.home)

	ref.mux = mux
	ref.buf = *bytes.NewReader(buf.Bytes())
	return &ref, nil
}

func (d *Ref) home(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, ".html", time.Time{}, &d.buf)
}

func (d *Ref) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.mux.ServeHTTP(w, r)
}

var ref_tmpl = template.Must(web.ParseChain(
	web.Base().Funcs(template.FuncMap{"lower": strings.ToLower}),
	`{{ define "head" }}<link rel="stylesheet" href="/toolkit/style/doc.css">{{ end }}`,
	web.MustText("doc/ref"),
))
