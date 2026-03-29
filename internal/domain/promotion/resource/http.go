package promotions

import (
	"context"
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain/promotion"

	"github.com/alan-b-lima/almodon/internal/support/resource"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Resource struct {
	http.ServeMux

	Promotions promotion.Service
}

func New(promotions promotion.Service) *Resource {
	rc := Resource{
		Promotions: promotions,
	}

	routes := map[string]http.HandlerFunc{
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

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) (promotion.Result, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return promotion.Result{}, resource.ErrBadUUID
		}

		ent, err := rc.Promotions.Get(ctx, uuid)
		if err != nil {
			return promotion.Result{}, err
		}

		return promotion.Result(ent), nil
	}, w, r)
}

func (rc *Resource) Create(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(r.Context(), func(ctx context.Context, req promotion.Create) (promotion.CreateResult, error) {
		res, err := rc.Promotions.Create(ctx, req)
		if err != nil {
			return promotion.CreateResult{}, err
		}

		return res, nil
	}, w, r)
}

func (rc *Resource) Update(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(r.Context(), func(ctx context.Context, req promotion.Update) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Promotions.Update(ctx, uuid, req)
	}, w, r)
}

func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(r.Context(), func(ctx context.Context) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Promotions.Delete(ctx, uuid)
	}, w, r)
}
