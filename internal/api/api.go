package api

import (
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain"
	"github.com/alan-b-lima/almodon/internal/support/session"
	"github.com/alan-b-lima/almodon/pkg/closer"
)

type Almodon struct {
	http.Handler
	bundle closer.Bundle
}

func New() (*Almodon, error) {
	domain, err := domain.New()
	if err != nil {
		return nil, err
	}

	handlers := map[string]http.Handler{
		"auth":       domain.Resources.Auth,
		"items":      domain.Resources.Items,
		"materials":  domain.Resources.Materials,
		"promotions": domain.Resources.Promotions,
		"users":      domain.Resources.Users,
	}

	mux := http.NewServeMux()
	for name, handler := range handlers {
		mux.Handle("/api/v1/"+name+"/", http.StripPrefix("/api/v1", handler))
	}

	return &Almodon{
		Handler: session.Wrap(mux),
		bundle:  domain.Bundle,
	}, nil
}

func (a *Almodon) Close() error {
	return a.bundle.Close()
}
