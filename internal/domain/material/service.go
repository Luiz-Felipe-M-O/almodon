package material

import (
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	List(act auth.Actor, req ListRequest) (ListResponse, error)

	Get(act auth.Actor, req GetRequest) (Response, error)
	GetBySIADS(act auth.Actor, req GetBySIADSRequest) (Response, error)
	GetByCATMAT(act auth.Actor, req GetByCATMATRequest) (Response, error)
	GetByECAMPUS(act auth.Actor, req GetByECAMPUSRequest) (Response, error)

	Create(act auth.Actor, req CreateRequest) (uuid.UUID, error)

	Patch(act auth.Actor, req PatchRequest) error
	UpdateMinQuantity(act auth.Actor, req UpdateMinQuantityRequest) error

	Delete(act auth.Actor, req DeleteRequest) error

	Search(act auth.Actor, req SearchRequest) (ListResponse, error)
}
