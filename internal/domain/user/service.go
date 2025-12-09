package user

import (
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	Lister
	Getter
	Creater
	Updater
	Deleter
}

type (
	Lister interface {
		List(req ListParams) (Entities, error)
	}

	Getter interface {
		Get(uuid uuid.UUID) (Entity, error)
		GetBySIAPE(siape string) (Entity, error)
	}

	Creater interface {
		Create(req Create) (uuid.UUID, error)
	}

	Updater interface {
		Patch(uuid uuid.UUID, req Patch) error
		UpdatePassword(uuid uuid.UUID, req UpdatePassword) error
		UpdateRole(uuid uuid.UUID, req UpdateRole) error
	}

	Deleter interface {
		Delete(uuid uuid.UUID) error
	}
)
