package promotionserve

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/promotion"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Gate struct {
	promotion.Service
	actor auth.Actor
}

func New(promotions promotion.Service) promotion.Service {
	return &Gate{
		Service: promotions,
	}
}

var permChief = auth.Permit(auth.Chief)

func (s *Gate) Allow(act auth.Actor) promotion.Service {
	return &Gate{
		Service: s.Service,
		actor:   act,
	}
}

func (a *Gate) List(req promotion.ListParams) (promotion.Entities, error) {
	if err := service.Authorize(permChief, a.actor); err != nil {
		return promotion.Entities{}, err
	}

	return a.Service.List(req)
}

func (a *Gate) Get(uuid uuid.UUID) (promotion.Entity, error) {
	if err := service.Authorize(permChief, a.actor); err != nil {
		return promotion.Entity{}, err
	}

	return a.Service.Get(uuid)
}

func (a *Gate) Create(req promotion.Create) (uuid.UUID, error) {
	if err := service.Authorize(permChief, a.actor); err != nil {
		return uuid.UUID{}, err
	}

	return a.Service.Create(req)
}

func (a *Gate) Update(uuid uuid.UUID, req promotion.Update) error {
	if err := service.Authorize(permChief, a.actor); err != nil {
		return err
	}

	return a.Service.Update(uuid, req)
}

func (a *Gate) Delete(uuid uuid.UUID) error {
	if err := service.Authorize(permChief, a.actor); err != nil {
		return err
	}

	return a.Service.Delete(uuid)
}
