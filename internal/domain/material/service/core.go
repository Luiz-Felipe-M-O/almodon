package materialserve

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/material"
	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Core struct {
	Materials material.Repository
}

var _ material.Service = &Core{}

func (c *Core) List(req material.ListParams) (material.Entities, error) {
	return c.Materials.List(req.Offset, req.Limit)
}

func (c *Core) ListByCATMAT(catmat string, req material.ListParams) (material.Entities, error) {
	return c.Materials.ListByCATMAT(catmat)
}

func (c *Core) ListByECampus(ecampus string, req material.ListParams) (material.Entities, error) {
	return c.Materials.ListByECampus(ecampus)
}

func (c *Core) ListBySIADS(siads string, req material.ListParams) (material.Entities, error) {
	return c.Materials.ListBySIADS(siads)
}

func (c *Core) Get(uuid uuid.UUID) (material.Entity, error) {
	return c.Materials.Get(uuid)
}

func (c *Core) Create(req material.Create) (uuid.UUID, error) {
	m, err := material.New(req.Name, req.SIADS, req.CATMAT, req.ECampus, req.Description, req.Unit, req.MinQuantity)
	if err != nil {
		return uuid.UUID{}, err
	}

	ent := translate(&m)

	now := time.Now()
	ent.Created = now
	ent.Updated = now

	return m.UUID(), c.Materials.Create(ent)
}

func (c *Core) Patch(uuid uuid.UUID, req material.Patch) error {
	var m material.PartialEntity

	err := errors.Join(
		entity.SomeThen(&m.Name, req.Name, material.ProcessName),
		entity.SomeThen(&m.SIADS, req.SIADS, material.ProcessSIADS),
		entity.SomeThen(&m.CATMAT, req.CATMAT, material.ProcessCATMAT),
		entity.SomeThen(&m.ECampus, req.ECampus, material.ProcessECampus),
		entity.SomeThen(&m.Description, req.Description, material.ProcessDescription),
		entity.SomeThen(&m.Unit, req.Unit, material.ProcessUnit),
		entity.SomeThen(&m.MinQuantity, req.MinQuantity, material.ProcessMinQuantity),
	)
	if err != nil {
		return xerrors.ErrUserUpdate.New(err)
	}

	m.Updated = time.Now()

	return c.Materials.Patch(uuid, m)
}

func (c *Core) Delete(uuid uuid.UUID) error {
	return c.Materials.Delete(uuid)
}

func translate(e *material.Material) material.Entity {
	return material.Entity{
		UUID:        e.UUID(),
		Name:        e.Name(),
		SIADS:       e.SIADS(),
		CATMAT:      e.CATMAT(),
		ECampus:     e.ECampus(),
		Description: e.Description(),
		Unit:        e.Unit(),
		MinQuantity: e.MinQuantity(),
	}
}
