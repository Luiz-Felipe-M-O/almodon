package material

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrCreate   = problem.Imp(problem.SemanticalError, "material-create").Message("could not create material")
	ErrUpdate   = problem.Imp(problem.SemanticalError, "material-update").Message("could not update material")
	ErrNotFound = problem.New(problem.NotFound, "material-not-found", "material not found", nil, nil)

	ErrNameEmpty   = problem.New(problem.SemanticalError, "name-empty", "name must not be empty", nil, nil)
	ErrNameTooLong = problem.Imp(problem.SemanticalError, "name-too-long").Format("name must be less than %d characters").Details(map[string]any{"max": NameMaxLen}).Make(NameMaxLen)

	ErrCATMATInvalid = problem.New(problem.SemanticalError, "catmat-invalid", "catmat must be a 6-digit number", nil, nil)

	ErrDescriptionTooLong = problem.Imp(problem.SemanticalError, "description-too-long").Format("description must be less than %d characters").Details(map[string]any{"max": DescriptionMaxLen}).Make(DescriptionMaxLen)

	ErrMinNegative = problem.New(problem.SemanticalError, "min-negative", "min quantity must not be negative", nil, nil)
)
