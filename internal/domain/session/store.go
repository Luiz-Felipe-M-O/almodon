package session

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Store interface {
	List(context.Context) ([]Record, error)

	Get(context.Context, uuid.UUID) (Record, error)
	GetByUser(context.Context, uuid.UUID) (Record, error)

	Create(context.Context, CreateRecord) error

	Update(context.Context, uuid.UUID, UpdateRecord) error

	Delete(context.Context, uuid.UUID) error
	DeleteExpired(context.Context, time.Time) error

	RunTx(context.Context, func(Store) error) error
}

type (
	Record struct {
		UUID    uuid.UUID
		User    uuid.UUID
		Renewed int
		Expires time.Time
		Created time.Time
	}
)

type (
	CreateRecord struct {
		UUID    uuid.UUID
		User    uuid.UUID
		Renewed int
		Expires time.Time
		Created time.Time
	}

	UpdateRecord struct {
		Renewed int
		Expires time.Time
	}
)
