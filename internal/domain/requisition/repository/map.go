package requisitionrepo

import (
	"cmp"
	"encoding/json"
	"os"
	"sync"
	"time"
	"unsafe"

	"github.com/alan-b-lima/almodon/internal/domain/requisition"
	"github.com/alan-b-lima/almodon/internal/support/repository"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Map struct {
	uuidIndex map[uuid.UUID]int
	repo      []requisition.Entity
	mu        sync.RWMutex
	datapath  string
}

func NewMap() requisition.Repository {
	repo := Map{
		uuidIndex: make(map[uuid.UUID]int),
	}

	return &repo
}

func NewPersistentMap(datapath string) (requisition.Repository, error) {
	repo := Map{
		uuidIndex: make(map[uuid.UUID]int),
		datapath:  datapath,
	}

	if err := repo.init(); err != nil {
		return nil, err
	}

	return &repo, nil
}

func (m *Map) init() error {
	f, err := os.Open(m.datapath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var repo []entity
	if err := json.NewDecoder(f).Decode(&repo); err != nil {
		return err
	}

	m.repo = entity_from_json(repo)
	for i, record := range m.repo {
		m.uuidIndex[record.UUID] = i
	}

	return nil
}

func (m *Map) Close() error {
	defer m.mu.Unlock()
	m.mu.Lock()

	if m.datapath == "" {
		return nil
	}

	f, err := os.OpenFile(m.datapath, os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := f.Truncate(0); err != nil {
		return err
	}

	repo := json_for_entity(m.repo)
	if err := json.NewEncoder(f).Encode(repo); err != nil {
		return err
	}

	return nil
}

func (m *Map) List(offset, limit int, filters requisition.ListFilters) (requisition.Entities, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	// Apply filters
	filtered := make([]requisition.Entity, 0, len(m.repo))
	for _, record := range m.repo {
		if !matchesFilters(record, filters) {
			continue
		}
		filtered = append(filtered, record)
	}

	lo := clamp(0, offset, len(filtered))
	hi := clamp(0, offset+limit, len(filtered))

	if lo >= hi {
		return requisition.Entities{
			Records:      []requisition.Entity{},
			TotalRecords: len(filtered),
		}, nil
	}

	res := make([]requisition.Entity, hi-lo)
	copy(res, filtered[lo:hi])

	return requisition.Entities{
		Offset:       lo,
		Length:       len(res),
		Records:      res,
		TotalRecords: len(filtered),
	}, nil
}

func matchesFilters(record requisition.Entity, filters requisition.ListFilters) bool {
	if status, ok := filters.Status.Unwrap(); ok {
		if record.Status != status {
			return false
		}
	}

	if author, ok := filters.Author.Unwrap(); ok {
		if record.Author != author {
			return false
		}
	}

	if dateFrom, ok := filters.DateFrom.Unwrap(); ok {
		if record.Created.Before(dateFrom) {
			return false
		}
	}

	if dateTo, ok := filters.DateTo.Unwrap(); ok {
		if record.Created.After(dateTo) {
			return false
		}
	}

	if destination, ok := filters.Destination.Unwrap(); ok {
		if record.Destination != destination {
			return false
		}
	}

	return true
}

func (m *Map) Get(uuid uuid.UUID) (requisition.Entity, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return requisition.Entity{}, xerrors.ErrRequisitionNotFound
	}

	return m.repo[index], nil
}

func (m *Map) Create(req requisition.Entity) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	m.uuidIndex[req.UUID] = len(m.repo)
	m.repo = append(m.repo, req)

	return nil
}

func (m *Map) Patch(uuid uuid.UUID, req requisition.PartialEntity) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return xerrors.ErrRequisitionNotFound
	}

	r := &m.repo[index]

	repository.SomeThen(&r.Notes, req.Notes)
	repository.SomeThen(&r.Destination, req.Destination)
	repository.SomeThen(&r.Status, req.Status)
	repository.SomeThen(&r.Approver, req.Approver)
	repository.SomeThen(&r.AnsweredAt, req.AnsweredAt)

	r.Updated = req.Updated

	return nil
}

func (m *Map) Delete(uuid uuid.UUID) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return nil
	}

	r := &m.repo[index]

	delete(m.uuidIndex, r.UUID)

	// Update index for the last element that will be moved
	if index < len(m.repo)-1 {
		lastUUID := m.repo[len(m.repo)-1].UUID
		m.uuidIndex[lastUUID] = index
	}

	m.repo[index] = m.repo[len(m.repo)-1]
	m.repo = m.repo[:len(m.repo)-1]

	return nil
}

func clamp[T cmp.Ordered](mn, val, mx T) T {
	return min(max(mn, val), mx)
}

type entity struct {
	UUID        uuid.UUID          `json:"uuid"`
	Author      uuid.UUID          `json:"author"`
	Notes       string             `json:"notes"`
	Destination string             `json:"destination"`
	Status      requisition.Status `json:"status"`
	Entries     []entryEntity      `json:"entries"`
	Answers     []answerEntity     `json:"answers"`
	Approver    uuid.UUID          `json:"approver"`
	AnsweredAt  time.Time          `json:"answered_at"`
	Created     time.Time          `json:"created"`
	Updated     time.Time          `json:"updated"`
}

type entryEntity struct {
	UUID     uuid.UUID `json:"uuid"`
	Material uuid.UUID `json:"material"`
	Quantity float64   `json:"quantity"`
}

type answerEntity struct {
	UUID       uuid.UUID           `json:"uuid"`
	Approver   uuid.UUID           `json:"approver"`
	Status     requisition.Status  `json:"status"`
	Notes      string              `json:"notes"`
	Entries    []answerEntryEntity `json:"entries"`
	AnsweredAt time.Time           `json:"answered_at"`
}

type answerEntryEntity struct {
	UUID             uuid.UUID `json:"uuid"`
	RequisitionEntry uuid.UUID `json:"requisition_entry"`
	ApprovedQuantity float64   `json:"approved_quantity"`
	Notes            string    `json:"notes"`
}

func json_for_entity(entities []requisition.Entity) []entity {
	result := make([]entity, len(entities))
	for i, e := range entities {
		entries := make([]entryEntity, len(e.Entries))
		for j, entry := range e.Entries {
			entries[j] = entryEntity{
				UUID:     entry.UUID,
				Material: entry.Material,
				Quantity: entry.Quantity,
			}
		}

		answers := make([]answerEntity, len(e.Answers))
		for j, answer := range e.Answers {
			answerEntries := make([]answerEntryEntity, len(answer.Entries))
			for k, ae := range answer.Entries {
				answerEntries[k] = answerEntryEntity{
					UUID:             ae.UUID,
					RequisitionEntry: ae.RequisitionEntry,
					ApprovedQuantity: ae.ApprovedQuantity,
					Notes:            ae.Notes,
				}
			}

			answers[j] = answerEntity{
				UUID:       answer.UUID,
				Approver:   answer.Approver,
				Status:     answer.Status,
				Notes:      answer.Notes,
				Entries:    answerEntries,
				AnsweredAt: answer.AnsweredAt,
			}
		}

		result[i] = entity{
			UUID:        e.UUID,
			Author:      e.Author,
			Notes:       e.Notes,
			Destination: e.Destination,
			Status:      e.Status,
			Entries:     entries,
			Answers:     answers,
			Approver:    e.Approver,
			AnsweredAt:  e.AnsweredAt,
			Created:     e.Created,
			Updated:     e.Updated,
		}
	}
	return result
}

func entity_from_json(entities []entity) []requisition.Entity {
	result := make([]requisition.Entity, len(entities))
	for i, e := range entities {
		entries := make([]requisition.EntryEntity, len(e.Entries))
		for j, entry := range e.Entries {
			entries[j] = requisition.EntryEntity{
				UUID:     entry.UUID,
				Material: entry.Material,
				Quantity: entry.Quantity,
			}
		}

		answers := make([]requisition.AnswerEntity, len(e.Answers))
		for j, answer := range e.Answers {
			answerEntries := make([]requisition.AnswerEntryEntity, len(answer.Entries))
			for k, ae := range answer.Entries {
				answerEntries[k] = requisition.AnswerEntryEntity{
					UUID:             ae.UUID,
					RequisitionEntry: ae.RequisitionEntry,
					ApprovedQuantity: ae.ApprovedQuantity,
					Notes:            ae.Notes,
				}
			}

			answers[j] = requisition.AnswerEntity{
				UUID:       answer.UUID,
				Approver:   answer.Approver,
				Status:     answer.Status,
				Notes:      answer.Notes,
				Entries:    answerEntries,
				AnsweredAt: answer.AnsweredAt,
			}
		}

		result[i] = requisition.Entity{
			UUID:        e.UUID,
			Author:      e.Author,
			Notes:       e.Notes,
			Destination: e.Destination,
			Status:      e.Status,
			Entries:     entries,
			Answers:     answers,
			Approver:    e.Approver,
			AnsweredAt:  e.AnsweredAt,
			Created:     e.Created,
			Updated:     e.Updated,
		}
	}
	return result
}
