package session

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Repository interface {
	Get(uuid.UUID) (Entity, error)

	Create(Entity) error
	
	Update(uuid.UUID, time.Time) error

	Delete(uuid.UUID) error
}

type (
	Entity struct {
		UUID    uuid.UUID
		User    uuid.UUID
		Expires time.Time
	}
)
