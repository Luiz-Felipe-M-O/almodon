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

// List retrieves all users.
//
//	GET /users/
//
// This function expects an UUID in the path, and no body.
//
// This function returns the following JSON array:
//
//	[
//		{
//			"uuid":    uuid,
//			"siape":   string,
//			"name":    string,
//			"email":   string,
//			"role":    role,
//			"logged":  bool,
//			"created": time,
//			"updated": time,
//		},
//		...
//	]
//
// The fields are as so:
//   - uuid: JSON string that represents a UUID.
//   - siape: JSON string that represents a SIAPE identifier.
//   - name: JSON string that is a name.
//   - email: JSON string that is an email address.
//   - role: JSON string that represents the user's role.
//   - logged: JSON boolean that reports whether the user is logged.
//   - created: JSON string that represents the time of creation of the user, formatted accourding to [RFC 3339].
//   - updated: JSON string that represents the time of the last update of the user, formatted accourding to [RFC 3339].
//
// Role may be one of: "maintainer", "chief", "promoted-admin", "admin" or "user".
//
// This function is limited to inheritors of Chief.
//
// [RFC 3339]: https://tools.ietf.org/html/rfc3339
func (rc *Resource) List(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), rc.Users.List, w, r)
}

// Get retrieves an user by their UUID.
//
//	GET /users/{uuid}
//
// This function expects an UUID in the path, and no body.
//
// This function returns the following JSON object:
//
//	{
//		"uuid":    uuid,
//		"siape":   string,
//		"name":    string,
//		"email":   string,
//		"role":    role,
//		"logged":  bool,
//		"created": time,
//		"updated": time,
//	}
//
// The fields are as so:
//   - uuid: JSON string that represents a UUID.
//   - siape: JSON string that represents a SIAPE identifier.
//   - name: JSON string that is a name.
//   - email: JSON string that is an email address.
//   - role: JSON string that represents the user's role.
//   - logged: JSON boolean that reports whether the user is logged.
//   - created: JSON string that represents the time of creation of the user, formatted accourding to [RFC 3339].
//   - updated: JSON string that represents the time of the last update of the user, formatted accourding to [RFC 3339].
//
// Role may be one of: "maintainer", "chief", "promoted-admin", "admin" or "user".
//
// This function is limited to inheritors of Chief. It also allows users to retrieve themselves.
//
// [RFC 3339]: https://tools.ietf.org/html/rfc3339
func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) (user.Result, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return user.Result{}, resource.ErrBadUUID
		}

		return rc.Users.Get(ctx, uuid)
	}, w, r)
}

// GetBySIAPE retrieves an user by their SIAPE.
//
//	GET /users/siape/{siape}
//
// This function expects a SIAPE in the path, and no body.
//
// This function returns the following JSON object:
//
//	{
//		"uuid":    uuid,
//		"siape":   string,
//		"name":    string,
//		"email":   string,
//		"role":    role,
//		"logged":  bool,
//		"created": time,
//		"updated": time,
//	}
//
// The fields are as so:
//   - uuid: JSON string that represents a UUID.
//   - siape: JSON string that represents a SIAPE identifier.
//   - name: JSON string that is a name.
//   - email: JSON string that is an email address.
//   - role: JSON string that represents the user's role.
//   - logged: JSON boolean that reports whether the user is logged.
//   - created: JSON string that represents the time of creation of the user, formatted accourding to [RFC 3339].
//   - updated: JSON string that represents the time of the last update of the user, formatted accourding to [RFC 3339].
//
// Role may be one of: "maintainer", "chief", "promoted-admin", "admin" or "user".
//
// This function is limited to inheritors of Chief. It also allows users to retrieve themselves.
//
// [RFC 3339]: https://tools.ietf.org/html/rfc3339
func (rc *Resource) GetBySIAPE(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) (user.Result, error) {
		siape := r.PathValue("siape")

		return rc.Users.GetBySIAPE(ctx, siape)
	}, w, r)
}

// Create creates a new user.
//
//	POST /users/
//
// This function expects the following JSON object:
//
//	{
//		"siape":    string,
//		"name":     string,
//		"email":    string,
//		"password": string,
//		"role":     role,
//	}
//
// The fields are as so:
//   - siape: JSON string that represents a SIAPE identifier.
//   - name: JSON string that is a name.
//   - email: JSON string that is an email address.
//   - password: JSON string that is a password.
//   - role: JSON string that represents the user's role.
//
// Role may be one of: "maintainer", "chief", "admin" or "user".
//
// This function returns the following JSON object:
//
//	{
//		"uuid": uuid
//	}
//
// The fields are as so:
//   - uuid: JSON string that represents the newly created user's UUID.
//
// This function is limited to inheritors of Chief.
func (rc *Resource) Create(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(r.Context(), rc.Users.Create, w, r)
}

// Patch updates an user selected by their UUID.
//
//	PATCH /users/{uuid}
//
// This function expects an UUID in the path, and the following JSON object in the body:
//
//	{
//		"name":  string | null,
//		"email": string | null,
//	}
//
// The fields are as so:
//   - name: JSON string that is a name. This field is optional, may be null or not defined.
//   - email: JSON string that is an email address. This field is optional, may be null or not defined.
//
// This function is limited to inheritors of Chief. It also allows users to update themselves.
func (rc *Resource) Patch(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(r.Context(), func(ctx context.Context, req user.Patch) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Users.Patch(ctx, uuid, req)
	}, w, r)
}

// Delete deletes an user selected by their UUID.
//
//	DELETE /users/{uuid}
//
// This function expects an UUID in the path, and no body.
//
// This function is limited to inheritors of Chief. It also allows users to delete themselves.
//
// This function is idempotent, meaning that if the user is already deleted, it will return a success response.
func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(r.Context(), func(ctx context.Context) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Users.Delete(ctx, uuid)
	}, w, r)
}

// Me retrieves the currently logged user.
//
//	GET /users/me/
//
// This function returns the following JSON object:
//
//	{
//		"uuid":    uuid,
//		"siape":   string,
//		"name":    string,
//		"email":   string,
//		"role":    role,
//		"logged":  bool,
//		"created": time,
//		"updated": time,
//	}
//
// The fields are as so:
//   - uuid: JSON string that represents a UUID.
//   - siape: JSON string that represents a SIAPE identifier.
//   - name: JSON string that is a name.
//   - email: JSON string that is an email address.
//   - role: JSON string that represents the user's role.
//   - logged: JSON boolean that reports whether the user is logged.
//   - created: JSON string that represents the time of creation of the user, formatted accourding to [RFC 3339].
//   - updated: JSON string that represents the time of the last update of the user, formatted accourding to [RFC 3339].
//
// Role may be one of: "maintainer", "chief", "promoted-admin", "admin" or "user".
//
// This function is available to all logged users.
//
// [RFC 3339]: https://tools.ietf.org/html/rfc3339
func (rc *Resource) Me(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), rc.Users.Me, w, r)
}
