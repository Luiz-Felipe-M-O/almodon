package item

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/money"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/opt"
)

type Store interface {
	List(context.Context) ([]Record, error)
	ListByMaterial(context.Context, uuid.UUID) ([]Record, error)
	ListByECampus(context.Context, int) ([]Record, error)
	ListByCATMAT(context.Context, int) ([]Record, error)
	ListBySIADS(context.Context, int) ([]Record, error)

	Get(context.Context, uuid.UUID) (Record, error)

	History(context.Context, uuid.UUID) (HistoryRecord, error)

	Create(context.Context, Entity) error

	Patch(context.Context, uuid.UUID, PatchEntity) error

	Delete(context.Context, uuid.UUID) error

	RunTx(context.Context, func(Store) error) error
}

type (
	Record struct {
		UUID     uuid.UUID
		Version  int
		Name     string
		ECampus  int
		CATMAT   int
		SIADS    int
		Material uuid.UUID
		Amount   float64
		UnitCost money.Money
		Unit     string
		Expires  time.Time
		Min      float64
		Created  time.Time
		Updated  time.Time
	}

	HistoryRecord struct {
		UUID     uuid.UUID
		Version  int
		Created  time.Time
		Updated  time.Time
		Versions []PastRecord
	}

	PastRecord struct {
		Version  int
		Material uuid.UUID
		Amount   float64
		UnitCost money.Money
		Expires  time.Time
		Created  time.Time
	}
)

type (
	Entity struct {
		UUID     uuid.UUID
		Material uuid.UUID
		Amount   float64
		UnitCost money.Money
		Expires  time.Time
		Created  time.Time
		Updated  time.Time
	}

	PatchEntity struct {
		Material opt.Opt[uuid.UUID]
		Amount   opt.Opt[float64]
		UnitCost opt.Opt[money.Money]
		Expires  opt.Opt[time.Time]
		Updated  time.Time
	}
)
