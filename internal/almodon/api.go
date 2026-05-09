package almodon

import (
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain"
	"github.com/alan-b-lima/almodon/pkg/closer"
)

type API struct {
	http.ServeMux
	bundle closer.Bundle
}

func NewAPI() (*API, error) {
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

	api := &API{bundle: domain.Bundle}
	for name, handler := range handlers {
		api.Handle("/api/v1/"+name+"/", http.StripPrefix("/api/v1", handler))
	}

	return api, nil
}
