package material

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Repository interface {
	Lister
	Getter
	ListerBySIADS
	ListerByCATMAT
	ListerByECAMPUS
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

	ListerBySIADS interface {
		ListBySIADS(siads string) (Entities, error)
	}

	ListerByCATMAT interface {
		ListByCATMAT(catmat string) (Entities, error)
	}

	ListerByECAMPUS interface {
		ListByECAMPUS(ecampus string) (Entities, error)
	}

	Creater interface {
		Create(Entity) error
	}

	Patcher interface {
		Patch(uuid uuid.UUID, PartialEntity) error
	}

	Deleter interface {
		Delete(uuid uuid.UUID) error
	}
)

type (
	Entity struct {
		UUID        uuid.UUID
		Name        string
		SIADS       string
		CATMAT      string
		ECAMPUS     string
		Description string
		Unit        string
		MinQuantity float64
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	PartialEntity struct {
		Name        opt.Opt[string]
		SIADS       opt.Opt[string]
		CATMAT      opt.Opt[string]
		ECAMPUS     opt.Opt[string]
		Description opt.Opt[string]
		Unit        opt.Opt[string]
		MinQuantity opt.Opt[float64]
		CreatedAt   opt.Opt[time.Time]
		UpdatedAt   opt.Opt[time.Time]
	}

	Entities struct {
		Offset       int
		Length       int
		Records      []Entity
		TotalRecords int
	}
)
