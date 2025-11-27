package material

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
		Name        string  `json:"name"`
		SIADS       string  `json:"siads"`
		CATMAT      string  `json:"catmat"`
		ECampus     string  `json:"ecampus"`
		Description string  `json:"description"`
		Unit        string  `json:"unit"`
		MinQuantity float64 `json:"min_quantity"`
	}

	Patch struct {
		Name        opt.Opt[string]  `json:"name"`
		SIADS       opt.Opt[string]  `json:"siads"`
		CATMAT      opt.Opt[string]  `json:"catmat"`
		ECampus     opt.Opt[string]  `json:"ecampus"`
		Description opt.Opt[string]  `json:"description"`
		Unit        opt.Opt[string]  `json:"unit"`
		MinQuantity opt.Opt[float64] `json:"min_quantity"`
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
		UUID        uuid.UUID `json:"uuid"`
		Name        string    `json:"name"`
		SIADS       string    `json:"siads"`
		CATMAT      string    `json:"catmat"`
		ECampus     string    `json:"ecampus"`
		Description string    `json:"description"`
		Unit        string    `json:"unit"`
		MinQuantity float64   `json:"min_quantity"`
		Created     time.Time `json:"created"`
		Updated     time.Time `json:"updated"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}
)
