package requisition

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type (
	ListParams struct {
		Offset      int             `query:"offset"`
		Limit       int             `query:"limit"`
		Status      opt.Opt[Status] `query:"status"`
		Author      opt.Opt[string] `query:"author"`    // UUID as string
		DateFrom    opt.Opt[string] `query:"date_from"` // ISO date
		DateTo      opt.Opt[string] `query:"date_to"`   // ISO date
		Destination opt.Opt[string] `query:"destination"`
	}

	Create struct {
		Notes       string        `json:"notes"`
		Destination string        `json:"destination"`
		Entries     []CreateEntry `json:"entries"`
	}

	CreateEntry struct {
		Material uuid.UUID `json:"material"`
		Quantity float64   `json:"quantity"`
	}

	Patch struct {
		Notes       opt.Opt[string] `json:"notes"`
		Destination opt.Opt[string] `json:"destination"`
	}

	AddEntry struct {
		Material uuid.UUID `json:"material"`
		Quantity float64   `json:"quantity"`
	}

	AnswerRequisition struct {
		Status  Status            `json:"status"`
		Notes   string            `json:"notes"`
		Entries []AnswerEntryData `json:"entries"`
	}

	AnswerEntryData struct {
		RequisitionEntry uuid.UUID `json:"requisition_entry"`
		ApprovedQuantity float64   `json:"approved_quantity"`
		Notes            string    `json:"notes"`
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
		UUID        uuid.UUID      `json:"uuid"`
		Author      uuid.UUID      `json:"author"`
		Notes       string         `json:"notes"`
		Destination string         `json:"destination"`
		Status      Status         `json:"status"`
		Entries     []EntryResult  `json:"entries"`
		Answers     []AnswerResult `json:"answers"`
		Approver    uuid.UUID      `json:"approver,omitempty"`
		AnsweredAt  time.Time      `json:"answered_at,omitempty"`
		Created     time.Time      `json:"created"`
		Updated     time.Time      `json:"updated"`
	}

	EntryResult struct {
		UUID     uuid.UUID `json:"uuid"`
		Material uuid.UUID `json:"material"`
		Quantity float64   `json:"quantity"`
	}

	AnswerResult struct {
		UUID       uuid.UUID           `json:"uuid"`
		Approver   uuid.UUID           `json:"approver"`
		Status     Status              `json:"status"`
		Notes      string              `json:"notes"`
		Entries    []AnswerEntryResult `json:"entries"`
		AnsweredAt time.Time           `json:"answered_at"`
	}

	AnswerEntryResult struct {
		UUID             uuid.UUID `json:"uuid"`
		RequisitionEntry uuid.UUID `json:"requisition_entry"`
		ApprovedQuantity float64   `json:"approved_quantity"`
		Notes            string    `json:"notes"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}

	AddEntryResult struct {
		UUID uuid.UUID `json:"uuid"`
	}
)
