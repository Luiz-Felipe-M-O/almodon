package requisition

import (
	"strings"

	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

const (
	destinationMaxLength   = 200
	notesMaxLength         = 5000
	approvalNotesMaxLength = 5000
)

type Status string

const (
	StatusPending        Status = "pending"
	StatusPendingNoStock Status = "pending_no_stock"
	StatusApproved       Status = "approved"
	StatusRejected       Status = "rejected"
	StatusFulfilled      Status = "fulfilled"
	StatusCancelled      Status = "cancelled"
)

type Entry struct {
	uuid     uuid.UUID
	material uuid.UUID
	quantity float64
}

type AnswerEntry struct {
	uuid             uuid.UUID
	requisitionEntry uuid.UUID
	approvedQuantity float64
	notes            string
}

type Answer struct {
	uuid     uuid.UUID
	approver uuid.UUID
	status   Status
	notes    string
	entries  map[uuid.UUID]AnswerEntry
}

type Requisition struct {
	uuid        uuid.UUID
	author      uuid.UUID
	notes       string
	destination string
	status      Status
	entries     map[uuid.UUID]Entry
	answers     map[uuid.UUID]Answer
	approver    uuid.UUID
	answeredAt  time.Time
}

func New(author uuid.UUID, notes, destination string, entries []CreateEntry) (Requisition, error) {
	var r Requisition

	err := errors.Join(
		r.SetAuthor(author),
		r.SetNotes(notes),
		r.SetDestination(destination),
	)
	if err != nil {
		return Requisition{}, xerrors.ErrRequisitionCreation.New(err)
	}

	r.entries = make(map[uuid.UUID]Entry)
	for _, entry := range entries {
		if err := r.AddEntry(entry.Material, entry.Quantity); err != nil {
			return Requisition{}, xerrors.ErrRequisitionCreation.New(err)
		}
	}

	r.status = StatusPending
	r.uuid = uuid.NewUUIDv7()
	r.answers = make(map[uuid.UUID]Answer)
	return r, nil
}

func (r *Requisition) UUID() uuid.UUID       { return r.uuid }
func (r *Requisition) Author() uuid.UUID     { return r.author }
func (r *Requisition) Notes() string         { return r.notes }
func (r *Requisition) Destination() string   { return r.destination }
func (r *Requisition) Status() Status        { return r.status }
func (r *Requisition) Entries() []Entry      { return r.getEntries() }
func (r *Requisition) Answers() []Answer     { return r.getAnswers() }
func (r *Requisition) Approver() uuid.UUID   { return r.approver }
func (r *Requisition) AnsweredAt() time.Time { return r.answeredAt }

func (r *Requisition) SetAuthor(author uuid.UUID) error {
	if author == uuid.Nil {
		return xerrors.ErrAuthorEmpty
	}
	r.author = author
	return nil
}

func (r *Requisition) SetNotes(notes string) error {
	return entity.Set(&r.notes, notes, ProcessNotes)
}

func (r *Requisition) SetDestination(destination string) error {
	return entity.Set(&r.destination, destination, ProcessDestination)
}

func (r *Requisition) SetStatus(status Status) error {
	validStatuses := map[Status]bool{
		StatusPending:        true,
		StatusPendingNoStock: true,
		StatusApproved:       true,
		StatusRejected:       true,
		StatusFulfilled:      true,
		StatusCancelled:      true,
	}

	if !validStatuses[status] {
		return xerrors.ErrInvalidStatus
	}

	r.status = status
	return nil
}

func (r *Requisition) AddEntry(material uuid.UUID, quantity float64) error {
	if material == uuid.Nil {
		return xerrors.ErrMaterialEmpty
	}

	if quantity <= 0 {
		return xerrors.ErrQuantityInvalid
	}

	entry := Entry{
		uuid:     uuid.NewUUIDv7(),
		material: material,
		quantity: quantity,
	}

	r.entries[entry.uuid] = entry
	return nil
}

func (r *Requisition) RemoveEntry(entryUUID uuid.UUID) error {
	if _, exists := r.entries[entryUUID]; !exists {
		return xerrors.ErrEntryNotFound
	}

	delete(r.entries, entryUUID)
	return nil
}

func (r *Requisition) Answer(approver uuid.UUID, status Status, notes string, answerEntries []AnswerEntryData) error {
	if err := r.SetStatus(status); err != nil {
		return err
	}

	answer := Answer{
		uuid:     uuid.NewUUIDv7(),
		approver: approver,
		status:   status,
		notes:    notes,
		entries:  make(map[uuid.UUID]AnswerEntry),
	}

	for _, ae := range answerEntries {
		if _, exists := r.entries[ae.RequisitionEntry]; !exists {
			return xerrors.ErrEntryNotFound
		}

		answerEntry := AnswerEntry{
			uuid:             uuid.NewUUIDv7(),
			requisitionEntry: ae.RequisitionEntry,
			approvedQuantity: ae.ApprovedQuantity,
			notes:            ae.Notes,
		}
		answer.entries[answerEntry.uuid] = answerEntry
	}

	r.answers[answer.uuid] = answer
	r.approver = approver
	r.answeredAt = time.Now()
	return nil
}

func (r *Requisition) getEntries() []Entry {
	entries := make([]Entry, 0, len(r.entries))
	for _, entry := range r.entries {
		entries = append(entries, entry)
	}
	return entries
}

func (r *Requisition) getAnswers() []Answer {
	answers := make([]Answer, 0, len(r.answers))
	for _, answer := range r.answers {
		answers = append(answers, answer)
	}
	return answers
}

func ProcessNotes(notes string) (string, error) {
	if len(notes) > notesMaxLength {
		return "", xerrors.ErrNotesTooLong
	}
	return strings.TrimSpace(notes), nil
}

func ProcessDestination(destination string) (string, error) {
	destination = strings.TrimSpace(destination)
	if destination == "" {
		return "", xerrors.ErrDestinationEmpty
	}
	if len(destination) > destinationMaxLength {
		return "", xerrors.ErrDestinationTooLong
	}
	return destination, nil
}

func ProcessApprovalNotes(notes string) (string, error) {
	if len(notes) > approvalNotesMaxLength {
		return "", xerrors.ErrApprovalNotesTooLong
	}
	return strings.TrimSpace(notes), nil
}

func ProcessQuantity(quantity float64) (float64, error) {
	if quantity <= 0 {
		return 0, xerrors.ErrQuantityInvalid
	}
	return quantity, nil
}

func ProcessApprovedQuantity(quantity float64) (float64, error) {
	if quantity < 0 {
		return 0, xerrors.ErrApprovedQuantityInvalid
	}
	return quantity, nil
}
