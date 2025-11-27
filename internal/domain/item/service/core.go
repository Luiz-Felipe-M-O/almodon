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

func (c *Core) ListByBatch(uuid uuidpkg.UUID) (item.Entities, error) {
	return c.Items.ListByBatch(uuid)
}

func (c *Core) ListByMaterial(uuid uuidpkg.UUID) (item.Entities, error) {
	return c.Items.ListByMaterial(uuid)
}

func (c *Core) Get(uuid uuidpkg.UUID) (item.Entity, error) {
	return c.Items.Get(uuid)
}

func (c *Core) Create(req item.Create) (uuidpkg.UUID, error) {
	u, err := item.New(req.Batch, req.Material, req.Quantity, req.Expiration)
	if err != nil {
		return uuidpkg.UUID{}, err
	}

	ent := translate(&u)

	now := time.Now()
	ent.Created = now
	ent.Updated = now

	return u.UUID(), c.Items.Create(ent)
}

func (c *Core) Patch(uuid uuidpkg.UUID, req item.Patch) error {
	return patch(c.Items, uuid, req.Batch, req.Material, req.Quantity, req.Expiration)
}

func (c *Core) UpdateQuantity(uuid uuidpkg.UUID, req item.UpdateQuantity) error {
	var uuid_ opt.Opt[uuidpkg.UUID]
	var expiration opt.Opt[time.Time]

	return patch(c.Items, uuid, uuid_, uuid_, opt.Some(req.Quantity), expiration)
}

func (c *Core) Delete(uuid uuidpkg.UUID) error {
	return c.Items.Delete(uuid)
}

func patch(repo item.Patcher, uuid uuidpkg.UUID, batch, material opt.Opt[uuidpkg.UUID], quantity opt.Opt[float64], expiration opt.Opt[time.Time]) error {
	var pi item.PartialEntity

	err := errors.Join(
		entity.SomeThen(&pi.Batch, batch, item.ProcessBatch),
		entity.SomeThen(&pi.Material, material, item.ProcessMaterial),
		entity.SomeThen(&pi.Quantity, quantity, item.ProcessQuantity),
		entity.SomeThen(&pi.Expiration, expiration, item.ProcessExpiration),
	)
	if err != nil {
		return err
	}

	pi.Updated = time.Now()

	return repo.Patch(uuid, pi)
}

func translate(i *item.ItemUnit) item.Entity {
	return item.Entity{
		UUID:       i.UUID(),
		Batch:      i.Batch(),
		Material:   i.Material(),
		Quantity:   i.Quantity(),
		Expiration: i.Expiration(),
	}
}
