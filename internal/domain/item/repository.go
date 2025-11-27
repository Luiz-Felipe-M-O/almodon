package item

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Repository interface {
	Lister
	Getter
	ListerByBatch
	ListerByMaterial
	Creater
	Patcher
	Deleter
}

type (
	Lister interface {
		List(offset, limit int) (Entities, error)
	}

	Getter interface {
		Get(uuid uuid.UUID) (Entity, error)
	}

	ListerByBatch interface {
		ListByBatch(batch uuid.UUID) (Entities, error)
	}

	ListerByMaterial interface {
		ListByMaterial(material uuid.UUID) (Entities, error)
	}

	Creater interface {
		Create(Entity) error
	}

	Patcher interface {
		Patch(uuid.UUID, PartialEntity) error
	}

	Deleter interface {
		Delete(uuid.UUID) error
	}
)

type (
	Entity struct {
		UUID       uuid.UUID
		Batch      uuid.UUID
		Material   uuid.UUID
		Quantity   float64
		Expiration time.Time
		Created    time.Time
		Updated    time.Time
	}

	PartialEntity struct {
		Batch      opt.Opt[uuid.UUID]
		Material   opt.Opt[uuid.UUID]
		Quantity   opt.Opt[float64]
		Expiration opt.Opt[time.Time]
		Updated    time.Time
	}

	Entities struct {
		Offset       int
		Length       int
		Records      []Entity
		TotalRecords int
	}
)
