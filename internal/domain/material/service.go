package material

import (
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	List(act auth.Actor, req ListRequest) (ListResponse, error)

	Get(act auth.Actor, req GetRequest) (Response, error)
	ListBySIADS(act auth.Actor, req ListBySIADSRequest) (Response, error)
	ListByCATMAT(act auth.Actor, req ListByCATMATRequest) (Response, error)
	ListByECampus(act auth.Actor, req ListByECampusRequest) (Response, error)

	Create(act auth.Actor, req CreateRequest) (uuid.UUID, error)

	Patch(act auth.Actor, req PatchRequest) error
	UpdateMinQuantity(act auth.Actor, req UpdateMinQuantityRequest) error

	Delete(act auth.Actor, req DeleteRequest) error
}
