package promotionserve

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/promotion"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Core struct {
	Promotions promotion.Repository
	Users      user.Service
}

const _MaxAge = 1 * 24 * time.Hour

func (c *Core) List(req promotion.ListParams) (promotion.Entities, error) {
	return c.Promotions.List(req.Offset, req.Limit)
}

func (c *Core) Get(uuid uuid.UUID) (promotion.Entity, error) {
	res, err := c.Promotions.Get(uuid)
	if err != nil {
		return promotion.Entity{}, err
	}

	if time.Now().After(res.Expires) {
		return promotion.Entity{}, xerrors.ErrPromotionNotFound
	}

	return res, err
}

func (c *Core) GetByUser(user uuid.UUID) (promotion.Entity, error) {
	res, err := c.Promotions.GetByUser(user)
	if err != nil {
		return promotion.Entity{}, err
	}

	if time.Now().After(res.Expires) {
		return promotion.Entity{}, xerrors.ErrPromotionNotFound
	}

	return res, err
}

// TODO: verify validity of _MaxAge and turn it to an internal error
func (c *Core) Create(req promotion.Create) (uuid.UUID, error) {
	maxAge := _MaxAge
	if v, ok := req.MaxAge.Unwrap(); ok {
		maxAge = v
	}

	if _, err := c.Users.Get(req.User); err != nil {
		return uuid.UUID{}, err
	}

	p, err := promotion.New(req.User, maxAge)
	if err != nil {
		return uuid.UUID{}, err
	}

	return p.UUID(), c.Promotions.Create(translate(&p))
}

// TODO: verify validity of _MaxAge and turn it to an internal error
func (c *Core) Update(uuid uuid.UUID, req promotion.Update) error {
	maxAge := _MaxAge
	if v, ok := req.MaxAge.Unwrap(); ok {
		maxAge = v
	}

	var p promotion.Promotion
	if err := p.SetMaxAge(maxAge); err != nil {
		return err
	}

	return c.Promotions.Update(uuid, p.Expires())
}

func (c *Core) Delete(uuid uuid.UUID) error {
	return c.Promotions.Delete(uuid)
}

func translate(p *promotion.Promotion) promotion.Entity {
	return promotion.Entity{
		UUID:    p.UUID(),
		User:    p.User(),
		Expires: p.Expires(),
	}
}
