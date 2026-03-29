package support

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrNilPointer = problem.New(problem.NotFound, "nil-pointer", "nil pointer", nil, nil)
	ErrTODO       = problem.New(problem.Unimplemented, "todo", "implement me", nil, nil)
)
