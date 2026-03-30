package promotion

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrCreate   = problem.Imp(problem.SemanticalError, "promotion-create").Message("could not create promotion")
	ErrNotFound = problem.New(problem.NotFound, "promotion-not-found", "promotion not found", nil, nil)

	ErrTooLong = problem.New(problem.SemanticalError, "session-too-long", "session too long", nil, map[string]any{"max": MaxAgeMax})
)
