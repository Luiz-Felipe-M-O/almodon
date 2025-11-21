package promotion

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type (
	ListParams struct {
		Offset int `query:"offset"`
		Limit  int `query:"limit"`
	}

	Create struct {
		User   uuid.UUID              `json:"user"`
		MaxAge opt.Opt[time.Duration] `json:"max_age"`
	}

	Update struct {
		UUID   uuid.UUID              `json:"user"`
		MaxAge opt.Opt[time.Duration] `json:"max_age"`
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
		User    uuid.UUID `json:"user"`
		Expires time.Time `json:"expires"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}
)
