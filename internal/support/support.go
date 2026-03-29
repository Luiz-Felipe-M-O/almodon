package support

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrNilPointer = problem.Imp(problem.NotFound, "nil-pointer").Message("nil pointer")
	ErrTODO       = problem.New(problem.Unimplemented, "todo", "implement me", nil, nil)
)
