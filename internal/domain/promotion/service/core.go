package promotionserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/promotion"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/support/entity"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Core struct {
	Promotions promotion.Store
	Users      user.Service
}

var _ promotion.Service = &Core{}

const _MaxAge = 1 * 24 * time.Hour

func (c *Core) Get(ctx context.Context, uuid uuid.UUID) (promotion.Result, error) {
	res, err := c.Promotions.Get(ctx, uuid)
	if err != nil {
		return promotion.Result{}, err
	}

	if time.Now().After(res.Expires) {
		return promotion.Result{}, promotion.ErrNotFound
	}

	return promotion.Result(res), err
}

func (c *Core) GetByUser(ctx context.Context, user uuid.UUID) (promotion.Result, error) {
	res, err := c.Promotions.GetByUser(ctx, user)
	if err != nil {
		return promotion.Result{}, err
	}

	if time.Now().After(res.Expires) {
		return promotion.Result{}, promotion.ErrNotFound
	}

	return promotion.Result(res), err
}

// TODO: verify validity of _MaxAge and turn it to an internal error
func (c *Core) Create(ctx context.Context, req promotion.Create) (promotion.CreateResult, error) {
	max_age := _MaxAge
	if v, ok := req.MaxAge.Unwrap(); ok {
		max_age = v
	}

	if _, err := c.Users.Get(ctx, req.User); err != nil {
		return promotion.CreateResult{}, err
	}

	p, err := promotion.New(req.User, max_age)
	if err != nil {
		return promotion.CreateResult{}, err
	}

	nreq := promotion.CreateRecord(p)
	return promotion.CreateResult{UUID: p.UUID}, c.Promotions.Create(ctx, nreq)
}

// TODO: verify validity of _MaxAge and turn it to an internal error
func (c *Core) Update(ctx context.Context, uuid uuid.UUID, req promotion.Update) error {
	max_age := _MaxAge
	if v, ok := req.MaxAge.Unwrap(); ok {
		max_age = v
	}

	var expires time.Time
	if err := entity.Set(&expires, max_age, promotion.ProcessMaxAge); err != nil {
		return err
	}

	nreq := promotion.UpdateRecord{Expires: expires}
	return c.Promotions.Update(ctx, uuid, nreq)
}

func (c *Core) Delete(ctx context.Context, uuid uuid.UUID) error {
	return c.Promotions.Delete(ctx, uuid)
}
