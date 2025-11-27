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

var permAdmin = auth.Permit(auth.Admin)

func (s *Gate) Allow(act auth.Actor) material.Service {
	return &Gate{
		Service: s.Service,
		actor:   act,
	}
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
