package itemrepo

import (
	"cmp"
	"encoding/json"
	"os"
	"sync"

	"github.com/alan-b-lima/almodon/internal/domain/item"
	repo "github.com/alan-b-lima/almodon/internal/support/repository"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Map struct {
	uuid     repo.Index[uuid.UUID, int]
	batch    repo.SliceIndex[uuid.UUID, int]
	material repo.SliceIndex[uuid.UUID, int]

	repo []item.Entity
	mu   sync.RWMutex

	datapath string
}

func NewMap() item.Repository {
	repo := Map{
		uuid:     make(repo.Index[uuid.UUID, int]),
		batch:    make(repo.SliceIndex[uuid.UUID, int]),
		material: make(repo.SliceIndex[uuid.UUID, int]),
	}

	return &repo
}

func NewPersistentMap(datapath string) (item.Repository, error) {
	repo := Map{
		uuid:     make(repo.Index[uuid.UUID, int]),
		batch:    make(repo.SliceIndex[uuid.UUID, int]),
		material: make(repo.SliceIndex[uuid.UUID, int]),
		datapath: datapath,
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

	if err := json.NewDecoder(f).Decode(&m.repo); err != nil {
		return err
	}

	for i, record := range m.repo {
		m.uuid[record.UUID] = i

		m.batch.Add(record.Batch, i)
		m.material.Add(record.Material, i)
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

	if err := json.NewEncoder(f).Encode(m.repo); err != nil {
		return err
	}

	return nil
}

func (m *Map) List(offset, limit int) (item.Entities, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	lo := clamp(0, offset, len(m.repo))
	hi := clamp(0, offset+limit, len(m.repo))

	if lo >= hi {
		return item.Entities{
			Records:      []item.Entity{},
			TotalRecords: len(m.repo),
		}, nil
	}

	res := make([]item.Entity, hi-lo)
	copy(res, m.repo[lo:hi])

	return item.Entities{
		Offset:       lo,
		Length:       len(res),
		Records:      res,
		TotalRecords: len(m.repo),
	}, nil
}

func (m *Map) ListByBatch(batch uuid.UUID) (item.Entities, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	indices, in := m.batch.Get(batch)
	if !in || len(indices) == 0 {
		return item.Entities{
			Records:      []item.Entity{},
			TotalRecords: 0,
		}, nil
	}

	res := make([]item.Entity, len(indices))
	for i, idx := range indices {
		res[i] = m.repo[idx]
	}

	return item.Entities{
		Offset:       0,
		Length:       len(res),
		Records:      res,
		TotalRecords: len(res),
	}, nil
}

func (m *Map) ListByMaterial(material uuid.UUID) (item.Entities, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	indices, in := m.material.Get(material)
	if !in || len(indices) == 0 {
		return item.Entities{
			Records:      []item.Entity{},
			TotalRecords: 0,
		}, nil
	}

	res := make([]item.Entity, len(indices))
	for i, idx := range indices {
		res[i] = m.repo[idx]
	}

	return item.Entities{
		Offset:       0,
		Length:       len(res),
		Records:      res,
		TotalRecords: len(res),
	}, nil
}

func (m *Map) Get(uuid uuid.UUID) (item.Entity, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.uuid[uuid]
	if !in {
		return item.Entity{}, xerrors.ErrItemNotFound
	}

	return m.repo[index], nil
}

func (m *Map) Create(itm item.Entity) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index := len(m.repo)

	m.uuid.Set(itm.UUID, index)
	m.batch.Add(itm.Batch, index)
	m.material.Add(itm.Material, index)

	m.repo = append(m.repo, itm)
	return nil
}

func (m *Map) Patch(uuid uuid.UUID, partial item.PartialEntity) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuid[uuid]
	if !in {
		return xerrors.ErrItemNotFound
	}

	item := &m.repo[index]

	if batch, ok := partial.Batch.Unwrap(); ok {
		m.batch.Del(item.Batch, index)
		m.batch.Add(batch, index)
		item.Batch = batch
	}

	if material, ok := partial.Material.Unwrap(); ok {
		m.material.Del(item.Material, index)
		m.material.Add(material, index)
		item.Material = material
	}

	repo.SomeThen(&item.Quantity, partial.Quantity)
	repo.SomeThen(&item.Expiration, partial.Expiration)
	item.Updated = partial.Updated

	return nil
}

func (m *Map) Delete(uuid uuid.UUID) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuid[uuid]
	if !in {
		return nil
	}

	itm := &m.repo[index]

	delete(m.uuid, itm.UUID)
	m.batch.Del(itm.Batch, index)
	m.material.Del(itm.Material, index)

	lastIndex := len(m.repo) - 1
	if index != lastIndex {
		m.repo[index] = m.repo[lastIndex]

		lastItm := &m.repo[index]

		m.uuid.Set(lastItm.UUID, index)
		m.batch.Swap(lastItm.Batch, lastIndex, index)
		m.material.Swap(lastItm.Material, lastIndex, index)
	}

	m.repo = m.repo[:lastIndex]

	return nil
}

func clamp[T cmp.Ordered](mn, val, mx T) T {
	return min(max(mn, val), mx)
}
