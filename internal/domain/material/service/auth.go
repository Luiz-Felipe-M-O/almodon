package materialserve

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/material"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Gate struct {
	material.Service
	actor auth.Actor
}

func New(service material.Service) material.Service {
	return &Gate{Service: service}
}

var (
	permUser  = auth.Permit(auth.User)
	permAdmin = auth.Permit(auth.Admin)
)

func (s *Gate) Allow(act auth.Actor) material.Service {
	return &Gate{
		Service: s.Service,
		actor:   act,
	}
}

func (s *Gate) List(req material.ListParams) (material.Entities, error) {
	if err := service.Authorize(permUser, s.actor); err != nil {
		return material.Entities{}, err
	}

	return s.Service.List(req)
}

func (s *Gate) ListByCATMAT(catmat string, req material.ListParams) (material.Entities, error) {
	if err := service.Authorize(permUser, s.actor); err != nil {
		return material.Entities{}, err
	}

	return s.Service.ListByCATMAT(catmat, req)
}

func (s *Gate) ListByECampus(ecampus string, req material.ListParams) (material.Entities, error) {
	if err := service.Authorize(permUser, s.actor); err != nil {
		return material.Entities{}, err
	}

	return s.Service.ListByECampus(ecampus, req)
}

func (s *Gate) ListBySIADS(siads string, req material.ListParams) (material.Entities, error) {
	if err := service.Authorize(permUser, s.actor); err != nil {
		return material.Entities{}, err
	}

	return s.Service.ListBySIADS(siads, req)
}

func (s *Gate) Get(uuid uuid.UUID) (material.Entity, error) {
	if err := service.Authorize(permUser, s.actor); err != nil {
		return material.Entity{}, err
	}

	return s.Service.Get(uuid)
}

func (s *Gate) Create(req material.Create) (uuid.UUID, error) {
	if err := service.Authorize(permAdmin, s.actor); err != nil {
		return uuid.UUID{}, err
	}

	return s.Service.Create(req)
}

func (s *Gate) Patch(uuid uuid.UUID, req material.Patch) error {
	if err := service.Authorize(permAdmin, s.actor); err != nil {
		return err
	}

	return s.Service.Patch(uuid, req)
}

func (s *Gate) Delete(uuid uuid.UUID) error {
	if err := service.Authorize(permAdmin, s.actor); err != nil {
		return err
	}

	return s.Service.Delete(uuid)
}
