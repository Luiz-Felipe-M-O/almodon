package items

import (
	"net/http"

	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/item"
	"github.com/alan-b-lima/almodon/internal/support/resource"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Resource struct {
	http.ServeMux
	Items *auth.Gatekeeper[item.Service]

	Ident auth.Identifier
}

func New(items item.Service, ident auth.Identifier) http.Handler {
	rc := Resource{
		Items: auth.NewGatekeeper(items),
		Ident: ident,
	}

	routes := map[string]http.HandlerFunc{
		"GET /items/{$}":                   rc.List,
		"GET /items/{uuid}":                rc.Get,
		"GET /items/material/{uuid}":       rc.ListByMaterial,
		"GET /items/supplier/{uuid}":       rc.ListBySupplier,
		"POST /items/{$}":                  rc.Create,
		"PATCH /items/{uuid}":              rc.Patch,
		"PATCH /items/{uuid}/quantity/{$}": rc.UpdateQuantity,
		"DELETE /items/{uuid}":             rc.Delete,
		"/":                                resource.NotFound,
	}

	for route, handler := range routes {
		rc.Handle(route, handler)
	}

	return &rc
}

func (rc *Resource) List(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (item.ListResult, error) {
		req := item.ListParams{Offset: 0, Limit: 10}
		if err := resource.QueryParams(r.URL.Query(), &req); err != nil {
			return item.ListResult{}, xerrors.ErrBadQueryParams.New(err)
		}

		ent, err := rc.Items.Permit(act).List(req)
		if err != nil {
			return item.ListResult{}, err
		}

		res := item.ListResult{
			Offset:       ent.Offset,
			Length:       ent.Length,
			Records:      make([]item.Result, len(ent.Records)),
			TotalRecords: ent.TotalRecords,
		}
		for i := range len(ent.Records) {
			res.Records[i] = transform(&ent.Records[i])
		}

		return res, nil
	}, w, r)
}

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (item.Result, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return item.Result{}, xerrors.ErrBadUUID
		}

		ent, err := rc.Items.Permit(act).Get(uuid)
		if err != nil {
			return item.Result{}, err
		}

		return transform(&ent), nil
	}, w, r)
}

func (rc *Resource) ListByMaterial(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (item.ListResult, error) {
		material, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return item.ListResult{}, xerrors.ErrBadUUID
		}

		ent, err := rc.Items.Permit(act).ListByMaterial(material)
		if err != nil {
			return item.ListResult{}, err
		}

		res := item.ListResult{
			Offset:       ent.Offset,
			Length:       ent.Length,
			Records:      make([]item.Result, len(ent.Records)),
			TotalRecords: ent.TotalRecords,
		}
		for i := range len(ent.Records) {
			res.Records[i] = transform(&ent.Records[i])
		}

		return res, nil
	}, w, r)
}

func (rc *Resource) ListBySupplier(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (item.ListResult, error) {
		supplier, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return item.ListResult{}, xerrors.ErrBadUUID
		}

		ent, err := rc.Items.Permit(act).ListBySupplier(supplier)
		if err != nil {
			return item.ListResult{}, err
		}

		res := item.ListResult{
			Offset:       ent.Offset,
			Length:       ent.Length,
			Records:      make([]item.Result, len(ent.Records)),
			TotalRecords: ent.TotalRecords,
		}
		for i := range len(ent.Records) {
			res.Records[i] = transform(&ent.Records[i])
		}

		return res, nil
	}, w, r)
}

func (rc *Resource) Create(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(rc.Ident, func(act auth.Actor, req item.Create) (item.CreateResult, error) {
		res, err := rc.Items.Permit(act).Create(req)
		if err != nil {
			return item.CreateResult{}, err
		}

		return item.CreateResult{UUID: res}, nil
	}, w, r)
}

func (rc *Resource) Patch(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(rc.Ident, func(act auth.Actor, req item.Patch) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return xerrors.ErrBadUUID
		}

		return rc.Items.Permit(act).Patch(uuid, req)
	}, w, r)
}

func (rc *Resource) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(rc.Ident, func(act auth.Actor, req item.UpdateQuantity) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return xerrors.ErrBadUUID
		}

		return rc.Items.Permit(act).UpdateQuantity(uuid, req)
	}, w, r)
}

func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(rc.Ident, func(act auth.Actor) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return xerrors.ErrBadUUID
		}

		return rc.Items.Permit(act).Delete(uuid)
	}, w, r)
}

func transform(e *item.Entity) item.Result {
	return item.Result{
		UUID:          e.UUID,
		Material:      e.Material,
		Supplier:      e.Supplier,
		Quantity:      e.Quantity,
		UnitCost:      e.UnitCost,
		Arrival:       e.Arrival,
		Expiration:    e.Expiration,
		Invoice:       e.Invoice,
		Lot:           e.Lot,
		Notes:         e.Notes,
		IsExpired:     item.IsExpired(e.Expiration),
		HasExpiration: item.HasExpiration(e.Expiration),
		IsAvailable:   e.Quantity > 0 && !item.IsExpired(e.Expiration),
		Created:       e.Created,
		Updated:       e.Updated,
	}
}
