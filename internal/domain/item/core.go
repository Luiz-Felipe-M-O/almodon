package item

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

func List(items Lister, offset, limit int) (Entities, error) {
	return items.List(offset, limit)
}

func Get(items Getter, uuid uuid.UUID) (Entity, error) {
	return items.Get(uuid)
}

func ListByBatch(items ListerByBatch, batch uuid.UUID) (Entities, error) {
	return items.ListByBatch(batch)
}

func ListByMaterial(items ListerByMaterial, material uuid.UUID) (Entities, error) {
	return items.ListByMaterial(material)
}

func Create(items Creater, batch, material uuid.UUID, quantity float64, expiration time.Time) (uuid.UUID, error) {
	i, err := New(batch, material, quantity, expiration)
	if err != nil {
		return uuid.UUID{}, err
	}

	return i.UUID(), items.Create(translate(&i))
}

func Patch(items Patcher, uuid uuid.UUID, batch, material opt.Opt[uuid.UUID],
	quantity opt.Opt[float64], expiration opt.Opt[time.Time]) error {
	var pi PartialEntity

	err := errors.Join(
		someThen(&pi.Batch, batch, ProcessBatch),
		someThen(&pi.Material, material, ProcessMaterial),
		someThen(&pi.Quantity, quantity, ProcessQuantity),
		someThen(&pi.Expiration, expiration, ProcessExpiration),
	)
	if err != nil {
		return xerrors.ErrItemUpdate.New(err)
	}

	return items.Patch(uuid, pi)
}

func Delete(items Deleter, uuid uuid.UUID) error {
	return items.Delete(uuid)
}

func translate(i *ItemUnit) Entity {
	return Entity{
		UUID:       i.UUID(),
		Batch:      i.Batch(),
		Material:   i.Material(),
		Quantity:   i.Quantity(),
		Expiration: i.Expiration(),
		CreatedAt:  i.CreatedAt(),
		UpdatedAt:  i.UpdatedAt(),
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
