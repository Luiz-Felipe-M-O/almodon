package userserve

import (
	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/promotion"
	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/hash"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Core struct {
	Users      user.Repository
	Sessions   session.Service
	Promotions promotion.Service
}

func (c *Core) List(req user.ListParams) (user.Entities, error) {
	return c.Users.List(req.Offset, req.Limit)
}

func (c *Core) Get(uuid uuid.UUID) (user.Entity, error) {
	return c.Users.Get(uuid)
}

func (c *Core) GetBySIAPE(siape int) (user.Entity, error) {
	return c.Users.GetBySIAPE(siape)
}

func (c *Core) Create(req user.Create) (uuid.UUID, error) {
	u, err := user.New(req.SIAPE, req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		return uuid.UUID{}, err
	}

	return u.UUID(), c.Users.Create(translate(&u))
}

func (c *Core) Patch(uuid uuid.UUID, req user.Patch) error {
	var string opt.Opt[string]
	var role opt.Opt[auth.Role]

	return patch(c.Users, uuid, req.Name, req.Email, string, role)
}

func (c *Core) UpdatePassword(uuid uuid.UUID, req user.UpdatePassword) error {
	return xerrors.ErrTODO
}

func (c *Core) UpdateRole(uuid uuid.UUID, req user.UpdateRole) error {
	var string opt.Opt[string]

	return patch(c.Users, uuid, string, string, string, opt.Some(req.Role))
}

func (c *Core) Delete(uuid uuid.UUID) error {
	return c.Users.Delete(uuid)
}

func (c *Core) Authenticate(siape int, password string) (user.AuthEntity, error) {
	res, err := c.Users.GetBySIAPE(siape)
	if err != nil {
		return user.AuthEntity{}, err
	}

	if !hash.Compare(res.Password[:], []byte(password)) {
		return user.AuthEntity{}, xerrors.ErrIncorrectPassword
	}

	sres, err := c.Sessions.CreateAndGet(session.Create{User: res.UUID})
	if err != nil {
		return user.AuthEntity{}, err
	}

	ares := user.AuthEntity{
		UUID:    sres.UUID,
		User:    res.UUID,
		Expires: sres.Expires,
	}
	return ares, nil
}

func (c *Core) Actor(session uuid.UUID) (auth.Actor, error) {
	res, err := c.Sessions.Get(session)
	if err != nil {
		return auth.NewUnlogged(), xerrors.ErrUnauthenticatedUser.New(err)
	}

	ures, err := c.Users.Get(res.User)
	if err != nil {
		return auth.NewUnlogged(), xerrors.ErrUnauthenticatedUser.New(err)
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

func patch(users user.Patcher, uuid uuid.UUID, name, email, password opt.Opt[string], role opt.Opt[auth.Role]) error {
	var u user.PartialEntity

	err := errors.Join(
		entity.SomeThen(&u.Name, name, user.ProcessName),
		entity.SomeThen(&u.Email, email, user.ProcessEmail),
		entity.SomeThen(&u.Password, password, user.ProcessPassword),
		entity.SomeThen(&u.Role, role, user.ProcessRole),
	)
	if err != nil {
		return xerrors.ErrUserUpdate.New(err)
	}

	return users.Patch(uuid, u)
}

func translate(e *user.User) user.Entity {
	return user.Entity{
		UUID:     e.UUID(),
		SIAPE:    e.SIAPE(),
		Name:     e.Name(),
		Email:    e.Email(),
		Password: e.Password(),
		Role:     e.Role(),
	}
}
