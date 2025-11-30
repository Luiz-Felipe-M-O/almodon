package item

import (
	"strings"
	"time"

	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

const (
	invoiceMaxLength = 50
	lotMaxLength     = 50
	notesMaxLength = 5000
)

type ItemBatch struct {
	uuid       uuid.UUID
	material   uuid.UUID
	supplier   uuid.UUID
	quantity   float64
	unitCost   float64
	arrival    time.Time
	expiration time.Time
	invoice    string
	lot        string
	notes      string
}

func New(
	material, supplier uuid.UUID,
	quantity, unitCost float64,
	arrival expiration time.Time,
	invoice, lot, notes string,
) (ItemBatch, error) {
	var i ItemBatch

	err := errors.Join(
		i.SetMaterial(material),
		i.SetSupplier(supplier),
		i.SetQuantity(quantity),
		i.SetExpiration(expiration),
		i.SetArrival(arrival),
		i.SetInvoice(invoice),
		i.SetLot(lot),
		i.SetNotes(notes),
	)
	if err != nil {
		return ItemBatch{}, xerrors.ErrItemCreation.New(err)
	}

	i.uuid = uuid.NewUUIDv7()
	return i, nil
}

func (i *ItemBatch) UUID() uuid.UUID       { return i.uuid }
func (i *ItemBatch) Material() uuid.UUID   { return i.material }
func (i *ItemBatch) Supplier() uuid.UUID   { return i.supplier }
func (i *ItemBatch) Quantity() float64     { return i.quantity }
func (i *ItemBatch) UnitCost() float64     { return i.unitCost }
func (i *ItemBatch) Expiration() time.Time { return i.expiration }
func (i *ItemBatch) Arrival() time.Time    { return i.arrival }
func (i *ItemBatch) Invoice() string       { return i.invoice }
func (i *ItemBatch) Lot() string           { return i.lot }
func (i *ItemBatch) Notes() string         { return i.notes }

func IsExpired(expiration time.Time) bool {
	if expiration.IsZero() {
		return false
	}
	return time.Now().After(expiration)
}

func HasExpiration(expiration time.Time) bool {
	return !expiration.IsZero()
}

func (i *ItemBatch) IsExpired() bool     { return IsExpired(i.expiration) }
func (i *ItemBatch) HasExpiration() bool { return HasExpiration(i.expiration) }
func (i *ItemBatch) IsAvailable() bool   { return i.quantity > 0 && !i.IsExpired() }

func (i *ItemBatch) SetMaterial(material uuid.UUID) error {
	return entity.Set(&i.material, material, ProcessMaterial)
}

func (i *ItemBatch) SetSupplier(supplier uuid.UUID) error {
	return entity.Set(&i.supplier, supplier, ProcessSupplier)
}

func (i *ItemBatch) SetQuantity(quantity float64) error {
	return entity.Set(&i.quantity, quantity, ProcessQuantity)
}

func (i *ItemBatch) SetExpiration(expiration time.Time) error {
	return entity.Set(&i.expiration, expiration, ProcessExpiration)
}

func (i *ItemBatch) SetArrival(arrival time.Time) error {
	return entity.Set(&i.arrival, arrival, ProcessArrival)
}

func (i *ItemBatch) SetInvoice(invoice string) error {
	return entity.Set(&i.invoice, invoice, ProcessInvoice)
}

func (i *ItemBatch) SetLot(lot string) error {
	return entity.Set(&i.lot, lot, ProcessLot)
}

func (i *ItemBatch) SetNotes(notes string) error {
	return entity.Set(&i.notes, notes, ProcessNotes)
}

func ProcessMaterial(material uuid.UUID) (uuid.UUID, error) {
	if material.IsNil() {
		return uuid.UUID{}, xerrors.ErrMaterialNotFound
	}
	return material, nil
}

func ProcessSupplier(supplier uuid.UUID) (uuid.UUID, error) {
	if supplier.IsNil() {
		return uuid.UUID{}, xerrors.ErrSupplierNotFound
	}
	return supplier, nil
}

func ProcessQuantity(quantity float64) (float64, error) {
	if quantity <= 0 {
		return 0, xerrors.ErrQuantityMustBePositive
	}
	return quantity, nil
}

func ProcessUnitCost(unitCost float64) (float64, error) {
	if unitCost <= 0 {
		return 0, xerrors.ErrUnitCostMustBePositive
	}
	return unitCost, nil
}

func ProcessArrival(arrival time.Time) (time.Time, error) {
	return arrival, nil
}

func ProcessExpiration(expiration time.Time) (time.Time, error) {
    // Zero means "no expiration"
    if expiration.IsZero() {
        return time.Time{}, nil
    }

    if time.Now().After(expiration) {
        return time.Time{}, xerrors.ErrExpirationInPast
    }

    return expiration, nil
}

func ProcessInvoice(invoice string) (string, error) {
	invoice = strings.TrimSpace(invoice)
	if invoice == "" {
		return "", xerrors.ErrInvoiceEmpty
	}
	if len(invoice) > invoiceMaxLength {
		return "", xerrors.ErrInvoiceTooLong
	}
	return invoice, nil
}

func ProcessLot(lot string) (string, error) {
	lot = strings.TrimSpace(lot)
	if lot == "" {
		return "", nil
	}

	if len(lot) > lotMaxLength {
		return "", xerrors.ErrLotTooLong
	}

	return lot, nil
}

func ProcessNotes(notes string) (string, error) {
	if notes == "" {
		return "", nil
	}

	if len(notes) > notesMaxLength {
		return "", xerrors.ErrLotTooLong
	}

	return notes, nil
}
