package promotionserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/promotion"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/support/service"
	entity "github.com/alan-b-lima/almodon/internal/support/service"

	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/pkg/problem"
	"github.com/alan-b-lima/pkg/scheduler"
)

type Core struct {
	Promotions promotion.Store
	Users      user.Service

	Scheduler *scheduler.Scheduler
}

var _ promotion.Service = &Core{}

func New(promotions promotion.Store, users user.Service, scheduler *scheduler.Scheduler) *Core {
	return &Core{
		Promotions: promotions,
		Users:      users,
		Scheduler:  scheduler,
	}
}

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

	var rec promotion.CreateRecord
	err := problem.Join(
		service.Set(&rec.Expires, max_age, promotion.ProcessMaxAge),
	)
	if err != nil {
		return promotion.CreateResult{}, promotion.ErrCreate.Cause(err).Make()
	}

	rec.UUID = uuid.NewUUIDv7()
	rec.User = req.User

	c.flush_at(rec.Expires)

	return promotion.CreateResult{UUID: rec.UUID}, c.Promotions.Create(ctx, rec)
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

	c.flush_at(expires)

	rec := promotion.UpdateRecord{Expires: expires}
	return c.Promotions.Update(ctx, uuid, rec)
}

func (c *Core) Delete(ctx context.Context, uuid uuid.UUID) error {
	return c.Promotions.Delete(ctx, uuid)
}

func (c *Core) Publish(ctx context.Context) error {
	err := c.Promotions.DeleteExpired(ctx, time.Now())
	if err != nil {
		return err
	}

	recs, err := c.Promotions.List(ctx)
	if err != nil {
		return err
	}

	for _, rec := range recs {
		c.flush_at(rec.Expires)
	}

	return nil
}

func (c *Core) flush_at(expires time.Time) {
	c.Scheduler.Post(func() {
		c.Promotions.DeleteExpired(context.TODO(), expires)
	}, expires)
}
