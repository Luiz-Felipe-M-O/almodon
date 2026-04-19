package item

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/opt"
)

type Service interface {
	List(context.Context) ([]Result, error)
	ListByMaterial(context.Context, uuid.UUID) ([]Result, error)
	ListByECampus(context.Context, int) ([]Result, error)
	ListByCATMAT(context.Context, int) ([]Result, error)
	ListBySIADS(context.Context, int) ([]Result, error)

	Get(context.Context, uuid.UUID) (Result, error)

	Create(context.Context, Create) (CreateResult, error)

	UpdateAmount(context.Context, uuid.UUID, UpdateAmount) error
	Patch(context.Context, uuid.UUID, Patch) error

	Delete(context.Context, uuid.UUID) error
}

type (
	Create struct {
		Material uuid.UUID `json:"material"`
		Amount   float64   `json:"amount"`
		UnitCost float64   `json:"unit_cost"`
		Arrival  time.Time `json:"arrival"`
		Expires  time.Time `json:"expires"`
	}

	UpdateAmount struct {
		Amount float64 `json:"amount"`
	}

	Patch struct {
		Material opt.Opt[uuid.UUID] `json:"material"`
		UnitCost opt.Opt[float64]   `json:"unit_cost"`
		Arrival  opt.Opt[time.Time] `json:"arrival"`
		Expires  opt.Opt[time.Time] `json:"expires"`
	}
)

type (
	Result struct {
		UUID        uuid.UUID  `json:"uuid"`
		Name        string     `json:"name"`
		ECampus     int        `json:"ecampus"`
		CATMAT      int        `json:"catmat"`
		SIADS       int        `json:"siads"`
		Material    uuid.UUID  `json:"material"`
		Amount      float64    `json:"amount"`
		AmountFlag  Stock      `json:"amount_flag"`
		UnitCost    float64    `json:"unit_cost"`
		Unit        string     `json:"unit"`
		Arrival     time.Time  `json:"arrival"`
		Expires     time.Time  `json:"expires"`
		ExpiresFlag Expiration `json:"expires_flag"`
		Created     time.Time  `json:"created"`
		Updated     time.Time  `json:"updated"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}
)
