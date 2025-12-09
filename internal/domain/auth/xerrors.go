package auth

import (
	"github.com/alan-b-lima/almodon/pkg/errors"
)

var (
	ErrUnauthenticated = errors.Imp(errors.Unauthentic, "unauthenticated").Message("unauthenticated user")
	ErrUnauthorized    = errors.Imp(errors.Forbidden, "unauthorized").Format("actor role %v does not inherit any roles in %v")
)
