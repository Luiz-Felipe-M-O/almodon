package itemserve

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/item"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type AuthService struct {
	item.Service
}

func New(service item.Service) item.Service {
	return &AuthService{
		Service: service,
	}
}

var permChief = auth.Permit(auth.Chief)

func (s *AuthService) List(act auth.Actor, req item.ListRequest) (item.ListResponse, error) {
	return s.Service.List(act, req)
}

func (s *AuthService) Get(act auth.Actor, req item.GetRequest) (item.Response, error) {
	return s.Service.Get(act, req)
}

func (s *AuthService) GetByBatch(act auth.Actor, req item.GetByBatchRequest) (item.ListResponse, error) {
	return s.Service.GetByBatch(act, req)
}

func (s *AuthService) GetByMaterial(act auth.Actor, req item.GetByMaterialRequest) (item.ListResponse, error) {
	return s.Service.GetByMaterial(act, req)
}

func (s *AuthService) Create(act auth.Actor, req item.CreateRequest) (uuid.UUID, error) {
	if err := service.Authorize(permChief, act); err != nil {
		return uuid.UUID{}, err
	}

	return s.Service.Create(act, req)
}

func (s *AuthService) Patch(act auth.Actor, req item.PatchRequest) error {
	if err := service.Authorize(permChief, act); err != nil {
		return err
	}

	return s.Service.Patch(act, req)
}

func (s *AuthService) UpdateQuantity(act auth.Actor, req item.UpdateQuantityRequest) error {
	if err := service.Authorize(permChief, act); err != nil {
		return err
	}

	return s.Service.UpdateQuantity(act, req)
}

func (s *AuthService) Delete(act auth.Actor, req item.DeleteRequest) error {
	if err := service.Authorize(permChief, act); err != nil {
		return err
	}

	return s.Service.Delete(act, req)
}
