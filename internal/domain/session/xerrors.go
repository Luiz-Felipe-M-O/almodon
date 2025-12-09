package session

import "github.com/alan-b-lima/almodon/pkg/errors"

var (
	ErrSessionTooLong  = errors.New(errors.InvalidInput, "session-too-long", "session too long", nil, map[string]any{"max": MaxAgeMax})

	ErrSessionNotFound = errors.New(errors.NotFound, "session-not-found", "session not found", nil, nil)
)
