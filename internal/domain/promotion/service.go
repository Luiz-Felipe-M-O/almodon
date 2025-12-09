package promotion

import "github.com/alan-b-lima/almodon/pkg/uuid"

type Service interface {
	Lister
	Getter
	Creater
	Updater
	Deleter
}

type (
	Lister interface {
		List(ListParams) (Entities, error)
	}

	Getter interface {
		Get(uuid.UUID) (Entity, error)
		GetByUser(uuid.UUID) (Entity, error)
	}

	Creater interface {
		Create(Create) (uuid.UUID, error)
	}

	Updater interface {
		Update(uuid.UUID, Update) error
	}

	Deleter interface {
		Delete(uuid.UUID) error
	}
)
