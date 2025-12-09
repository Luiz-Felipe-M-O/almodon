package session

import (
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	Getter
	Creater
	Updater
	Deleter
}

type (
	Getter interface {
		Get(uuid.UUID) (Entity, error)
	}

	Creater interface {
		CreateAndGet(Create) (Entity, error)
	}

	Updater interface {
		Update(uuid.UUID, Update) error
	}

	Deleter interface {
		Delete(uuid.UUID) error
	}
)
