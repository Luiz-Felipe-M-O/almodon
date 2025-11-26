package material

import (
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

func List(materials Lister, offset, limit int) (Entities, error) {
	return materials.List(offset, limit)
}

func Get(materials Getter, uuid uuid.UUID) (Entity, error) {
	return materials.Get(uuid)
}

func ListBySIADS(materials ListerBySIADS, siads string) (Entities, error) {
	return materials.ListBySIADS(siads)
}

func ListByCATMAT(materials ListerByCATMAT, catmat string) (Entities, error) {
	return materials.ListByCATMAT(catmat)
}

func ListByECampus(materials ListerByECAMPUS, ecampus string) (Entities, error) {
	return materials.ListByECampus(ecampus)
}

func Create(materials Creater, name, siads, catmat, ecampus, description,
	unit string, minQuantity float64) (uuid.UUID, error) {
	m, err := New(name, siads, catmat, ecampus, description, unit, minQuantity)
	if err != nil {
		return uuid.UUID{}, err
	}

	return m.UUID(), materials.Create(translate(&m))
}

func Patch(materials Patcher, uuid uuid.UUID, name, siads, catmat, ecampus,
	description, unit opt.Opt[string], minQuantity opt.Opt[float64]) error {
	var pm PartialEntity

	err := errors.Join(
		someThen(&pm.Name, name, ProcessName),
		someThen(&pm.SIADS, siads, ProcessSIADS),
		someThen(&pm.CATMAT, catmat, ProcessCATMAT),
		someThen(&pm.ECAMPUS, ecampus, ProcessECAMPUS),
		someThen(&pm.Description, description, ProcessDescription),
		someThen(&pm.Unit, unit, ProcessUnit),
		someThen(&pm.MinQuantity, minQuantity, ProcessMinQuantity),
	)
	if err != nil {
		return xerrors.ErrMaterialUpdate.New(err)
	}

	return materials.Patch(uuid, pm)
}

func Delete(materials Deleter, uuid uuid.UUID) error {
	return materials.Delete(uuid)
}

func translate(m *Material) Entity {
	return Entity{
		UUID:        m.UUID(),
		Name:        m.Name(),
		SIADS:       m.SIADS(),
		CATMAT:      m.CATMAT(),
		ECAMPUS:     m.ECampus(),
		Description: m.Description(),
		Unit:        m.Unit(),
		MinQuantity: m.MinQuantity(),
		CreatedAt:   m.CreatedAt(),
		UpdatedAt:   m.UpdatedAt(),
	}
}

func someThen[F, R any](dst *opt.Opt[R], src opt.Opt[F], fn func(F) (R, error)) error {
	val, ok := src.Unwrap()
	if !ok {
		return nil
	}

	res, err := fn(val)
	if err != nil {
		return err
	}

	*dst = opt.Some(res)
	return nil
}
