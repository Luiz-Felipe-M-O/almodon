package promotion

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrTooLong  = problem.New(problem.SemanticalError, "promotion-too-long", "promotion too long", nil, map[string]any{"max": MaxAgeMax})
	ErrNotFound = problem.New(problem.NotFound, "promotion-not-found", "promotion not found", nil, nil)
)
