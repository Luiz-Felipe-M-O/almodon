package userserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/support"
	"github.com/alan-b-lima/almodon/internal/support/entity"

	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/pkg/problem"
)

type Core struct {
	Users user.Store
}

var _ user.Service = (*Core)(nil)

func New(users user.Store) *Core {
	return &Core{
		Users: users,
	}
}

func (c *Core) List(ctx context.Context) ([]user.Result, error) {
	recs, err := c.Users.List(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]user.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, user.Result(rec))
	}

	return res, nil
}

func (c *Core) Get(ctx context.Context, uuid uuid.UUID) (user.Result, error) {
	rec, err := c.Users.Get(ctx, uuid)
	if err != nil {
		return user.Result{}, err
	}

	return user.Result(rec), nil
}

func (c *Core) GetBySIAPE(ctx context.Context, siape string) (user.Result, error) {
	rec, err := c.Users.GetBySIAPE(ctx, siape)
	if err != nil {
		return user.Result{}, err
	}

	return user.Result(rec), nil
}

func (c *Core) Me(ctx context.Context) (user.Result, error) {
	return user.Result{}, support.ErrTODO
}

func (c *Core) Create(ctx context.Context, req user.Create) (user.CreateResult, error) {
	u, err := user.New(req.SIAPE, req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		return user.CreateResult{}, err
	}

	rec := user.CreateRecord{
		UUID:     u.UUID,
		SIAPE:    u.SIAPE,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Role:     u.Role,
		Created:  time.Now(),
		Updated:  time.Now(),
	}

	return user.CreateResult{UUID: u.UUID}, c.Users.Create(ctx, rec)
}

func (c *Core) Patch(ctx context.Context, uuid uuid.UUID, req user.Patch) error {
	var rec user.PatchRecord
	err := problem.Join(
		entity.SetOpt(&rec.Name, req.Name, user.ProcessName),
		entity.SetOpt(&rec.Email, req.Email, user.ProcessEmail),
	)
	if err != nil {
		return user.ErrUpdate.Cause(err).Make()
	}

	rec.Updated = time.Now()

	return c.Users.Patch(ctx, uuid, rec)
}

func (c *Core) Delete(ctx context.Context, uuid uuid.UUID) error {
	return c.Users.Delete(ctx, uuid)
}
