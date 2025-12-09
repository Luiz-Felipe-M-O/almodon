package service

import (
	authpkg "github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/pkg/auth"
)

func Authorize(auth auth.Permission[authpkg.Role], actor authpkg.Actor) error {
	if role := actor.Role(); !auth.Authorize(role) {
		return authpkg.ErrUnauthorized.Metadata(map[string]any{"allowed": auth}).Make(role, auth)
	}

	return nil
}
