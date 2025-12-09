package user

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type (
	ListParams struct {
		Offset int `query:"offset"`
		Limit  int `query:"limit"`
	}

	Create struct {
		SIAPE    string    `json:"siape"`
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
)

type (
	ListResult struct {
		Offset       int      `json:"offset"`
		Length       int      `json:"length"`
		Records      []Result `json:"records"`
		TotalRecords int      `json:"total_records"`
	}

	Result struct {
		UUID    uuid.UUID `json:"uuid"`
		SIAPE   string    `json:"siape"`
		Name    string    `json:"name"`
		Email   string    `json:"email"`
		Role    auth.Role `json:"role"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}
)
