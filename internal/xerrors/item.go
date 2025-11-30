package xerrors

import "github.com/alan-b-lima/almodon/pkg/errors"

var (
	ErrItemCreation = errors.Imp(errors.InvalidInput, "item-creation", "given data does not satisfy the item batch requirements")
	ErrItemUpdate   = errors.Imp(errors.InvalidInput, "item-update", "given data does not satisfy the item batch requirements")

	ErrMaterialNotFound       = errors.New(errors.InvalidInput, "material-not-found", "material UUID cannot be nil", nil)
	ErrSupplierNotFound       = errors.New(errors.InvalidInput, "supplier-not-found", "supplier UUID cannot be nil", nil)
	ErrQuantityMustBePositive = errors.New(errors.InvalidInput, "quantity-must-be-positive", "quantity must be greater than zero", nil)
	ErrUnitCostMustBePositive = errors.New(errors.InvalidInput, "unit-cost-must-be-positive", "unit cost must be greater than zero", nil)
	ErrExpirationInPast       = errors.New(errors.InvalidInput, "expiration-in-past", "expiration date cannot be in the past", nil)
	ErrInvoiceEmpty           = errors.New(errors.InvalidInput, "invoice-empty", "invoice cannot be empty", nil)
	ErrInvoiceTooLong         = errors.New(errors.InvalidInput, "invoice-too-long", "invoice exceeds maximum length of 50 characters", nil)
	ErrLotTooLong             = errors.New(errors.InvalidInput, "lot-too-long", "lot exceeds maximum length of 50 characters", nil)
	ErrNotesTooLong           = errors.New(errors.InvalidInput, "notes-too-long", "notes exceed maximum length of 5000 characters", nil)

	ErrBatchNotFound = errors.New(errors.NotFound, "batch-not-found", "item batch not found", nil)
	ErrItemNotFound  = errors.New(errors.NotFound, "item-not-found", "item not found", nil)
)
