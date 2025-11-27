package item

import "github.com/alan-b-lima/almodon/pkg/uuid"

type Service interface {
	List(req ListParams) (Entities, error)
	ListByBatch(uuid uuid.UUID) (Entities, error)
	ListByMaterial(uuid uuid.UUID) (Entities, error)

	Get(uuid uuid.UUID) (Entity, error)

	Create(req Create) (uuid.UUID, error)

	Patch(uuid uuid.UUID, req Patch) error
	UpdateQuantity(uuid uuid.UUID,req UpdateQuantity) error

	Delete(uuid uuid.UUID) error
}
