package auth

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	Login(ctx context.Context, siape string, password string) (Result, error)
	Logout(ctx context.Context, session uuid.UUID) error

	Identifier
}

type Identifier interface {
	Actor(ctx context.Context, session uuid.UUID) (Actor, error)
}

type (
	Request struct {
		SIAPE    string `json:"siape"`
		Password string `json:"password"`
	}
)

type (
	Result struct {
		UUID    uuid.UUID `json:"uuid"`
		User    uuid.UUID `json:"user"`
		Expires time.Time `json:"expires"`
	}
)
