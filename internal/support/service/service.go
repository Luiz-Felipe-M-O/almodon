package service

import (
	"context"
	"errors"

	"github.com/alan-b-lima/almodon/internal/domain/auth"

	"github.com/alan-b-lima/almodon/pkg/rbac"
	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/pkg/problem"
)

func AuthorizeFromContext(ctx context.Context, gate auth.Authenticator, perms rbac.Permission[auth.Role]) (auth.Actor, error) {
	actor, err := ActorFromContext(ctx, gate)
	if err != nil {
		return auth.NewUnlogged(), err
	}

	if err := Authorize(perms, actor); err != nil {
		if actor.Role == auth.Unlogged {
			return actor, auth.ErrUnauthenticated.Make()
		}

		return actor, err
	}

	return actor, nil
}

func ActorFromContext(ctx context.Context, gate auth.Authenticator) (auth.Actor, error) {
	session, ok := ctx.Value("session").(uuid.UUID)
	if !ok {
		return auth.NewUnlogged(), nil
	}

	actor, err := gate.Actor(ctx, session)
	if err != nil {
		if err, ok := errors.AsType[*problem.Error](err); ok && err.IsExternal() {
			return auth.NewUnlogged(), auth.ErrUnauthenticated.Cause(err).Make()
		}

		return auth.NewUnlogged(), err
	}

	return actor, nil
}

func Authorize(perms rbac.Permission[auth.Role], actor auth.Actor) error {
	if role := actor.Role; !perms.Authorize(role) {
		return auth.ErrUnauthorized.Details(map[string]any{"allowed": perms}).Make(role, perms)
	}

	return nil
}
