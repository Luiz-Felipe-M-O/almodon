package authserve

import (
	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/promotion"
	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/domain/user"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Core struct {
	Users      user.Getter
	Sessions   session.Service
	Promotions promotion.Getter
}

var _ auth.Service = &Core{}

func (c *Core) Login(siape, password string) (auth.Result, error) {
	res, err := c.Users.GetBySIAPE(siape)
	if err != nil {
		return auth.Result{}, err
	}

	if err := user.ComparePassword(res.Password[:], []byte(password)); err != nil {
		return auth.Result{}, err
	}

	sres, err := c.Sessions.CreateAndGet(session.Create{User: res.UUID})
	if err != nil {
		return auth.Result{}, err
	}

	ares := auth.Result{
		UUID:    sres.UUID,
		User:    res.UUID,
		Expires: sres.Expires,
	}
	return ares, nil
}

func (c *Core) Logout(session uuid.UUID) error {
	return c.Sessions.Delete(session)
}

func (c *Core) Actor(session uuid.UUID) (auth.Actor, error) {
	res, err := c.Sessions.Get(session)
	if err != nil {
		return auth.NewUnlogged(), auth.ErrUnauthenticated.Cause(err).Make()
	}

	ures, err := c.Users.Get(res.User)
	if err != nil {
		return auth.NewUnlogged(), auth.ErrUnauthenticated.Cause(err).Make()
	}

	role := ures.Role
	if ures.Role == auth.Admin {
		_, err := c.Promotions.GetByUser(ures.UUID)
		if err == nil {
			role = auth.Promoted
		}
	}

	return auth.NewLogged(
		ures.UUID,
		role,
	), nil
}
