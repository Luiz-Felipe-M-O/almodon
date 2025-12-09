package support

import "github.com/alan-b-lima/almodon/pkg/errors"

var ErrTODO = errors.New(errors.Unimplemented, "todo", "implement me", nil, nil)
