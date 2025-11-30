package requisition

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Repository interface {
	Lister
	Getter
	Creater
	Patcher
	Deleter
}

type (
	Lister interface {
		List(offset, limit int, filters ListFilters) (Entities, error)
	}

	Getter interface {
		Get(uuid.UUID) (Entity, error)
	}

	Creater interface {
		Create(Entity) error
	}

	Patcher interface {
		Patch(uuid.UUID, PartialEntity) error
	}

	Deleter interface {
		Delete(uuid.UUID) error
	}
)

type ListFilters struct {
	Status      opt.Opt[Status]
	Author      opt.Opt[uuid.UUID]
	DateFrom    opt.Opt[time.Time]
	DateTo      opt.Opt[time.Time]
	Destination opt.Opt[string]
}

type (
	Entities struct {
		Offset       int
		Length       int
		Records      []Entity
		TotalRecords int
	}

	Entity struct {
		UUID        uuid.UUID
		Author      uuid.UUID
		Notes       string
		Destination string
		Status      Status
		Entries     []EntryEntity
		Answers     []AnswerEntity
		Approver    uuid.UUID
		AnsweredAt  time.Time
		Created     time.Time
		Updated     time.Time
	}

	EntryEntity struct {
		UUID     uuid.UUID
		Material uuid.UUID
		Quantity float64
	}

	AnswerEntity struct {
		UUID       uuid.UUID
		Approver   uuid.UUID
		Status     Status
		Notes      string
		Entries    []AnswerEntryEntity
		AnsweredAt time.Time
	}

	AnswerEntryEntity struct {
		UUID             uuid.UUID
		RequisitionEntry uuid.UUID
		ApprovedQuantity float64
		Notes            string
	}

	PartialEntity struct {
		Notes       opt.Opt[string]
		Destination opt.Opt[string]
		Status      opt.Opt[Status]
		Approver    opt.Opt[uuid.UUID]
		AnsweredAt  opt.Opt[time.Time]
		Updated     time.Time
	}
)
