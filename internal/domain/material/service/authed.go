package materialserve

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/material"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type AuthService struct {
	material.Service
}

func New(service material.Service) material.Service {
	return &AuthService{
		Service: service,
	}
}

var permChief = auth.Permit(auth.Chief)

func (s *AuthService) List(act auth.Actor, req material.ListRequest) (material.ListResponse, error) {
	return s.Service.List(act, req)
}

func (s *AuthService) Get(act auth.Actor, req material.GetRequest) (material.Response, error) {
	return s.Service.Get(act, req)
}

func (s *AuthService) ListBySIADS(act auth.Actor, req material.ListBySIADSRequest) (material.ListResponse, error) {
	return s.Service.ListBySIADS(act, req)
}

func (s *AuthService) ListByCATMAT(act auth.Actor, req material.ListByCATMATRequest) (material.ListResponse, error) {
	return s.Service.ListByCATMAT(act, req)
}

func (s *AuthService) ListByECampus(act auth.Actor, req material.ListByECampusRequest) (material.ListResponse, error) {
	return s.Service.ListByECampus(act, req)
}

func (s *AuthService) Create(act auth.Actor, req material.CreateRequest) (uuid.UUID, error) {
	if err := service.Authorize(permChief, act); err != nil {
		return uuid.UUID{}, err
	}

	return s.Service.Create(act, req)
}

func (s *AuthService) Patch(act auth.Actor, req material.PatchRequest) error {
	if err := service.Authorize(permChief, act); err != nil {
		return err
	}

	return s.Service.Patch(act, req)
}

func (s *AuthService) UpdateMinQuantity(act auth.Actor, req material.UpdateMinQuantityRequest) error {
	if err := service.Authorize(permChief, act); err != nil {
		return err
	}

	return s.Service.UpdateMinQuantity(act, req)
}

func (s *AuthService) Delete(act auth.Actor, req material.DeleteRequest) error {
	if err := service.Authorize(permChief, act); err != nil {
		return err
	}

	return s.Service.Delete(act, req)
}

func (s *AuthService) Search(act auth.Actor, req material.SearchRequest) (material.ListResponse, error) {
	return s.Service.Search(act, req)
}
