package auths

import (
	"context"
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/support/resource"
	"github.com/alan-b-lima/almodon/internal/support/session"
)

type Resource struct {
	http.ServeMux
	Auth auth.Service
}

func New(auth auth.Service) *Resource {
	rc := Resource{
		Auth: auth,
	}

	routes := map[string]http.HandlerFunc{
		"POST /auth/{$}":   rc.Login,
		"DELETE /auth/{$}": rc.Logout,
		"/":                resource.NotFound,
	}

	for route, handler := range routes {
		rc.Handle(route, handler)
	}

	return &rc
}

func (rc *Resource) Login(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(r.Context(), func(ctx context.Context, req auth.Create) (auth.Result, error) {
		res, err := rc.Auth.Login(r.Context(), req.SIAPE, req.Password)
		if err != nil {
			return auth.Result{}, err
		}

		session.SetCookie(w, res.UUID, res.Expires)
		return res, nil
	}, w, r)
}

func (rc *Resource) Logout(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(r.Context(), func(ctx context.Context) error {
		s, err := session.Cookie(r)
		if err != nil {
			return nil
		}

		if err := rc.Auth.Logout(r.Context(), s); err != nil {
			return nil
		}

		session.DeleteCookie(w)
		return nil
	}, w, r)
}
