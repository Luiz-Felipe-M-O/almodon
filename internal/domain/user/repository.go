package user

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Repository interface {
	List(offset, limit int) (Entities, error)

	Get(uuid.UUID) (Entity, error)
	GetBySIAPE(string) (Entity, error)
	
	Create(Entity) error
	
	Patch(uuid.UUID, PartialEntity) error
	
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
