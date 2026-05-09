package almodon

import (
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain"
	"github.com/alan-b-lima/almodon/internal/support/session"
	"github.com/alan-b-lima/almodon/ui/web"
)

type Almodon struct {
	mux http.ServeMux
	api *API
}

func New() (*Almodon, error) {
	var a Almodon

	toolkit := web.ToolkitDyn()

	docs, err := domain.Reference()
	if err != nil {
		return nil, err
	}

	// mount the API last not to handle closing
	// see [domain.New] for an example
	api, err := NewAPI()
	if err != nil {
		return nil, err
	}
	a.api = api

	a.mux.Handle("/toolkit/", toolkit)
	a.mux.Handle("/docs/", docs)
	a.mux.Handle("/api/", session.Wrap(api))

	return &a, nil
}

func (a *Almodon) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *Almodon) Close() error {
	return a.api.bundle.Close()
}
