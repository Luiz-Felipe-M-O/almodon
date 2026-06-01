package almodon

import (
	"net/http"

	"github.com/alan-b-lima/almodon/internal/client"
	sessions "github.com/alan-b-lima/almodon/internal/domain/session/resource"
)

type Almodon struct {
	http.Handler
	api *API
}

func New() (*Almodon, error) {
	var a Almodon

	client, err := client.New()
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

	mux := http.NewServeMux()
	mux.Handle("/api/", api)
	mux.Handle("/", client)

	a.Handler = sessions.Wrap(mux)
	return &a, nil
}

func (a *Almodon) Close() error {
	return a.api.Close()
}
