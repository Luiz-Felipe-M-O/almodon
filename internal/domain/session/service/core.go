package sessionserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/support/service"

	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/pkg/problem"
	"github.com/alan-b-lima/pkg/scheduler"
)

type Core struct {
	Sessions  session.Store
	Scheduler *scheduler.Scheduler
}

var _ session.Service = (*Core)(nil)

func New(sessions session.Store, scheduler *scheduler.Scheduler) *Core {
	return &Core{
		Sessions:  sessions,
		Scheduler: scheduler,
	}
}

const _MaxAge = 1 * time.Hour

func (c *Core) Get(ctx context.Context, uuid uuid.UUID) (session.Result, error) {
	res, err := c.Sessions.Get(ctx, uuid)
	if err != nil {
		return session.Result{}, err
	}

	if time.Now().After(res.Expires) {
		return session.Result{}, session.ErrNotFound
	}

	return session.Result(res), nil
}

// TODO: verify validity of [_MaxAge] and turn it to an internal error
func (c *Core) CreateAndGet(ctx context.Context, req session.Create) (session.Result, error) {
	max_age := _MaxAge
	if v, ok := req.MaxAge.Unwrap(); ok {
		max_age = v
	}

	var rec session.CreateRecord
	err := problem.Join(
		service.Set(&rec.Renewed, 0, session.ProcessRenewed),
		service.Set(&rec.Expires, max_age, session.ProcessMaxAge),
	)
	if err != nil {
		return session.Result{}, session.ErrCreate.Cause(err).Make()
	}

	rec.User = req.User
	rec.UUID = uuid.NewUUIDv7()
	rec.Created = time.Now()

	err = c.Sessions.RunTx(ctx, func(store session.Store) error {
		s, err := store.GetByUser(ctx, req.User)
		if err != session.ErrNotFound {
			if err != nil {
				return err
			}

			if err := store.Delete(ctx, s.UUID); err != nil {
				return err
			}
		}

		return store.Create(ctx, rec)
	})
	if err != nil {
		return session.Result{}, err
	}

	c.flush_at(rec.Expires)

	return session.Result(rec), nil
}

// TODO: verify validity of _MaxAge and turn it to an internal error
func (c *Core) Update(ctx context.Context, uuid uuid.UUID, req session.Update) error {
	max_age := _MaxAge
	if v, ok := req.MaxAge.Unwrap(); ok {
		max_age = v
	}

	var expires time.Time
	err := c.Sessions.RunTx(ctx, func(store session.Store) error {
		s, err := store.Get(ctx, uuid)
		if err != nil {
			return err
		}

		var rec session.UpdateRecord
		err = problem.Join(
			service.Set(&rec.Renewed, s.Renewed, session.ProcessRenewed),
			service.Set(&rec.Expires, max_age, session.ProcessMaxAge),
		)
		if err != nil {
			return session.ErrUpdate.Cause(err).Make()
		}

		expires = rec.Expires
		return c.Sessions.Update(ctx, uuid, rec)
	})
	if err != nil {
		return err
	}

	c.flush_at(expires)
	return nil
}

func (c *Core) Delete(ctx context.Context, uuid uuid.UUID) error {
	return c.Sessions.Delete(ctx, uuid)
}

func (c *Core) flush_at(expires time.Time) {
	c.Scheduler.Post(func() {
		c.Sessions.DeleteExpired(context.TODO(), time.Now())
	}, expires)
}
