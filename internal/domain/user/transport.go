package user

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type (
	ListParams struct {
		Offset int `query:"offset"`
		Limit  int `query:"limit"`
	}

	Create struct {
		SIAPE    int       `json:"siape"`
		Name     string    `json:"name"`
		Email    string    `json:"email"`
		Password string    `json:"password"`
		Role     auth.Role `json:"role"`
	}

	Patch struct {
		Name  opt.Opt[string] `json:"name"`
		Email opt.Opt[string] `json:"email"`
	}

	UpdatePassword struct {
		Password string `json:"password"`
	}

	UpdateRole struct {
		Role auth.Role `json:"role"`
	}

	Authenticate struct {
		SIAPE    int    `json:"siape"`
		Password string `json:"password"`
	}
)

type (
	ListResult struct {
		Offset       int      `json:"offset"`
		Length       int      `json:"length"`
		Records      []Result `json:"records"`
		TotalRecords int      `json:"total_records"`
	}

	Result struct {
		UUID  uuid.UUID `json:"uuid"`
		SIAPE int       `json:"siape"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Role  auth.Role `json:"role"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}

	AuthResult struct {
		UUID    uuid.UUID `json:"uuid"`
		User    uuid.UUID `json:"user"`
		Expires time.Time `json:"expires"`
	}
)
