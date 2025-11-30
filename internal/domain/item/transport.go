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
		Material   uuid.UUID `json:"material"`
		Supplier   uuid.UUID `json:"supplier"`
		Quantity   float64   `json:"quantity"`
		UnitCost   float64   `json:"unit_cost"`
		Arrival    time.Time `json:"arrival"`
		Expiration time.Time `json:"expiration"`
		Invoice    string    `json:"invoice"`
		Lot        string    `json:"lot"`
		Notes      string    `json:"notes"`
	}

	Patch struct {
		Material   opt.Opt[uuid.UUID] `json:"material"`
		Supplier   opt.Opt[uuid.UUID] `json:"supplier"`
		Quantity   opt.Opt[float64]   `json:"quantity"`
		UnitCost   opt.Opt[float64]   `json:"unit_cost"`
		Arrival    opt.Opt[time.Time] `json:"arrival"`
		Expiration opt.Opt[time.Time] `json:"expiration"`
		Invoice    opt.Opt[string]    `json:"invoice"`
		Lot        opt.Opt[string]    `json:"lot"`
		Notes      opt.Opt[string]    `json:"notes"`
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
		Material      uuid.UUID `json:"material"`
		Supplier      uuid.UUID `json:"supplier"`
		Quantity      float64   `json:"quantity"`
		UnitCost      float64   `json:"unit_cost"`
		Arrival       time.Time `json:"arrival"`
		Expiration    time.Time `json:"expiration,omitzero"`
		Invoice       string    `json:"invoice"`
		Lot           string    `json:"lot,omitempty"`
		Notes         string    `json:"notes,omitempty"`
		IsExpired     bool      `json:"is_expired"`
		HasExpiration bool      `json:"has_expiration"`
		IsAvailable   bool      `json:"is_available"`
		Created       time.Time `json:"created"`
		Updated       time.Time `json:"updated"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}
)
