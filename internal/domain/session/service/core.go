package sessionserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/pkg/problem"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Core struct {
	Sessions session.Store
}

var _ session.Service = (*Core)(nil)

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

		return c.Sessions.Update(ctx, uuid, rec)
	})

	return err
}

func (c *Core) Delete(ctx context.Context, uuid uuid.UUID) error {
	return c.Sessions.Delete(ctx, uuid)
}
