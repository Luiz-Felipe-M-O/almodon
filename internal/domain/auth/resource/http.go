package auths

import (
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/support/resource"
)

type Resource struct {
	http.ServeMux
	Auth auth.Service

	Ident auth.Identifier
}

func New(ident auth.Identifier) http.Handler {
	rc := Resource{
		Ident: ident,
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
	var req auth.Request
	if err := resource.DecodeJSON(&req, r); err != nil {
		resource.WriteError(w, err)
		return
	}

	res, err := rc.Auth.Login(req.SIAPE, req.Password)
	if err != nil {
		resource.WriteError(w, err)
		return
	}

	resource.SetSessionCookie(w, res.UUID, res.Expires)

	if err := resource.EncodeJSON(&res, http.StatusCreated, w, r); err != nil {
		resource.WriteError(w, err)
		return
	}
}

func (rc *Resource) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := resource.SessionCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := rc.Auth.Logout(session); err != nil {
		resource.WriteError(w, err)
		return
	}

	resource.DeleteSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}
