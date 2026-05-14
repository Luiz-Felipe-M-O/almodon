package session

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/opt"
)

type Service interface {
	Get(context.Context, Token) (Result, error)

	Create(context.Context, Create) (Result, error)

	Update(context.Context, Token, Update) error

	Delete(context.Context, Token) error
}

type (
	Create struct {
		User   uuid.UUID              `json:"user"`
		MaxAge opt.Opt[time.Duration] `json:"max_age"`
	}

	Update struct {
		MaxAge opt.Opt[time.Duration] `json:"max_age"`
	}
)

type (
	Result struct {
		Token   Token     `json:"-"`
		User    uuid.UUID `json:"user"`
		Renewed int       `json:"renewed"`
		Expires time.Time `json:"expires"`
		Created time.Time `json:"created"`
	}
)
