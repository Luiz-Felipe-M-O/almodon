package item

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type (
	ListRequest struct {
		Offset int `query:"offset"`
		Limit  int `query:"limit"`
	}

	GetRequest struct {
		UUID uuid.UUID `json:"-"`
	}

	GetByBatchRequest struct {
		Batch uuid.UUID `json:"-"`
	}

	GetByMaterialRequest struct {
		Material uuid.UUID `json:"-"`
	}

	CreateRequest struct {
		Batch      uuid.UUID `json:"batch"`
		Material   uuid.UUID `json:"material"`
		Quantity   float64   `json:"quantity"`
		Expiration time.Time `json:"expiration"`
	}

	PatchRequest struct {
		UUID       uuid.UUID          `json:"-"`
		Batch      opt.Opt[uuid.UUID] `json:"batch"`
		Material   opt.Opt[uuid.UUID] `json:"material"`
		Quantity   opt.Opt[float64]   `json:"quantity"`
		Expiration opt.Opt[time.Time] `json:"expiration"`
	}

	UpdateQuantityRequest struct {
		UUID     uuid.UUID `json:"-"`
		Quantity float64   `json:"quantity"`
	}

	DeleteRequest struct {
		UUID uuid.UUID `json:"-"`
	}
)

type (
	ListResponse struct {
		Offset       int        `json:"offset"`
		Length       int        `json:"length"`
		Records      []Response `json:"records"`
		TotalRecords int        `json:"total_records"`
	}

	Response struct {
		UUID       uuid.UUID `json:"uuid"`
		Batch      uuid.UUID `json:"batch"`
		Material   uuid.UUID `json:"material"`
		Quantity   float64   `json:"quantity"`
		Expiration time.Time `json:"expiration"`
		IsExpired  bool      `json:"is_expired"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}
)
