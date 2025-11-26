package item

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	List(act auth.Actor, req ListRequest) (ListResponse, error)

	Get(act auth.Actor, req GetRequest) (Response, error)
	GetByBatch(act auth.Actor, req GetByBatchRequest) (ListResponse, error)
	GetByMaterial(act auth.Actor, req GetByMaterialRequest) (ListResponse, error)

	Create(act auth.Actor, req CreateRequest) (uuid.UUID, error)

	Patch(act auth.Actor, req PatchRequest) error
	UpdateQuantity(act auth.Actor, req UpdateQuantityRequest) error

	Delete(act auth.Actor, req DeleteRequest) error
}
