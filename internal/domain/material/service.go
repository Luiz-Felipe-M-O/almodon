package material

import (
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	List(req ListParams) (Entities, error)
	ListBySIADS(siads string, req ListParams) (Entities, error)
	ListByCATMAT(catmat string, req ListParams) (Entities, error)
	ListByECampus(ecampus string, req ListParams) (Entities, error)

	Get(uuid uuid.UUID) (Entity, error)

	Create(req Create) (uuid.UUID, error)

	Patch(uuid uuid.UUID, req Patch) error

	Delete(uuid uuid.UUID) error
}
