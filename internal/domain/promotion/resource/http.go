package promotions

import (
	"net/http"

	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/promotion"
	"github.com/alan-b-lima/almodon/internal/support/resource"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Resource struct {
	http.ServeMux
	Promotions *auth.Gatekeeper[promotion.Service]

	Ident auth.Identifier
}

func New(promotions promotion.Service, authorizer auth.Identifier) http.Handler {
	rc := Resource{
		Promotions: auth.NewGatekeeper(promotions),
		Ident:      authorizer,
	}

	routes := map[string]http.HandlerFunc{
		"GET /promotions/{$}":       rc.List,
		"GET /promotions/{uuid}":    rc.Get,
		"POST /promotions/{$}":      rc.Create,
		"PUT /promotions/{uuid}":    rc.Update,
		"DELETE /promotions/{uuid}": rc.Delete,
		"/":                         resource.NotFound,
	}

	for route, handler := range routes {
		rc.Handle(route, handler)
	}

	return &rc
}

func (rc *Resource) List(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (promotion.ListResult, error) {
		req := promotion.ListParams{Offset: 0, Limit: 10}
		if err := resource.QueryParams(r.URL.Query(), &req); err != nil {
			return promotion.ListResult{}, xerrors.ErrBadQueryParams.New(err)
		}

		ent, err := rc.Promotions.Permit(act).List(req)
		if err != nil {
			return promotion.ListResult{}, err
		}

		res := promotion.ListResult{
			Offset:       ent.Offset,
			Length:       ent.Length,
			Records:      make([]promotion.Result, len(ent.Records)),
			TotalRecords: ent.TotalRecords,
		}
		for i := range len(ent.Records) {
			transpose(&res.Records[i], &ent.Records[i])
		}

		return res, nil
	}, w, r)
}

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (promotion.Result, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return promotion.Result{}, xerrors.ErrBadUUID
		}

		ent, err := rc.Promotions.Permit(act).Get(uuid)
		if err != nil {
			return promotion.Result{}, err
		}

		return transform(&ent), nil
	}, w, r)
}

func (rc *Resource) Create(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(rc.Ident, func(act auth.Actor, req promotion.Create) (promotion.CreateResult, error) {
		res, err := rc.Promotions.Permit(act).Create(req)
		if err != nil {
			return promotion.CreateResult{}, err
		}

		return promotion.CreateResult{UUID: res}, nil
	}, w, r)
}

func (rc *Resource) Update(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(rc.Ident, func(act auth.Actor, req promotion.Update) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return xerrors.ErrBadUUID
		}

		return rc.Promotions.Permit(act).Update(uuid, req)
	}, w, r)
}

func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(rc.Ident, func(act auth.Actor) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return xerrors.ErrBadUUID
		}

		return rc.Promotions.Permit(act).Delete(uuid)
	}, w, r)
}

func transform(e *promotion.Entity) promotion.Result {
	return promotion.Result(*e)
}

func transpose(r *promotion.Result, e *promotion.Entity) {
	*r = promotion.Result(*e)
}
