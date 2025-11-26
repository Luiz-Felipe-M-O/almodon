package item

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type ItemUnit struct {
	uuid       uuid.UUID
	batch      uuid.UUID
	material   uuid.UUID
	quantity   float64
	expiration time.Time // isZero() means, by convention, that an item does not expire
	createdAt  time.Time
	updatedAt  time.Time
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
	i.createdAt = time.Now()
	i.updatedAt = time.Now()
	return i, nil
}

func (i *ItemUnit) IsExpired() bool {
	if i.expiration.IsZero() {
		return false
	}
	return time.Now().After(i.expiration)
}

func (i *ItemUnit) HasExpiration() bool {
	return !i.expiration.IsZero()
}

func (i *ItemUnit) UUID() uuid.UUID       { return i.uuid }
func (i *ItemUnit) Batch() uuid.UUID      { return i.batch }
func (i *ItemUnit) Material() uuid.UUID   { return i.material }
func (i *ItemUnit) Quantity() float64     { return i.quantity }
func (i *ItemUnit) Expiration() time.Time { return i.expiration }
func (i *ItemUnit) CreatedAt() time.Time  { return i.createdAt }
func (i *ItemUnit) UpdatedAt() time.Time  { return i.updatedAt }

func (i *ItemUnit) SetBatch(batch uuid.UUID) error {
	return entity.SetWithUpdate(&i.batch, batch, ProcessBatch, &i.updatedAt)
}

func (i *ItemUnit) SetMaterial(material uuid.UUID) error {
	return entity.SetWithUpdate(&i.material, material, ProcessMaterial, &i.updatedAt)
}

func (i *ItemUnit) SetQuantity(quantity float64) error {
	return entity.SetWithUpdate(&i.quantity, quantity, ProcessQuantity, &i.updatedAt)
}

func (i *ItemUnit) SetExpiration(expiration time.Time) error {
	return entity.SetWithUpdate(&i.expiration, expiration, ProcessExpiration, &i.updatedAt)
}

func ProcessBatch(batch uuid.UUID) (uuid.UUID, error) {
	if batch == (uuid.UUID{}) {
		return uuid.UUID{}, xerrors.ErrBatchEmpty
	}
	return batch, nil
}

func ProcessMaterial(material uuid.UUID) (uuid.UUID, error) {
	if material == (uuid.UUID{}) {
		return uuid.UUID{}, xerrors.ErrMaterialEmpty
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
