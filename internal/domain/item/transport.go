package item

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
		Batch      uuid.UUID `json:"batch"`
		Material   uuid.UUID `json:"material"`
		Quantity   float64   `json:"quantity"`
		Expiration time.Time `json:"expiration"`
	}

	Patch struct {
		Batch      opt.Opt[uuid.UUID] `json:"batch"`
		Material   opt.Opt[uuid.UUID] `json:"material"`
		Quantity   opt.Opt[float64]   `json:"quantity"`
		Expiration opt.Opt[time.Time] `json:"expiration"`
	}

	UpdateQuantity struct {
		Quantity float64 `json:"quantity"`
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
		UUID          uuid.UUID `json:"uuid"`
		Batch         uuid.UUID `json:"batch"`
		Material      uuid.UUID `json:"material"`
		Quantity      float64   `json:"quantity"`
		Expiration    time.Time `json:"expiration,omitzero"`
		IsExpired     bool      `json:"is_expired"`
		HasExpiration bool      `json:"has_expired"`
		Created       time.Time `json:"created"`
		Updated       time.Time `json:"updated"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}
)
