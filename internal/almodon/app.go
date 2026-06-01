package almodon

import (
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain"
	sessions "github.com/alan-b-lima/almodon/internal/domain/session/resource"
	"github.com/alan-b-lima/almodon/ui/web"
)

type Almodon struct {
	mux http.ServeMux
	api *API
}

func New() (*Almodon, error) {
	var a Almodon

	glob, err := web.NewGlobDyn()
	if err != nil {
		return nil, err
	}

	toolkit := web.NewToolkitDyn(glob)

	docs, err := domain.Reference(glob)
	if err != nil {
		return nil, err
	}

	// mount the API last not to handle conditinal
	// closing, see [domain.New] for an example
	api, err := NewAPI()
	if err != nil {
		return nil, err
	}
	a.api = api

	a.mux.Handle("/toolkit/", toolkit)
	a.mux.Handle("/docs/", docs)
	a.mux.Handle("/api/", sessions.Wrap(api))

	return &a, nil
}

func (a *Almodon) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *Almodon) Close() error {
	return a.api.Close()
}
