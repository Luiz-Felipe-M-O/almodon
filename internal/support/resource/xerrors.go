package resource

import "github.com/alan-b-lima/almodon/pkg/errors"

var (
	ErrBadUUID        = errors.New(errors.InvalidInput, "bad-uuid", "invalid uuid", nil, nil)
	ErrBadQueryParams = errors.Imp(errors.InvalidInput, "bad-uuid").Message("invalid query params")

	ErrNoContentType          = errors.New(errors.PreconditionFailed, "no-content-type", "content type must be informed", nil, nil)
	ErrUnsupportedContentType = errors.Imp(errors.PreconditionFailed, "unsupported-content-type").Format("content type must be %s")
	ErrNotAcceptableJson      = errors.Imp(errors.PreconditionFailed, "unacceptable-type").Format("client does not accept %s")

	ErrResourceNotFound = errors.Imp(errors.NotFound, "resource-not-found").Format("resource %+q not found")
	ErrNotAcceptable    = errors.Imp(errors.PreconditionFailed, "not-acceptable").Format("client does not accept %s")
	ErrJSON             = errors.Imp(errors.InvalidInput, "json-error").Message("unexpected error ocurred while processing")
)
