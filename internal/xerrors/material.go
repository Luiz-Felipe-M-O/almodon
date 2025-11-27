package xerrors

import "github.com/alan-b-lima/almodon/pkg/errors"

var (
	ErrMaterialCreation = errors.Imp(errors.InvalidInput, "material-creation", "given data does not satisfy the material type")
	ErrMaterialUpdate   = errors.Imp(errors.InvalidInput, "material-update", "given data does not satisfy the material type")

	ErrUnitNotFound        = errors.New(errors.InvalidInput, "unit-not-found", "given unit was not found in allowed units", nil)
	ErrCATMATInvalid       = errors.Imp(errors.InvalidInput, "catmat-invalid", "given CATMAT is invalid")
	ErrSIADSInvalid        = errors.Imp(errors.InvalidInput, "siads-invalid", "given SIADS is invalid")
	ErrECampusInvalid      = errors.Imp(errors.InvalidInput, "ecampus-invalid", "given eCampus code is invalid")
	ErrDescriptionTooLong  = errors.New(errors.InvalidInput, "description-too-long", "given description exceeds maximum length", nil)
	ErrUnitEmpty           = errors.New(errors.InvalidInput, "unit-empty", "given unit is empty", nil)
	ErrMinQuantityNegative = errors.New(errors.InvalidInput, "min-quantity-negative", "given quantity must be non-negative", nil)

	ErrInvalidIdLength      = errors.Fmt(errors.InvalidInput, "id-length-invalid", "given id does not match expected length of %d")
	ErrInvalidIdLengthRange = errors.Fmt(errors.InvalidInput, "id-length-range-invalid", "given id is not within expected %d-%d length")
	ErrIdContainsNonDigits  = errors.New(errors.InvalidInput, "non-digit-chars", "given id contains non digit characters", nil)

	ErrMaterialNotFound = errors.New(errors.NotFound, "material-not-found", "material not found", nil)
)
