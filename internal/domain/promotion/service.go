package promotion

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/opt"
)

type Service interface {
	Get(context.Context, uuid.UUID) (Result, error)
	GetByUser(context.Context, uuid.UUID) (Result, error)

	Create(context.Context, Create) (CreateResult, error)

	Update(context.Context, uuid.UUID, Update) error

	Delete(context.Context, uuid.UUID) error
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
		UUID    uuid.UUID `json:"uuid"`
		User    uuid.UUID `json:"user"`
		Expires time.Time `json:"expires"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}
)
