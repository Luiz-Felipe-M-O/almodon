package requisitionserve

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/requisition"
	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Core struct {
	Requisitions requisition.Repository
}

var _ requisition.Service = &Core{}

func (c *Core) List(req requisition.ListParams) (requisition.Entities, error) {
	filters := requisition.ListFilters{
		Status:      req.Status,
		Author:      parseUUID(req.Author),
		DateFrom:    parseTime(req.DateFrom),
		DateTo:      parseTime(req.DateTo),
		Destination: req.Destination,
	}

	return c.Requisitions.List(req.Offset, req.Limit, filters)
}

func (c *Core) Get(uuid uuid.UUID) (requisition.Entity, error) {
	return c.Requisitions.Get(uuid)
}

func (c *Core) Create(req requisition.Create) (uuid.UUID, error) {
	entries := make([]requisition.CreateEntry, len(req.Entries))
	for i, e := range req.Entries {
		entries[i] = requisition.CreateEntry{
			Material: e.Material,
			Quantity: e.Quantity,
		}
	}

	r, err := requisition.New(
		uuid.UUID{}, // Will be set by auth layer
		req.Notes,
		req.Destination,
		entries,
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	ent := translate(&r)

	now := time.Now()
	ent.Created = now
	ent.Updated = now

	return r.UUID(), c.Requisitions.Create(ent)
}

func (c *Core) Patch(uuid uuid.UUID, req requisition.Patch) error {
	// Check if requisition can be modified
	existing, err := c.Requisitions.Get(uuid)
	if err != nil {
		return err
	}

	if existing.Status != requisition.StatusPending &&
		existing.Status != requisition.StatusPendingNoStock {
		return xerrors.ErrCannotModifyAnswered
	}

	return patch(c.Requisitions, uuid, req.Notes, req.Destination)
}

func (c *Core) Delete(uuid uuid.UUID) error {
	// Check if requisition can be deleted
	existing, err := c.Requisitions.Get(uuid)
	if err != nil {
		return err
	}

	if existing.Status != requisition.StatusPending &&
		existing.Status != requisition.StatusPendingNoStock {
		return xerrors.ErrCannotModifyAnswered
	}

	return c.Requisitions.Delete(uuid)
}

func (c *Core) AddEntry(requisitionUUID uuid.UUID, req requisition.AddEntry) (uuid.UUID, error) {
	existing, err := c.Requisitions.Get(requisitionUUID)
	if err != nil {
		return uuid.UUID{}, err
	}

	if existing.Status != requisition.StatusPending &&
		existing.Status != requisition.StatusPendingNoStock {
		return uuid.UUID{}, xerrors.ErrCannotModifyAnswered
	}

	// Validate entry
	if req.Material == uuid.Nil {
		return uuid.UUID{}, xerrors.ErrMaterialEmpty
	}
	if req.Quantity <= 0 {
		return uuid.UUID{}, xerrors.ErrQuantityInvalid
	}

	entryUUID := uuid.NewUUIDv7()
	newEntry := requisition.EntryEntity{
		UUID:     entryUUID,
		Material: req.Material,
		Quantity: req.Quantity,
	}

	existing.Entries = append(existing.Entries, newEntry)
	existing.Updated = time.Now()

	partial := requisition.PartialEntity{
		Updated: existing.Updated,
	}

	if err := c.Requisitions.Patch(requisitionUUID, partial); err != nil {
		return uuid.UUID{}, err
	}

	return entryUUID, nil
}

func (c *Core) RemoveEntry(requisitionUUID, entryUUID uuid.UUID) error {
	existing, err := c.Requisitions.Get(requisitionUUID)
	if err != nil {
		return err
	}

	if existing.Status != requisition.StatusPending &&
		existing.Status != requisition.StatusPendingNoStock {
		return xerrors.ErrCannotModifyAnswered
	}

	// Find and remove entry
	found := false
	newEntries := make([]requisition.EntryEntity, 0, len(existing.Entries)-1)
	for _, entry := range existing.Entries {
		if entry.UUID == entryUUID {
			found = true
			continue
		}
		newEntries = append(newEntries, entry)
	}

	if !found {
		return xerrors.ErrEntryNotFound
	}

	if len(newEntries) == 0 {
		return xerrors.ErrRequisitionMustHaveEntries
	}

	existing.Entries = newEntries
	existing.Updated = time.Now()

	partial := requisition.PartialEntity{
		Updated: existing.Updated,
	}

	return c.Requisitions.Patch(requisitionUUID, partial)
}

func (c *Core) Answer(requisitionUUID uuid.UUID, req requisition.AnswerRequisition) error {
	existing, err := c.Requisitions.Get(requisitionUUID)
	if err != nil {
		return err
	}

	if existing.Status != requisition.StatusPending &&
		existing.Status != requisition.StatusPendingNoStock {
		return xerrors.ErrCannotAnswerTwice
	}

	// Validate status transition
	if req.Status != requisition.StatusApproved && req.Status != requisition.StatusRejected {
		return xerrors.ErrInvalidStatus
	}

	// Validate answer entries match requisition entries
	entryMap := make(map[uuid.UUID]bool)
	for _, entry := range existing.Entries {
		entryMap[entry.UUID] = true
	}

	answerEntries := make([]requisition.AnswerEntryEntity, len(req.Entries))
	for i, ae := range req.Entries {
		if !entryMap[ae.RequisitionEntry] {
			return xerrors.ErrEntryNotFound
		}

		if ae.ApprovedQuantity < 0 {
			return xerrors.ErrApprovedQuantityInvalid
		}

		answerEntries[i] = requisition.AnswerEntryEntity{
			UUID:             uuid.NewUUIDv7(),
			RequisitionEntry: ae.RequisitionEntry,
			ApprovedQuantity: ae.ApprovedQuantity,
			Notes:            ae.Notes,
		}
	}

	now := time.Now()
	answer := requisition.AnswerEntity{
		UUID:       uuid.NewUUIDv7(),
		Approver:   uuid.UUID{}, // Will be set by auth layer
		Status:     req.Status,
		Notes:      req.Notes,
		Entries:    answerEntries,
		AnsweredAt: now,
	}

	existing.Answers = append(existing.Answers, answer)
	existing.Status = req.Status
	existing.Approver = answer.Approver
	existing.AnsweredAt = now
	existing.Updated = now

	partial := requisition.PartialEntity{
		Status:     opt.Some(req.Status),
		Approver:   opt.Some(answer.Approver),
		AnsweredAt: opt.Some(now),
		Updated:    now,
	}

	return c.Requisitions.Patch(requisitionUUID, partial)
}

func (c *Core) Cancel(requisitionUUID uuid.UUID) error {
	existing, err := c.Requisitions.Get(requisitionUUID)
	if err != nil {
		return err
	}

	if existing.Status == requisition.StatusFulfilled ||
		existing.Status == requisition.StatusCancelled {
		return xerrors.ErrCannotModifyAnswered
	}

	partial := requisition.PartialEntity{
		Status:  opt.Some(requisition.StatusCancelled),
		Updated: time.Now(),
	}

	return c.Requisitions.Patch(requisitionUUID, partial)
}

func (c *Core) MarkFulfilled(requisitionUUID uuid.UUID) error {
	existing, err := c.Requisitions.Get(requisitionUUID)
	if err != nil {
		return err
	}

	if existing.Status != requisition.StatusApproved {
		return xerrors.ErrCannotFulfillUnapproved
	}

	partial := requisition.PartialEntity{
		Status:  opt.Some(requisition.StatusFulfilled),
		Updated: time.Now(),
	}

	return c.Requisitions.Patch(requisitionUUID, partial)
}

func patch(
	repo requisition.Patcher,
	uuid uuid.UUID,
	notes, destination opt.Opt[string],
) error {
	var r requisition.PartialEntity

	err := errors.Join(
		entity.SomeThen(&r.Notes, notes, requisition.ProcessNotes),
		entity.SomeThen(&r.Destination, destination, requisition.ProcessDestination),
	)
	if err != nil {
		return xerrors.ErrRequisitionUpdate.New(err)
	}

	r.Updated = time.Now()

	return repo.Patch(uuid, r)
}

func translate(e *requisition.Requisition) requisition.Entity {
	entries := e.Entries()
	entryEntities := make([]requisition.EntryEntity, len(entries))
	for i, entry := range entries {
		entryEntities[i] = requisition.EntryEntity{
			UUID:     entry.UUID(),
			Material: entry.Material(),
			Quantity: entry.Quantity(),
		}
	}

	return requisition.Entity{
		UUID:        e.UUID(),
		Author:      e.Author(),
		Notes:       e.Notes(),
		Destination: e.Destination(),
		Status:      e.Status(),
		Entries:     entryEntities,
		Answers:     []requisition.AnswerEntity{},
		Approver:    uuid.Nil,
		AnsweredAt:  time.Time{},
	}
}

func parseUUID(s opt.Opt[string]) opt.Opt[uuid.UUID] {
	str, ok := s.Unwrap()
	if !ok {
		return opt.None[uuid.UUID]()
	}

	id, err := uuid.FromString(str)
	if err != nil {
		return opt.None[uuid.UUID]()
	}

	return opt.Some(id)
}

func parseTime(s opt.Opt[string]) opt.Opt[time.Time] {
	str, ok := s.Unwrap()
	if !ok {
		return opt.None[time.Time]()
	}

	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return opt.None[time.Time]()
	}

	return opt.Some(t)
}
