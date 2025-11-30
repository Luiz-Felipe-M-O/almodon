package itemserve

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/item"
	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/opt"
	uuidpkg "github.com/alan-b-lima/almodon/pkg/uuid"
)

type Core struct {
	Items item.Repository
}

var _ item.Service = &Core{}

func (c *Core) List(req item.ListParams) (item.Entities, error) {
	return c.Items.List(req.Offset, req.Limit)
}

func (c *Core) ListByMaterial(uuid uuidpkg.UUID) (item.Entities, error) {
	return c.Items.ListByMaterial(uuid)
}

func (c *Core) ListBySupplier(uuid uuidpkg.UUID) (item.Entities, error) {
	return c.Items.ListBySupplier(uuid)
}

func (c *Core) Get(uuid uuidpkg.UUID) (item.Entity, error) {
	return c.Items.Get(uuid)
}

func (c *Core) Create(req item.Create) (uuidpkg.UUID, error) {
	batch, err := item.New(
		req.Material,
		req.Supplier,
		req.Quantity,
		req.UnitCost,
		req.Arrival,
		req.Expiration,
		req.Invoice,
		req.Lot,
		req.Notes,
	)
	if err != nil {
		return uuidpkg.UUID{}, err
	}

	ent := translate(&batch)

	now := time.Now()
	ent.Created = now
	ent.Updated = now

	return batch.UUID(), c.Items.Create(ent)
}

func (c *Core) Patch(uuid uuidpkg.UUID, req item.Patch) error {
	return patch(
		c.Items,
		uuid,
		req.Material,
		req.Supplier,
		req.Quantity,
		req.UnitCost,
		req.Arrival,
		req.Expiration,
		req.Invoice,
		req.Lot,
		req.Notes,
	)
}

func (c *Core) UpdateQuantity(uuid uuidpkg.UUID, req item.UpdateQuantity) error {
	var emptyUUID opt.Opt[uuidpkg.UUID]
	var emptyFloat opt.Opt[float64]
	var emptyTime opt.Opt[time.Time]
	var emptyString opt.Opt[string]

	return patch(
		c.Items,
		uuid,
		emptyUUID,
		emptyUUID,
		opt.Some(req.Quantity),
		emptyFloat,
		emptyTime,
		emptyTime,
		emptyString,
		emptyString,
		emptyString,
	)
}

func (c *Core) Delete(uuid uuidpkg.UUID) error {
	return c.Items.Delete(uuid)
}

func patch(
	repo item.Patcher,
	uuid uuidpkg.UUID,
	material, supplier opt.Opt[uuidpkg.UUID],
	quantity, unitCost opt.Opt[float64],
	arrival, expiration opt.Opt[time.Time],
	invoice, lot, notes opt.Opt[string],
) error {
	var pi item.PartialEntity

	err := errors.Join(
		entity.SomeThen(&pi.Material, material, item.ProcessMaterial),
		entity.SomeThen(&pi.Supplier, supplier, item.ProcessSupplier),
		entity.SomeThen(&pi.Quantity, quantity, item.ProcessQuantity),
		entity.SomeThen(&pi.UnitCost, unitCost, item.ProcessUnitCost),
		entity.SomeThen(&pi.Arrival, arrival, item.ProcessArrival),
		entity.SomeThen(&pi.Expiration, expiration, item.ProcessExpiration),
		entity.SomeThen(&pi.Invoice, invoice, item.ProcessInvoice),
		entity.SomeThen(&pi.Lot, lot, item.ProcessLot),
		entity.SomeThen(&pi.Notes, notes, item.ProcessNotes),
	)
	if err != nil {
		return err
	}

	pi.Updated = time.Now()

	return repo.Patch(uuid, pi)
}

func translate(i *item.ItemBatch) item.Entity {
	return item.Entity{
		UUID:       i.UUID(),
		Material:   i.Material(),
		Supplier:   i.Supplier(),
		Quantity:   i.Quantity(),
		UnitCost:   i.UnitCost(),
		Arrival:    i.Arrival(),
		Expiration: i.Expiration(),
		Invoice:    i.Invoice(),
		Lot:        i.Lot(),
		Notes:      i.Notes(),
	}
}
