package material

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Repository interface {
	Lister
	ListerBySIADS
	ListerByCATMAT
	ListerByECampus
	Getter
	Creater
	Patcher
	Deleter
}

type (
	Lister interface {
		List(offset, limit int) (Entities, error)
	}

	ListerBySIADS interface {
		ListBySIADS(siads string) (Entities, error)
	}

	ListerByCATMAT interface {
		ListByCATMAT(catmat string) (Entities, error)
	}

	ListerByECampus interface {
		ListByECampus(ecampus string) (Entities, error)
	}

	Getter interface {
		Get(uuid uuid.UUID) (Entity, error)
	}

	Creater interface {
		Create(Entity) error
	}

	Patcher interface {
		Patch(uuid.UUID, PartialEntity) error
	}

	Deleter interface {
		Delete(uuid uuid.UUID) error
	}
)

type (
	Entities = ListResult
	Entity   = Result

	PartialEntity struct {
		Name        opt.Opt[string]
		SIADS       opt.Opt[string]
		CATMAT      opt.Opt[string]
		ECampus     opt.Opt[string]
		Description opt.Opt[string]
		Unit        opt.Opt[string]
		MinQuantity opt.Opt[float64]
		Updated     time.Time
	}
)
