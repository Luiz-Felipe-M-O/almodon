package xerrors

import "github.com/alan-b-lima/almodon/pkg/errors"

var (
	ErrRequisitionCreation = errors.Imp(errors.InvalidInput, "requisition-creation", "given data does not satisfy requisition requirements")
	ErrRequisitionNotFound = errors.New(errors.NotFound, "requisition-not-found", "requisition not found", nil)

	ErrAuthorEmpty          = errors.New(errors.InvalidInput, "author-empty", "author cannot be empty", nil)
	ErrDestinationEmpty     = errors.New(errors.InvalidInput, "destination-empty", "destination cannot be empty", nil)
	ErrDestinationTooLong   = errors.Fmt(errors.InvalidInput, "destination-too-long", "destination exceeds maximum length of %d characters")
	ErrNotesTooLong         = errors.Fmt(errors.InvalidInput, "notes-too-long", "notes exceed maximum length of %d characters")
	ErrApprovalNotesTooLong = errors.Fmt(errors.InvalidInput, "approval-notes-too-long", "approval notes exceed maximum length of %d characters")

	ErrMaterialEmpty           = errors.New(errors.InvalidInput, "material-empty", "material cannot be empty", nil)
	ErrQuantityInvalid         = errors.New(errors.InvalidInput, "quantity-invalid", "quantity must be greater than zero", nil)
	ErrApprovedQuantityInvalid = errors.New(errors.InvalidInput, "approved-quantity-invalid", "approved quantity must be non-negative", nil)

	ErrEntryNotFound = errors.New(errors.NotFound, "entry-not-found", "requisition entry not found", nil)
	ErrInvalidStatus = errors.New(errors.InvalidInput, "invalid-status", "invalid requisition status", nil)

	ErrCannotModifyAnswered = errors.New(errors.InvalidInput, "cannot-modify-answered", "cannot modify requisition that has been answered", nil)
	ErrCannotAnswerTwice    = errors.New(errors.InvalidInput, "cannot-answer-twice", "requisition has already been answered", nil)

	ErrRequisitionUpdate          = errors.Imp(errors.InvalidInput, "requisition-update", "given data does not satisfy requisition requirements")
	ErrRequisitionMustHaveEntries = errors.New(errors.InvalidInput, "requisition-must-have-entries", "requisition must have at least one entry", nil)
	ErrCannotFulfillUnapproved    = errors.New(errors.InvalidInput, "cannot-fulfill-unapproved", "can only fulfill approved requisitions", nil)
)
