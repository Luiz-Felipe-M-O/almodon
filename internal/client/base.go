package client

import (
	"fmt"
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain"
	"github.com/alan-b-lima/almodon/ui/web"
)

type Client struct {
	http.ServeMux
}

func New() (*Client, error) {
	glob, err := web.NewGlobDyn()
	if err != nil {
		return nil, fmt.Errorf("glob: %w", err)
	}

	toolkit := web.NewToolkitDyn(glob)

	var client Client
	client.Handle("/toolkit/", toolkit)

	docs, err := domain.Reference(glob)
	if err != nil {
		return nil, fmt.Errorf("docs: %w", err)
	}

	about, err := About(glob)
	if err != nil {
		return nil, fmt.Errorf("about: %w", err)
	}

	client.Handle("/docs/", docs)
	client.Handle("/about/{$}", about)

	return &client, nil
}
