package user

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/opt"
)

type Store interface {
	List(context.Context) ([]Record, error)
	CountChiefs(context.Context) (int, error)

	Get(context.Context, uuid.UUID) (Record, error)
	GetBySIAPE(context.Context, string) (Record, error)

	Create(context.Context, CreateRecord) error

	Patch(context.Context, uuid.UUID, PatchRecord) error

	Delete(context.Context, uuid.UUID) error

	RunTx(context.Context, func(Store) error) error
}

type (
	Record struct {
		UUID     uuid.UUID
		SIAPE    string
		Name     string
		Email    string
		Password []byte
		Role     auth.Role
		Logged   bool
		Created  time.Time
		Updated  time.Time
	}

	CreateRecord struct {
		UUID     uuid.UUID
		SIAPE    string
		Name     string
		Email    string
		Password []byte
		Role     auth.Role
		Created  time.Time
		Updated  time.Time
	}

	PatchRecord struct {
		Name    opt.Opt[string]
		Email   opt.Opt[string]
		Updated time.Time
	}
)
