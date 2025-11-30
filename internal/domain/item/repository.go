package item

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Repository interface {
	Lister
	Getter
	ListerByMaterial
	ListerBySupplier
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

	ListerByMaterial interface {
		ListByMaterial(material uuid.UUID) (Entities, error)
	}

	ListerBySupplier interface {
		ListBySupplier(supplier uuid.UUID) (Entities, error)
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
		Material   uuid.UUID
		Supplier   uuid.UUID
		Quantity   float64
		UnitCost   float64
		Arrival    time.Time
		Expiration time.Time
		Invoice    string
		Lot        string
		Notes      string
		Created    time.Time
		Updated    time.Time
	}

	PartialEntity struct {
		Material   opt.Opt[uuid.UUID]
		Supplier   opt.Opt[uuid.UUID]
		Quantity   opt.Opt[float64]
		UnitCost   opt.Opt[float64]
		Arrival    opt.Opt[time.Time]
		Expiration opt.Opt[time.Time]
		Invoice    opt.Opt[string]
		Lot        opt.Opt[string]
		Notes      opt.Opt[string]
		Updated    time.Time
	}

	Entities struct {
		Offset       int
		Length       int
		Records      []Entity
		TotalRecords int
	}
)
