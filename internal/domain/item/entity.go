package item

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type ItemUnit struct {
	uuid     uuid.UUID
	batch    uuid.UUID
	material uuid.UUID
	quantity float64

	// zero value, [time.Time.IsZero](), means, by convention, that
	// an item does not expire.
	expiration time.Time
}

func New(batch, material uuid.UUID, quantity float64, expiration time.Time) (ItemUnit, error) {
	var i ItemUnit

	err := errors.Join(
		i.SetBatch(batch),
		i.SetMaterial(material),
		i.SetQuantity(quantity),
		i.SetExpiration(expiration),
	)
	if err != nil {
		return ItemUnit{}, xerrors.ErrItemCreation.New(err)
	}

	i.uuid = uuid.NewUUIDv7()
	return i, nil
}

func (i *ItemUnit) UUID() uuid.UUID       { return i.uuid }
func (i *ItemUnit) Batch() uuid.UUID      { return i.batch }
func (i *ItemUnit) Material() uuid.UUID   { return i.material }
func (i *ItemUnit) Quantity() float64     { return i.quantity }
func (i *ItemUnit) Expiration() time.Time { return i.expiration }

func IsExpired(expiration time.Time) bool {
	if expiration.IsZero() {
		return false
	}
	return time.Now().After(expiration)
}

func HasExpiration(expiration time.Time) bool {
	return !expiration.IsZero()
}

func (i *ItemUnit) SetBatch(batch uuid.UUID) error {
	return entity.Set(&i.batch, batch, ProcessBatch)
}

func (i *ItemUnit) SetMaterial(material uuid.UUID) error {
	return entity.Set(&i.material, material, ProcessMaterial)
}

func (i *ItemUnit) SetQuantity(quantity float64) error {
	return entity.Set(&i.quantity, quantity, ProcessQuantity)
}

func (i *ItemUnit) SetExpiration(expiration time.Time) error {
	return entity.Set(&i.expiration, expiration, ProcessExpiration)
}

func ProcessBatch(batch uuid.UUID) (uuid.UUID, error) {
	if batch.IsNil() {
		return uuid.UUID{}, xerrors.ErrBatchNotFound
	}
	return batch, nil
}

func ProcessMaterial(material uuid.UUID) (uuid.UUID, error) {
	if material.IsNil() {
		return uuid.UUID{}, xerrors.ErrMaterialNotFound
	}
	return material, nil
}

func ProcessQuantity(quantity float64) (float64, error) {
	if quantity <= 0 {
		return 0, xerrors.ErrQuantityMustBePositive
	}
	return quantity, nil
}

func ProcessExpiration(expiration time.Time) (time.Time, error) {
	return expiration, nil
}
