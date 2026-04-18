package materials

import (
	"context"
	"net/http"
	"strconv"

	"github.com/alan-b-lima/almodon/internal/domain/material"
	"github.com/alan-b-lima/almodon/internal/support/resource"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Resource struct {
	http.ServeMux

	Materials material.Service
}

func New(materials material.Service) *Resource {
	rc := Resource{
		Materials: materials,
	}

	routes := map[string]http.HandlerFunc{
		"GET /materials/{$}":               rc.List,
		"GET /materials/ecampus/{ecampus}": rc.ListByECampus,
		"GET /materials/catmat/{catmat}":   rc.ListByCATMAT,
		"GET /materials/siads/{siads}":     rc.ListBySIADS,
		"GET /materials/{uuid}":            rc.Get,
		"POST /materials/{$}":              rc.Create,
		"PATCH /materials/{uuid}":          rc.Patch,
		"DELETE /materials/{uuid}":         rc.Delete,
		"/":                                resource.NotFound,
	}

	for route, handler := range routes {
		rc.Handle(route, handler)
	}

	return &rc
}

func (rc *Resource) List(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) ([]material.Result, error) {
		return rc.Materials.List(ctx)
	}, w, r)
}

func (rc *Resource) ListByECampus(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) ([]material.Result, error) {
		ecampus, err := strconv.Atoi(r.PathValue("ecampus"))
		if err != nil {
			return nil, resource.ErrBadUUID
		}

		ent, err := rc.Materials.ListByECampus(ctx, ecampus)
		if err != nil {
			return nil, err
		}

		return ent, nil
	}, w, r)
}

func (rc *Resource) ListByCATMAT(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) ([]material.Result, error) {
		catmat, err := strconv.Atoi(r.PathValue("catmat"))
		if err != nil {
			return nil, resource.ErrBadUUID
		}

		ent, err := rc.Materials.ListByCATMAT(ctx, catmat)
		if err != nil {
			return nil, err
		}

		return ent, nil
	}, w, r)
}

func (rc *Resource) ListBySIADS(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) ([]material.Result, error) {
		siads, err := strconv.Atoi(r.PathValue("siads"))
		if err != nil {
			return nil, resource.ErrBadUUID
		}

		ent, err := rc.Materials.ListBySIADS(ctx, siads)
		if err != nil {
			return nil, err
		}

		return ent, nil
	}, w, r)
}

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) (material.Result, error) {
		uuid, res := uuid.FromString(r.PathValue("uuid"))
		if res != nil {
			return material.Result{}, resource.ErrBadUUID
		}

		ent, res := rc.Materials.Get(ctx, uuid)
		if res != nil {
			return material.Result{}, res
		}

		return ent, nil
	}, w, r)
}

func (rc *Resource) Create(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(r.Context(), func(ctx context.Context, req material.Create) (material.CreateResult, error) {
		res, err := rc.Materials.Create(ctx, req)
		if err != nil {
			return material.CreateResult{}, err
		}

		return res, nil
	}, w, r)
}

func (rc *Resource) Patch(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(r.Context(), func(ctx context.Context, req material.Patch) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Materials.Patch(ctx, uuid, req)
	}, w, r)
}

func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(r.Context(), func(ctx context.Context) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Materials.Delete(ctx, uuid)
	}, w, r)
}
