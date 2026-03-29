package users

import (
	"context"
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain/user"

	"github.com/alan-b-lima/almodon/internal/support/resource"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Resource struct {
	http.ServeMux

	Users user.Service
}

func New(users user.Service) *Resource {
	rc := Resource{
		Users: users,
	}

	routes := map[string]http.HandlerFunc{
		"GET /users/{$}":           rc.List,
		"GET /users/{uuid}":        rc.Get,
		"GET /users/siape/{siape}": rc.GetBySIAPE,
		"POST /users/{$}":          rc.Create,
		"PATCH /users/{uuid}":      rc.Patch,
		"DELETE /users/{uuid}":     rc.Delete,
		"GET /users/me/{$}":        rc.Me,
		"/":                        resource.NotFound,
	}

	for route, handler := range routes {
		rc.Handle(route, handler)
	}

	return &rc
}

func (rc *Resource) List(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) ([]user.Result, error) {
		return rc.Users.List(ctx)
	}, w, r)
}

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) (user.Result, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return user.Result{}, resource.ErrBadUUID
		}

		ent, err := rc.Users.Get(ctx, uuid)
		if err != nil {
			return user.Result{}, err
		}

		return ent, nil
	}, w, r)
}

func (rc *Resource) GetBySIAPE(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) (user.Result, error) {
		siape := r.PathValue("siape")

		ent, err := rc.Users.GetBySIAPE(ctx, siape)
		if err != nil {
			return user.Result{}, err
		}

		return ent, nil
	}, w, r)
}

func (rc *Resource) Create(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(r.Context(), func(ctx context.Context, req user.Create) (user.CreateResult, error) {
		res, err := rc.Users.Create(ctx, req)
		if err != nil {
			return user.CreateResult{}, err
		}

		return res, nil
	}, w, r)
}

func (rc *Resource) Patch(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(r.Context(), func(ctx context.Context, req user.Patch) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Users.Patch(ctx, uuid, req)
	}, w, r)
}

func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(r.Context(), func(ctx context.Context) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Users.Delete(ctx, uuid)
	}, w, r)
}

func (rc *Resource) Me(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) (user.Result, error) {
		return rc.Users.Me(ctx)
	}, w, r)
}
