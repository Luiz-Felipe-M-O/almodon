package sessionserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/support"
	"github.com/alan-b-lima/almodon/internal/support/entity"

	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/pkg/problem"
	"github.com/alan-b-lima/pkg/scheduler"
)

type Core struct {
	Sessions  session.Store
	Scheduler *scheduler.Scheduler
}

var _ session.Service = (*Core)(nil)

func New(sessions session.Store, scheduler *scheduler.Scheduler) (*Core, error) {
	if scheduler == nil {
		return nil, support.ErrNilPointer.Message("scheduler required").Make()
	}

	return &Core{
		Sessions:  sessions,
		Scheduler: scheduler,
	}, nil
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
	maxAge := _MaxAge
	if v, ok := req.MaxAge.Unwrap(); ok {
		maxAge = v
	}

	s, err := session.New(req.User, maxAge)
	if err != nil {
		return session.Result{}, err
	}

	ss := session.CreateRecord{
		UUID:    s.UUID,
		User:    s.User,
		Renewed: s.Renewed,
		Expires: s.Expires,
		Created: time.Now(),
	}

	return session.Result(ss), c.Sessions.Create(ctx, ss)
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
			entity.Set(&rec.Renewed, s.Renewed, session.ProcessRenewed),
			entity.Set(&rec.Expires, max_age, session.ProcessMaxAge),
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

	c.post(uuid, expires)
	return nil
}

func (c *Core) Delete(ctx context.Context, uuid uuid.UUID) error {
	return c.Sessions.Delete(ctx, uuid)
}

func (c *Core) post(uuid uuid.UUID, expires time.Time) {
	c.Scheduler.Post(func() {
		ctx := context.TODO()

		c.Sessions.RunTx(ctx, func(sessions session.Store) error {
			session, err := sessions.Get(ctx, uuid)
			if err != nil {
				return err
			}

			if time.Now().Before(session.Expires) {
				return sessions.Delete(ctx, uuid)
			}
			return nil
		})
	}, expires)
}
