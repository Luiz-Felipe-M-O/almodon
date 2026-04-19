package item

import (
	"context"
	"time"

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

	Create(context.Context, CreateRecord) error

	UpdateAmount(context.Context, uuid.UUID, float64) error
	Patch(context.Context, uuid.UUID, PatchRecord) error

	Delete(context.Context, uuid.UUID) error
}

type (
	Record struct {
		UUID     uuid.UUID
		Name     string
		ECampus  int
		CATMAT   int
		SIADS    int
		Material uuid.UUID
		Amount   float64
		UnitCost float64
		Unit     string
		Arrival  time.Time
		Expires  time.Time
		Min      float64
		Created  time.Time
		Updated  time.Time
	}

	CreateRecord struct {
		UUID     uuid.UUID
		Material uuid.UUID
		Amount   float64
		UnitCost float64
		Arrival  time.Time
		Expires  time.Time
		Created  time.Time
		Updated  time.Time
	}

	PatchRecord struct {
		Material opt.Opt[uuid.UUID]
		UnitCost opt.Opt[float64]
		Arrival  opt.Opt[time.Time]
		Expires  opt.Opt[time.Time]
		Updated  time.Time
	}
)
