package itemrepo

import (
	"cmp"
	"encoding/json"
	"os"
	"sync"

	"github.com/alan-b-lima/almodon/internal/domain/item"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Map struct {
	uuidIndex     map[uuid.UUID]int
	batchIndex    map[uuid.UUID][]int
	materialIndex map[uuid.UUID][]int

	repo     []item.Entity
	mu       sync.RWMutex
	datapath string
}

func NewMap() item.Repository {
	repo := Map{
		uuidIndex:     make(map[uuid.UUID]int),
		batchIndex:    make(map[uuid.UUID][]int),
		materialIndex: make(map[uuid.UUID][]int),
	}

	return &repo
}

func NewPersistentMap(datapath string) (item.Repository, error) {
	repo := Map{
		uuidIndex:     make(map[uuid.UUID]int),
		batchIndex:    make(map[uuid.UUID][]int),
		materialIndex: make(map[uuid.UUID][]int),
		datapath:      datapath,
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

	var repo []item.Entity
	if err := json.NewDecoder(f).Decode(&repo); err != nil {
		return err
	}

	m.repo = repo
	for i, record := range m.repo {
		m.uuidIndex[record.UUID] = i
		m.batchIndex[record.Batch] = append(m.batchIndex[record.Batch], i)
		m.materialIndex[record.Material] = append(m.materialIndex[record.Material], i)
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

func (m *Map) Get(uuid uuid.UUID) (item.Entity, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return item.Entity{}, xerrors.ErrItemNotFound
	}

	return m.repo[index], nil
}

func (m *Map) ListByBatch(batch uuid.UUID) (item.Entities, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	indices, in := m.batchIndex[batch]
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

	indices, in := m.materialIndex[material]
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

func (m *Map) Create(itm item.Entity) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index := len(m.repo)
	m.uuidIndex[itm.UUID] = index
	m.batchIndex[itm.Batch] = append(m.batchIndex[itm.Batch], index)
	m.materialIndex[itm.Material] = append(m.materialIndex[itm.Material], index)
	m.repo = append(m.repo, itm)

	return nil
}

func (m *Map) Patch(uuid uuid.UUID, itm item.PartialEntity) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return xerrors.ErrItemNotFound
	}

	entity := &m.repo[index]

	if batch, ok := itm.Batch.Unwrap(); ok && batch != entity.Batch {
		m.removeFromIndex(m.batchIndex, entity.Batch, index)
		m.batchIndex[batch] = append(m.batchIndex[batch], index)
		entity.Batch = batch
	}

	if material, ok := itm.Material.Unwrap(); ok && material != entity.Material {
		m.removeFromIndex(m.materialIndex, entity.Material, index)
		m.materialIndex[material] = append(m.materialIndex[material], index)
		entity.Material = material
	}

	someThen(&entity.Quantity, itm.Quantity)
	someThen(&entity.Expiration, itm.Expiration)
	someThen(&entity.CreatedAt, itm.CreatedAt)
	someThen(&entity.UpdatedAt, itm.UpdatedAt)

	return nil
}

func (m *Map) Delete(uuid uuid.UUID) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuidIndex[uuid]
	if !in {
		return nil
	}

	entity := &m.repo[index]

	delete(m.uuidIndex, entity.UUID)
	m.removeFromIndex(m.batchIndex, entity.Batch, index)
	m.removeFromIndex(m.materialIndex, entity.Material, index)

	lastIndex := len(m.repo) - 1
	if index != lastIndex {
		m.repo[index] = m.repo[lastIndex]

		movedEntity := &m.repo[index]
		m.uuidIndex[movedEntity.UUID] = index
		m.updateIndexPosition(m.batchIndex, movedEntity.Batch, lastIndex, index)
		m.updateIndexPosition(m.materialIndex, movedEntity.Material, lastIndex, index)
	}

	m.repo = m.repo[:lastIndex]

	return nil
}

func (m *Map) removeFromIndex(index map[uuid.UUID][]int, key uuid.UUID, pos int) {
	indices := index[key]
	for i, idx := range indices {
		if idx == pos {
			index[key] = append(indices[:i], indices[i+1:]...)
			break
		}
	}
	if len(index[key]) == 0 {
		delete(index, key)
	}
}

func (m *Map) updateIndexPosition(index map[uuid.UUID][]int, key uuid.UUID, oldPos, newPos int) {
	indices := index[key]
	for i, idx := range indices {
		if idx == oldPos {
			indices[i] = newPos
			break
		}
	}
}

func someThen[F any](dst *F, src opt.Opt[F]) {
	val, ok := src.Unwrap()
	if !ok {
		return
	}

	*dst = val
}

func clamp[T cmp.Ordered](mn, val, mx T) T {
	return min(max(mn, val), mx)
}
