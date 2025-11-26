package userserve

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Gate struct {
	user.Service
	actor auth.Actor
}

func New(service user.Service) user.Service {
	return &Gate{Service: service}
}

var permChief = auth.Permit(auth.Chief)

func (s *Gate) Allow(act auth.Actor) user.Service {
	return &Gate{
		Service: s.Service,
		actor:   act,
	}
}

func (s *Gate) List(req user.ListParams) (user.Entities, error) {
	if err := service.Authorize(permChief, s.actor); err != nil {
		return user.Entities{}, err
	}

	return s.Service.List(req)
}

func (s *Gate) Get(uuid uuid.UUID) (user.Entity, error) {
	if s.actor.User() == uuid {
		goto Do
	}

	if err := service.Authorize(permChief, s.actor); err != nil {
		return user.Entity{}, err
	}

Do:
	return s.Service.Get(uuid)
}

func (s *Gate) GetBySIAPE(siape string) (user.Entity, error) {
	res, err := s.Service.GetBySIAPE(siape)
	if err != nil {
		return user.Entity{}, err
	}

	if s.actor.User() == res.UUID {
		goto Do
	}

	if err := service.Authorize(permChief, s.actor); err != nil {
		return user.Entity{}, err
	}

Do:
	return res, nil
}

func (s *Gate) Create(req user.Create) (uuid.UUID, error) {
	if err := service.Authorize(permChief, s.actor); err != nil {
		return uuid.UUID{}, err
	}

	return s.Service.Create(req)
}

func (s *Gate) Patch(uuid uuid.UUID, req user.Patch) error {
	if s.actor.User() == uuid {
		goto Do
	}

	if err := service.Authorize(permChief, s.actor); err != nil {
		return err
	}

Do:
	return s.Service.Patch(uuid, req)
}

func (s *Gate) UpdatePassword(uuid uuid.UUID, req user.UpdatePassword) error {
	if s.actor.User() == uuid {
		goto Do
	}

	if err := service.Authorize(permChief, s.actor); err != nil {
		return err
	}

Do:
	return s.Service.UpdatePassword(uuid, req)
}

func (s *Gate) UpdateRole(uuid uuid.UUID, req user.UpdateRole) error {
	if err := service.Authorize(permChief, s.actor); err != nil {
		return err
	}

	return s.Service.UpdateRole(uuid, req)
}

func (s *Gate) Delete(uuid uuid.UUID) error {
	if s.actor.User() == uuid {
		goto Do
	}

	if err := service.Authorize(permChief, s.actor); err != nil {
		return err
	}

Do:
	return s.Service.Delete(uuid)
}
