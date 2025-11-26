package user

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Repository interface {
	Lister
	Getter
	GetterBySIAPE
	Creater
	Patcher
	Deleter
}

type (
	Lister interface {
		List(offset, limit int) (Entities, error)
	}

	Getter interface {
		Get(uuid.UUID) (Entity, error)
	}

	GetterBySIAPE interface {
		GetBySIAPE(string) (Entity, error)
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
	Entities struct {
		Offset       int
		Length       int
		Records      []Entity
		TotalRecords int
	}

	Entity struct {
		UUID     uuid.UUID
		SIAPE    string
		Name     string
		Email    string
		Password [60]byte
		Role     auth.Role
		Created  time.Time
		Updated  time.Time
	}

	PartialEntity struct {
		SIAPE    opt.Opt[string]
		Name     opt.Opt[string]
		Email    opt.Opt[string]
		Password opt.Opt[[60]byte]
		Role     opt.Opt[auth.Role]
		Updated  time.Time
	}

	AuthEntity struct {
		UUID    uuid.UUID
		User    uuid.UUID
		Expires time.Time
	}
)
