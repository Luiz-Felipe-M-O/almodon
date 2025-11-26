package material

import (
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"time"
)

type (
	ListRequest struct {
		Offset int `query:"offset"`
		Limit  int `query:"limit"`
	}

	GetRequest struct {
		UUID uuid.UUID `json:"-"`
	}

	GetBySIADSRequest struct {
		SIADS string `json:"-"`
	}

	GetByCATMATRequest struct {
		CATMAT string `json:"-"`
	}

	GetByECAMPUSRequest struct {
		ECAMPUS string `json:"-"`
	}

	CreateRequest struct {
		Name        string  `json:"name"`
		SIADS       string  `json:"siads"`
		CATMAT      string  `json:"catmat"`
		ECAMPUS     string  `json:"ecampus"`
		Description string  `json:"description"`
		Unit        string  `json:"unit"`
		MinQuantity float64 `json:"min_quantity"`
	}

	PatchRequest struct {
		UUID        uuid.UUID        `json:"-"`
		Name        opt.Opt[string]  `json:"name"`
		SIADS       opt.Opt[string]  `json:"siads"`
		CATMAT      opt.Opt[string]  `json:"catmat"`
		ECAMPUS     opt.Opt[string]  `json:"ecampus"`
		Description opt.Opt[string]  `json:"description"`
		Unit        opt.Opt[string]  `json:"unit"`
		MinQuantity opt.Opt[float64] `json:"min_quantity"`
	}

	UpdateMinQuantityRequest struct {
		UUID        uuid.UUID `json:"-"`
		MinQuantity float64   `json:"min_quantity"`
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
		UUID        uuid.UUID `json:"uuid"`
		Name        string    `json:"name"`
		SIADS       string    `json:"siads"`
		CATMAT      string    `json:"catmat"`
		ECAMPUS     string    `json:"ecampus"`
		Description string    `json:"description"`
		Unit        string    `json:"unit"`
		MinQuantity float64   `json:"min_quantity"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}
)
