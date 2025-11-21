package session

import (
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	Get(uuid.UUID) (Entity, error)

	CreateAndGet(Create) (Entity, error)

	Update(uuid.UUID, Update) error

	Delete(uuid.UUID) error
}
