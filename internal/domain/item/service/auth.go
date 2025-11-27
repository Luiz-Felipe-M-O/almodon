package itemserve

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/item"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Gate struct {
	item.Service
	actor auth.Actor
}

func New(service item.Service) item.Service {
	return &Gate{Service: service}
}

var (
	permUser  = auth.Permit(auth.User)
	permAdmin = auth.Permit(auth.Admin)
)

func (s *Gate) Allow(act auth.Actor) item.Service {
	return &Gate{
		Service: s.Service,
		actor:   act,
	}
}

func (s *Gate) List(req item.ListParams) (item.Entities, error) {
	if err := service.Authorize(permUser, s.actor); err != nil {
		return item.Entities{}, err
	}

	return s.Service.List(req)
}

func (s *Gate) ListByBatch(uuid uuid.UUID) (item.Entities, error) {
	if err := service.Authorize(permUser, s.actor); err != nil {
		return item.Entities{}, err
	}

	return s.Service.ListByBatch(uuid)
}

func (s *Gate) ListByMaterial(uuid uuid.UUID) (item.Entities, error) {
	if err := service.Authorize(permUser, s.actor); err != nil {
		return item.Entities{}, err
	}

	return s.Service.ListByMaterial(uuid)
}

func (s *Gate) Get(uuid uuid.UUID) (item.Entity, error) {
	if err := service.Authorize(permUser, s.actor); err != nil {
		return item.Entity{}, err
	}

	return s.Service.Get(uuid)
}

func (s *Gate) Create(req item.Create) (uuid.UUID, error) {
	if err := service.Authorize(permAdmin, s.actor); err != nil {
		return uuid.UUID{}, err
	}

	return s.Service.Create(req)
}

func (s *Gate) Patch(uuid uuid.UUID, req item.Patch) error {
	if err := service.Authorize(permAdmin, s.actor); err != nil {
		return err
	}

	return s.Service.Patch(uuid, req)
}

func (s *Gate) UpdateQuantity(uuid uuid.UUID, req item.UpdateQuantity) error {
	if err := service.Authorize(permAdmin, s.actor); err != nil {
		return err
	}

	return s.Service.UpdateQuantity(uuid, req)
}

func (s *Gate) Delete(uuid uuid.UUID) error {
	if err := service.Authorize(permAdmin, s.actor); err != nil {
		return err
	}

	return s.Service.Delete(uuid)
}
