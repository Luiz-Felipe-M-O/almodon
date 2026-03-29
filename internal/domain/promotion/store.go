package promotion

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Store interface {
	Get(context.Context, uuid.UUID) (Record, error)
	GetByUser(context.Context, uuid.UUID) (Record, error)

	Create(context.Context, CreateRecord) error

	Update(context.Context, uuid.UUID, UpdateRecord) error

	Delete(context.Context, uuid.UUID) error

	RunTx(context.Context, func(Store) error) error
}

type (
	Record struct {
		UUID    uuid.UUID
		User    uuid.UUID
		Expires time.Time
	}

	CreateRecord struct {
		UUID    uuid.UUID
		User    uuid.UUID
		Expires time.Time
	}

	UpdateRecord struct {
		Expires time.Time
	}
)
