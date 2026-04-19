package itemserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/item"
	"github.com/alan-b-lima/almodon/internal/support/service"

	"github.com/alan-b-lima/almodon/pkg/uuid"

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
	return translate_list(c.Items.List(ctx))
}

func (c *Core) ListByMaterial(ctx context.Context, material uuid.UUID) ([]item.Result, error) {
	return translate_list(c.Items.ListByMaterial(ctx, material))
}

func (c *Core) ListByECampus(ctx context.Context, ecampus int) ([]item.Result, error) {
	return translate_list(c.Items.ListByECampus(ctx, ecampus))
}

func (c *Core) ListByCATMAT(ctx context.Context, catmat int) ([]item.Result, error) {
	return translate_list(c.Items.ListByCATMAT(ctx, catmat))
}

func (c *Core) ListBySIADS(ctx context.Context, siads int) ([]item.Result, error) {
	return translate_list(c.Items.ListBySIADS(ctx, siads))
}

func (c *Core) Get(ctx context.Context, id uuid.UUID) (item.Result, error) {
	rec, err := c.Items.Get(ctx, id)
	if err != nil {
		return item.Result{}, err
	}

	return translate(&rec), nil
}

func (c *Core) Create(ctx context.Context, req item.Create) (item.CreateResult, error) {
	var rec item.CreateRecord
	err := problem.Join(
		service.Set(&rec.Amount, req.Amount, item.ProcessAmount),
		service.Set(&rec.UnitCost, req.UnitCost, item.ProcessUnitCost),
		service.Set(&rec.Arrival, req.Arrival, item.ProcessArrival),
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
		return err
	}

	return c.Items.UpdateAmount(ctx, uuid, amount)
}

func (c *Core) Patch(ctx context.Context, uuid uuid.UUID, req item.Patch) error {
	var rec item.PatchRecord
	err := problem.Join(
		service.SetOpt(&rec.UnitCost, req.UnitCost, item.ProcessUnitCost),
		service.SetOpt(&rec.Arrival, req.Arrival, item.ProcessArrival),
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

func translate_list(recs []item.Record, err error) ([]item.Result, error) {
	if err != nil {
		return nil, err
	}

	res := make([]item.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, translate(&rec))
	}

	return res, nil
}

func translate(rec *item.Record) item.Result {
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
		Arrival:     rec.Arrival,
		Expires:     rec.Expires,
		ExpiresFlag: item.StatusExpires(rec.Expires),
		Created:     rec.Created,
		Updated:     rec.Updated,
	}
}
