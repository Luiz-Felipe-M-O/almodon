package xerrors

import "github.com/alan-b-lima/almodon/pkg/errors"

var (
	ErrItemCreation  = errors.Imp(errors.InvalidInput, "material-creation", "given data does not satisfy the material type")
	ErrItemUpdate    = errors.Imp(errors.InvalidInput, "material-update", "given data does not satisfy the material type")
	
	ErrQuantityMustBePositive = errors.New(errors.InvalidInput, "quantity-must-be-positive", "given quantity must be positive", nil)
	ErrBatchNotFound = errors.New(errors.NotFound, "batch-not-found", "given batch was not found", nil)

	ErrItemNotFound = errors.New(errors.NotFound, "item-not-found", "item not found", nil)
)
