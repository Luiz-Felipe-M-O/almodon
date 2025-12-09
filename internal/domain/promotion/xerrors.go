package promotion

import "github.com/alan-b-lima/almodon/pkg/errors"

var (
	ErrPromotionTooLong  = errors.New(errors.InvalidInput, "promotion-too-long", "promotion too long", nil, map[string]any{"max": MaxAgeMax})

	ErrPromotionNotFound = errors.New(errors.NotFound, "promotion-not-found", "promotion not found", nil, nil)
)
