package materials

import (
	"net/http"

	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/material"
	"github.com/alan-b-lima/almodon/internal/support/resource"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Resource struct {
	http.ServeMux
	Materials *auth.Gatekeeper[material.Service]

	Ident auth.Identifier
}

func New(materials material.Service, ident auth.Identifier) http.Handler {
	rc := Resource{
		Materials: auth.NewGatekeeper(materials),
		Ident:     ident,
	}

	routes := map[string]http.HandlerFunc{
		"GET /materials/{$}":               rc.List,
		"GET /materials/siads/{siads}":     rc.ListBySIADS,
		"GET /materials/catmat/{catmat}":   rc.ListByCATMAT,
		"GET /materials/ecampus/{ecampus}": rc.ListByECAMPUS,
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
	resource.GetHandler(rc.Ident, func(act auth.Actor) (material.ListResult, error) {
		req := material.ListParams{Offset: 0, Limit: 10}
		if err := resource.QueryParams(r.URL.Query(), &req); err != nil {
			return material.ListResult{}, xerrors.ErrBadQueryParams.New(err)
		}

		return rc.Materials.Permit(act).List(req)
	}, w, r)
}

func (rc *Resource) ListBySIADS(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (material.ListResult, error) {
		req := material.ListParams{Offset: 0, Limit: 10}
		if err := resource.QueryParams(r.URL.Query(), &req); err != nil {
			return material.ListResult{}, xerrors.ErrBadQueryParams.New(err)
		}

		ent, err := rc.Materials.Permit(act).ListBySIADS(r.PathValue("siads"), req)
		if err != nil {
			return material.ListResult{}, err
		}

		return ent, nil
	}, w, r)
}

func (rc *Resource) ListByCATMAT(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (material.ListResult, error) {
		req := material.ListParams{Offset: 0, Limit: 10}
		if err := resource.QueryParams(r.URL.Query(), &req); err != nil {
			return material.ListResult{}, xerrors.ErrBadQueryParams.New(err)
		}

		ent, err := rc.Materials.Permit(act).ListByCATMAT(r.PathValue("catmat"), req)
		if err != nil {
			return material.ListResult{}, err
		}

		return ent, nil
	}, w, r)
}

func (rc *Resource) ListByECAMPUS(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (material.ListResult, error) {
		req := material.ListParams{Offset: 0, Limit: 10}
		if err := resource.QueryParams(r.URL.Query(), &req); err != nil {
			return material.ListResult{}, xerrors.ErrBadQueryParams.New(err)
		}

		ent, err := rc.Materials.Permit(act).ListByECampus(r.PathValue("ecampus"), req)
		if err != nil {
			return material.ListResult{}, err
		}

		return ent, nil
	}, w, r)
}

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (material.Result, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return material.Result{}, xerrors.ErrBadUUID
		}

		ent, err := rc.Materials.Permit(act).Get(uuid)
		if err != nil {
			return material.Result{}, err
		}

		return ent, nil
	}, w, r)
}

func (rc *Resource) Create(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(rc.Ident, func(act auth.Actor, req material.Create) (material.CreateResult, error) {
		res, err := rc.Materials.Permit(act).Create(req)
		if err != nil {
			return material.CreateResult{}, err
		}

		return material.CreateResult{UUID: res}, nil
	}, w, r)
}

func (rc *Resource) Patch(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(rc.Ident, func(act auth.Actor, req material.Patch) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return xerrors.ErrBadUUID
		}

		return rc.Materials.Permit(act).Patch(uuid, req)
	}, w, r)
}

func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(rc.Ident, func(act auth.Actor) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return xerrors.ErrBadUUID
		}

		return rc.Materials.Permit(act).Delete(uuid)
	}, w, r)
}
