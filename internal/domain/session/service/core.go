package sessionserve

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Core struct {
	Sessions session.Repository
}

const _MaxAge = 10 * time.Minute

func (c *Core) Get(uuid uuid.UUID) (session.Entity, error) {
	res, err := c.Sessions.Get(uuid)
	if err != nil {
		return session.Entity{}, err
	}

	if time.Now().After(res.Expires) {
		return session.Entity{}, xerrors.ErrSessionNotFound
	}

	return res, nil
}

// TODO: verify validity of [_MaxAge] and turn it to an internal error
func (c *Core) CreateAndGet(req session.Create) (session.Entity, error) {
	maxAge := _MaxAge
	if v, ok := req.MaxAge.Unwrap(); ok {
		maxAge = v
	}

	s, err := session.New(req.User, maxAge)
	if err != nil {
		return session.Entity{}, err
	}

	session := session.Entity{
		UUID:    s.UUID(),
		User:    s.User(),
		Expires: s.Expires(),
	}

	return session, c.Sessions.Create(session)
}

// TODO: verify validity of _MaxAge and turn it to an internal error
func (c *Core) Update(uuid uuid.UUID, req session.Update) error {
	maxAge := _MaxAge
	if v, ok := req.MaxAge.Unwrap(); ok {
		maxAge = v
	}

	var s session.Session
	s.SetMaxAge(maxAge)

	return c.Sessions.Update(uuid, s.Expires())
}

func (c *Core) Delete(uuid uuid.UUID) error {
	return c.Sessions.Delete(uuid)
}
