package userserve

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/support"
	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Core struct {
	Users user.Repository
}

var _ user.Service = &Core{}

func (c *Core) List(req user.ListParams) (user.Entities, error) {
	return c.Users.List(req.Offset, req.Limit)
}

func (c *Core) Get(uuid uuid.UUID) (user.Entity, error) {
	return c.Users.Get(uuid)
}

func (c *Core) GetBySIAPE(siape string) (user.Entity, error) {
	return c.Users.GetBySIAPE(siape)
}

func (c *Core) Create(req user.Create) (uuid.UUID, error) {
	u, err := user.New(req.SIAPE, req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		return uuid.UUID{}, err
	}

	ent := translate(&u)

	now := time.Now()
	ent.Created = now
	ent.Updated = now

	return u.UUID(), c.Users.Create(ent)
}

func (c *Core) Patch(uuid uuid.UUID, req user.Patch) error {
	var string opt.Opt[string]
	var role opt.Opt[auth.Role]

	return patch(c.Users, uuid, req.Name, req.Email, string, role)
}

func (c *Core) UpdatePassword(uuid uuid.UUID, req user.UpdatePassword) error {
	return support.ErrTODO
}

func (c *Core) UpdateRole(uuid uuid.UUID, req user.UpdateRole) error {
	var string opt.Opt[string]

	return patch(c.Users, uuid, string, string, string, opt.Some(req.Role))
}

func (c *Core) Delete(uuid uuid.UUID) error {
	return c.Users.Delete(uuid)
}

func patch(users user.Repository, uuid uuid.UUID, name, email, password opt.Opt[string], role opt.Opt[auth.Role]) error {
	var u user.PartialEntity

	err := errors.Join(
		entity.SomeThen(&u.Name, name, user.ProcessName),
		entity.SomeThen(&u.Email, email, user.ProcessEmail),
		entity.SomeThen(&u.Password, password, user.ProcessPassword),
		entity.SomeThen(&u.Role, role, user.ProcessRole),
	)
	if err != nil {
		return user.ErrUserUpdate.Cause(err).Make()
	}

	u.Updated = time.Now()

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
