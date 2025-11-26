package user

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	List(req ListParams) (Entities, error)

	Get(uuid uuid.UUID) (Entity, error)
	GetBySIAPE(siape string) (Entity, error)

	Create(req Create) (uuid.UUID, error)

	Patch(uuid uuid.UUID, req Patch) error
	UpdatePassword(uuid uuid.UUID, req UpdatePassword) error
	UpdateRole(uuid uuid.UUID, req UpdateRole) error

	Delete(uuid uuid.UUID) error

	Authenticate(siape string, password string) (AuthEntity, error)
	Actor(session uuid.UUID) (auth.Actor, error)
}
