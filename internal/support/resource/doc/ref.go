package doc

import (
	"net/http"

	"github.com/alan-b-lima/almodon/ui/web"
)

type Ref struct {
	Title string
	Docs  []*Doc

	mux *http.ServeMux
}

func NewRef(glob *web.Glob, title string, docs []*Doc) (*Ref, error) {
	ref := Ref{
		Title: title,
		Docs:  docs,
	}

	tmpl, err := glob.Parse("index", "doc/ref")
	if err != nil {
		return nil, err
	}

	page, err := web.MakePage(tmpl, ref)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	for _, doc := range docs {
		mux.Handle("/docs/"+doc.Path, doc)
	}
	mux.Handle("/docs/", page)

	ref.mux = mux
	return &ref, nil
}

func (d *Ref) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.mux.ServeHTTP(w, r)
}
