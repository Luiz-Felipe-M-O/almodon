package materialrepo

import (
	"cmp"
	"encoding/json"
	"os"
	"sync"

	"github.com/alan-b-lima/almodon/internal/domain/material"
	repo "github.com/alan-b-lima/almodon/internal/support/repository"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Map struct {
	uuid    repo.Index[uuid.UUID, int]
	siads   repo.SliceIndex[string, int]
	catmat  repo.SliceIndex[string, int]
	ecampus repo.SliceIndex[string, int]

	repo []material.Entity
	mu   sync.RWMutex

	datapath string
}

func NewMap() material.Repository {
	repo := Map{
		uuid:    make(repo.Index[uuid.UUID, int]),
		siads:   make(repo.SliceIndex[string, int]),
		catmat:  make(repo.SliceIndex[string, int]),
		ecampus: make(repo.SliceIndex[string, int]),
	}

	return &repo
}

func NewPersistentMap(datapath string) (material.Repository, error) {
	repo := Map{
		uuid:     make(repo.Index[uuid.UUID, int]),
		siads:    make(repo.SliceIndex[string, int]),
		catmat:   make(repo.SliceIndex[string, int]),
		ecampus:  make(repo.SliceIndex[string, int]),
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

		m.siads.Add(record.SIADS, i)
		m.catmat.Add(record.CATMAT, i)
		m.ecampus.Add(record.ECampus, i)
	}

	return nil
}

func (m *Map) Close() error {
	if m.datapath == "" {
		return nil
	}

	defer m.mu.Unlock()
	m.mu.Lock()

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

func (m *Map) List(offset, limit int) (material.Entities, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	lo := clamp(0, offset, len(m.repo))
	hi := clamp(0, offset+limit, len(m.repo))

	if lo >= hi {
		return material.Entities{
			Records:      []material.Entity{},
			TotalRecords: len(m.repo),
		}, nil
	}

	res := make([]material.Entity, hi-lo)
	copy(res, m.repo[lo:hi])

	return material.Entities{
		Offset:       lo,
		Length:       len(res),
		Records:      res,
		TotalRecords: len(m.repo),
	}, nil
}

func (m *Map) ListBySIADS(siads string) (material.Entities, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	indices, in := m.siads.Get(siads)
	if !in || len(indices) == 0 {
		return material.Entities{
			Records:      []material.Entity{},
			TotalRecords: 0,
		}, nil
	}

	res := make([]material.Entity, len(indices))
	for i, idx := range indices {
		res[i] = m.repo[idx]
	}

	return material.Entities{
		Offset:       0,
		Length:       len(res),
		Records:      res,
		TotalRecords: len(res),
	}, nil
}

func (m *Map) ListByCATMAT(catmat string) (material.Entities, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	indices, in := m.catmat.Get(catmat)
	if !in || len(indices) == 0 {
		return material.Entities{
			Records:      []material.Entity{},
			TotalRecords: 0,
		}, nil
	}

	res := make([]material.Entity, len(indices))
	for i, idx := range indices {
		res[i] = m.repo[idx]
	}

	return material.Entities{
		Offset:       0,
		Length:       len(res),
		Records:      res,
		TotalRecords: len(res),
	}, nil
}

func (m *Map) ListByECampus(ecampus string) (material.Entities, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	indexes, in := m.ecampus.Get(ecampus)
	if !in || len(indexes) == 0 {
		return material.Entities{
			Records:      []material.Entity{},
			TotalRecords: 0,
		}, nil
	}

	res := make([]material.Entity, len(indexes))
	for i, index := range indexes {
		res[i] = m.repo[index]
	}

	return material.Entities{
		Offset:       0,
		Length:       len(res),
		Records:      res,
		TotalRecords: len(res),
	}, nil
}

func (m *Map) Get(uuid uuid.UUID) (material.Entity, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.uuid[uuid]
	if !in {
		return material.Entity{}, xerrors.ErrMaterialNotFound
	}

	return m.repo[index], nil
}

func (m *Map) Create(material material.Entity) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index := len(m.repo)

	m.uuid.Set(material.UUID, index)
	m.siads.Add(material.SIADS, index)
	m.catmat.Add(material.CATMAT, index)
	m.ecampus.Add(material.ECampus, index)

	m.repo = append(m.repo, material)
	return nil
}

func (m *Map) Patch(uuid uuid.UUID, partial material.PartialEntity) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuid[uuid]
	if !in {
		return xerrors.ErrMaterialNotFound
	}

	material := &m.repo[index]

	if siads, ok := partial.SIADS.Unwrap(); ok {
		m.siads.Del(material.SIADS, index)
		m.siads.Add(siads, index)
		material.SIADS = siads
	}

	if catmat, ok := partial.CATMAT.Unwrap(); ok {
		m.catmat.Del(material.CATMAT, index)
		m.catmat.Add(catmat, index)
		material.CATMAT = catmat
	}

	if ecampus, ok := partial.ECampus.Unwrap(); ok {
		m.ecampus.Del(material.ECampus, index)
		m.ecampus.Add(ecampus, index)
		material.ECampus = ecampus
	}

	repo.SomeThen(&material.Name, partial.Name)
	repo.SomeThen(&material.Description, partial.Description)
	repo.SomeThen(&material.Unit, partial.Unit)
	repo.SomeThen(&material.MinQuantity, partial.MinQuantity)
	material.Updated = partial.Updated

	return nil
}

func (m *Map) Delete(uuid uuid.UUID) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuid[uuid]
	if !in {
		return nil
	}

	mat := &m.repo[index]

	delete(m.uuid, mat.UUID)
	m.siads.Del(mat.SIADS, index)
	m.catmat.Del(mat.CATMAT, index)
	m.ecampus.Del(mat.ECampus, index)

	lastIndex := len(m.repo) - 1
	if index != lastIndex {
		m.repo[index] = m.repo[lastIndex]

		lastMat := &m.repo[index]

		m.uuid.Set(lastMat.UUID, index)
		m.siads.Swap(lastMat.SIADS, lastIndex, index)
		m.catmat.Swap(lastMat.CATMAT, lastIndex, index)
		m.ecampus.Swap(lastMat.ECampus, lastIndex, index)
	}

	m.repo = m.repo[:lastIndex]

	return nil
}

func clamp[T cmp.Ordered](mn, val, mx T) T {
	return min(max(mn, val), mx)
}
