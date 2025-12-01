package users

import (
	"net/http"
	"time"

	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/support/resource"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Resource struct {
	http.ServeMux
	Users *auth.Gatekeeper[user.Service]

	Ident auth.Identifier
}

func New(users user.Service) http.Handler {
	rc := Resource{
		Users: auth.NewGatekeeper(users),
		Ident: users,
	}

	routes := map[string]http.HandlerFunc{
		"GET /users/{$}":           rc.List,
		"GET /users/{uuid}":        rc.Get,
		"GET /users/siape/{siape}": rc.GetBySIAPE,
		"POST /users/{$}":          rc.Create,
		"PATCH /users/{uuid}":      rc.Patch,
		"DELETE /users/{uuid}":     rc.Delete,
		"POST /users/auth/{$}":     rc.Authenticate,
		"DELETE /users/auth/{$}":   rc.Logout,
		"GET /users/me/{$}":        rc.Me,
		"/":                        resource.NotFound,
	}

	for route, handler := range routes {
		rc.Handle(route, handler)
	}

	return &rc
}

func (rc *Resource) List(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (user.ListResult, error) {
		req := user.ListParams{Offset: 0, Limit: 10}
		if err := resource.QueryParams(r.URL.Query(), &req); err != nil {
			return user.ListResult{}, xerrors.ErrBadQueryParams.New(err)
		}

		ent, err := rc.Users.Permit(act).List(req)
		if err != nil {
			return user.ListResult{}, err
		}

		res := user.ListResult{
			Offset:       ent.Offset,
			Length:       ent.Length,
			Records:      make([]user.Result, len(ent.Records)),
			TotalRecords: ent.TotalRecords,
		}
		for i := range len(ent.Records) {
			res.Records[i] = transform(&ent.Records[i])
		}

		return res, nil
	}, w, r)
}

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (user.Result, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return user.Result{}, xerrors.ErrBadUUID
		}

		ent, err := rc.Users.Permit(act).Get(uuid)
		if err != nil {
			return user.Result{}, err
		}

		return transform(&ent), nil
	}, w, r)
}

func (rc *Resource) GetBySIAPE(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (user.Result, error) {
		ent, err := rc.Users.Permit(act).GetBySIAPE(r.PathValue("siape"))
		if err != nil {
			return user.Result{}, err
		}

		return transform(&ent), nil
	}, w, r)
}

func (rc *Resource) Create(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(rc.Ident, func(act auth.Actor, req user.Create) (user.CreateResult, error) {
		res, err := rc.Users.Permit(act).Create(req)
		if err != nil {
			return user.CreateResult{}, err
		}

		return user.CreateResult{UUID: res}, nil
	}, w, r)
}

func (rc *Resource) Patch(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(rc.Ident, func(act auth.Actor, req user.Patch) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return xerrors.ErrBadUUID
		}

		return rc.Users.Permit(act).Patch(uuid, req)
	}, w, r)
}

func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(rc.Ident, func(act auth.Actor) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return xerrors.ErrBadUUID
		}

		return rc.Users.Permit(act).Delete(uuid)
	}, w, r)
}

func (rc *Resource) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req user.Authenticate
	if err := resource.DecodeJSON(&req, r); err != nil {
		resource.WriteError(w, err)
		return
	}

	res, err := rc.Users.Service.Authenticate(req.SIAPE, req.Password)
	if err != nil {
		resource.WriteError(w, err)
		return
	}

	resource.SetSession(w, res.UUID, res.Expires)

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

	if err := rc.Users.Service.Logout(session); err != nil {
		resource.WriteError(w, err)
		return
	}

	resource.SetSession(w, session, time.Time{})
	w.WriteHeader(http.StatusNoContent)
}

func (rc *Resource) Me(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(rc.Ident, func(act auth.Actor) (user.Result, error) {
		ent, err := rc.Users.Permit(act).Get(act.User())
		if err != nil {
			return user.Result{}, err
		}

		return transform(&ent), nil
	}, w, r)
}

func transform(e *user.Entity) user.Result {
	return user.Result{
		UUID:    e.UUID,
		SIAPE:   e.SIAPE,
		Name:    e.Name,
		Email:   e.Email,
		Role:    e.Role,
		Created: e.Created,
		Updated: e.Updated,
	}
}
