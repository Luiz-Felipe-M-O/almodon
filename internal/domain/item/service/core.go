package itemserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/item"
	"github.com/alan-b-lima/almodon/internal/support/service"

	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/pkg/opt"
	"github.com/alan-b-lima/pkg/problem"
)

type Core struct {
	Items item.Store
}

var _ item.Service = (*Core)(nil)

func New(items item.Store) *Core {
	return &Core{
		Items: items,
	}
}

func (c *Core) List(ctx context.Context) ([]item.Result, error) {
	recs, err := c.Items.List(ctx)
	if err != nil {
		return nil, err
	}

	return TranslateList(recs), nil
}

func (c *Core) ListByMaterial(ctx context.Context, material uuid.UUID) ([]item.Result, error) {
	recs, err := c.Items.ListByMaterial(ctx, material)
	if err != nil {
		return nil, err
	}

	return TranslateList(recs), nil
}

func (c *Core) ListByECampus(ctx context.Context, ecampus int) ([]item.Result, error) {
	recs, err := c.Items.ListByECampus(ctx, ecampus)
	if err != nil {
		return nil, err
	}

	return TranslateList(recs), nil
}

func (c *Core) ListByCATMAT(ctx context.Context, catmat int) ([]item.Result, error) {
	recs, err := c.Items.ListByCATMAT(ctx, catmat)
	if err != nil {
		return nil, err
	}

	return TranslateList(recs), nil
}

func (c *Core) ListBySIADS(ctx context.Context, siads int) ([]item.Result, error) {
	recs, err := c.Items.ListBySIADS(ctx, siads)
	if err != nil {
		return nil, err
	}

	return TranslateList(recs), nil
}

func (c *Core) Get(ctx context.Context, id uuid.UUID) (item.Result, error) {
	rec, err := c.Items.Get(ctx, id)
	if err != nil {
		return item.Result{}, err
	}

	return Translate(&rec), nil
}

func (c *Core) History(ctx context.Context, uuid uuid.UUID) (item.HistoryResult, error) {
	rec, err := c.Items.History(ctx, uuid)
	if err != nil {
		return item.HistoryResult{}, err
	}

	res := item.HistoryResult{
		UUID:     rec.UUID,
		Version:  rec.Version,
		Created:  rec.Created,
		Updated:  rec.Updated,
		Versions: make([]item.PastResult, 0, len(rec.Versions)),
	}

	for _, version := range rec.Versions {
		res.Versions = append(res.Versions, item.PastResult(version))
	}

	return res, nil
}

func (c *Core) Create(ctx context.Context, req item.Create) (item.CreateResult, error) {
	var rec item.Entity
	err := problem.Join(
		service.Set(&rec.Amount, req.Amount, item.ProcessAmount),
		service.Set(&rec.UnitCost, req.UnitCost, item.ProcessUnitCost),
		service.Set(&rec.Expires, req.Expires, item.ProcessExpires),
	)
	if err != nil {
		return item.CreateResult{}, item.ErrCreate.Cause(err).Make()
	}

	now := time.Now()

	rec.UUID = uuid.NewUUIDv7()
	rec.Material = req.Material
	rec.Created = now
	rec.Updated = now

	return item.CreateResult{UUID: rec.UUID}, c.Items.Create(ctx, rec)
}

func (c *Core) UpdateAmount(ctx context.Context, uuid uuid.UUID, req item.UpdateAmount) error {
	amount, err := item.ProcessAmount(req.Amount)
	if err != nil {
		return item.ErrUpdate.Cause(err).Make()
	}

	rec := item.PatchEntity{
		Amount:  opt.Some(amount),
		Updated: time.Now(),
	}

	return c.Items.Patch(ctx, uuid, rec)
}

func (c *Core) Patch(ctx context.Context, uuid uuid.UUID, req item.Patch) error {
	var rec item.PatchEntity
	err := problem.Join(
		service.SetOpt(&rec.UnitCost, req.UnitCost, item.ProcessUnitCost),
		service.SetOpt(&rec.Expires, req.Expires, item.ProcessExpires),
	)
	if err != nil {
		return item.ErrUpdate.Cause(err).Make()
	}

	rec.Material = req.Material
	rec.Updated = time.Now()

	return c.Items.Patch(ctx, uuid, rec)
}

func (c *Core) Delete(ctx context.Context, uuid uuid.UUID) error {
	return c.Items.Delete(ctx, uuid)
}

func Translate(rec *item.Record) item.Result {
	return item.Result{
		UUID:        rec.UUID,
		Name:        rec.Name,
		ECampus:     rec.ECampus,
		CATMAT:      rec.CATMAT,
		SIADS:       rec.SIADS,
		Material:    rec.Material,
		Amount:      rec.Amount,
		AmountFlag:  item.StatusAmount(rec.Amount, rec.Min),
		UnitCost:    rec.UnitCost,
		Unit:        rec.Unit,
		Expires:     rec.Expires,
		ExpiresFlag: item.StatusExpires(rec.Expires),
		Created:     rec.Created,
		Updated:     rec.Updated,
	}
}

func TranslateList(recs []item.Record) []item.Result {
	res := make([]item.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, Translate(&rec))
	}

	return res
}
