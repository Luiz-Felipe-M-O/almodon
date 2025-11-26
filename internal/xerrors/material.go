package xerrors

import "github.com/alan-b-lima/almodon/pkg/errors"

var (
	ErrMaterialCreation = errors.Imp(errors.InvalidInput, "material-creation", "given data does not satisfy the material type")
	ErrMaterialUpdate   = errors.Imp(errors.InvalidInput, "material-update", "given data does not satisfy the material type")

	ErrInvalidIdLength     = errors.New(errors.InvalidInput, "invalid-id-length", "given id length is not expected", nil)
	ErrIdContainsNonDigits = errors.New(errors.InvalidInput, "non-digit-characters", "given id contains non digit characters", nil)
	ErrNegativeMinQuantity = errors.New(errors.InvalidInput, "negative-min-quantity", "given quantity must be non-negative", nil)
)
