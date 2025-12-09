package promotion

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Repository interface {
	List(offset, limit int) (Entities, error)
	
	Get(uuid.UUID) (Entity, error)
	GetByUser(uuid.UUID) (Entity, error)
	
	Create(Entity) error
	
	Update(uuid.UUID, time.Time) error
	
	Delete(uuid.UUID) error
}

type (
	Entities struct {
		Offset       int
		Length       int
		Records      []Entity
		TotalRecords int
	}

	Entity struct {
		UUID    uuid.UUID
		User    uuid.UUID
		Expires time.Time
	}
)
