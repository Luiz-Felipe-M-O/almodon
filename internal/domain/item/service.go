package item

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/money"
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
	History(context.Context, uuid.UUID) (HistoryResult, error)

	Create(context.Context, Create) (CreateResult, error)

	UpdateAmount(context.Context, uuid.UUID, UpdateAmount) error
	Patch(context.Context, uuid.UUID, Patch) error

	Delete(context.Context, uuid.UUID) error
}

type (
	Create struct {
		Material uuid.UUID   `json:"material"`
		Amount   float64     `json:"amount"`
		UnitCost money.Money `json:"unit_cost"`
		Expires  time.Time   `json:"expires"`
	}

	UpdateAmount struct {
		Amount float64 `json:"amount"`
	}

	Patch struct {
		Material opt.Opt[uuid.UUID]   `json:"material"`
		UnitCost opt.Opt[money.Money] `json:"unit_cost"`
		Expires  opt.Opt[time.Time]   `json:"expires"`
	}
)

type (
	Result struct {
		UUID        uuid.UUID   `json:"uuid"`
		Name        string      `json:"name"`
		ECampus     int         `json:"ecampus"`
		CATMAT      int         `json:"catmat"`
		SIADS       int         `json:"siads"`
		Material    uuid.UUID   `json:"material"`
		Amount      float64     `json:"amount"`
		AmountFlag  Stock       `json:"amount_flag"`
		UnitCost    money.Money `json:"unit_cost"`
		Unit        string      `json:"unit"`
		Expires     time.Time   `json:"expires"`
		ExpiresFlag Expiration  `json:"expires_flag"`
		Created     time.Time   `json:"created"`
		Updated     time.Time   `json:"updated"`
	}

	HistoryResult struct {
		UUID     uuid.UUID    `json:"uuid"`
		Version  int          `json:"version"`
		Created  time.Time    `json:"created"`
		Updated  time.Time    `json:"updated"`
		Versions []PastResult `json:"versions"`
	}

	PastResult struct {
		Version  int         `json:"version"`
		Material uuid.UUID   `json:"material"`
		Amount   float64     `json:"amount"`
		UnitCost money.Money `json:"unit_cost"`
		Expires  time.Time   `json:"expires"`
		Created  time.Time   `json:"created"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}
)
