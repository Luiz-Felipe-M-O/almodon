package userrepo

import (
	"cmp"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
	"unsafe"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	userpkg "github.com/alan-b-lima/almodon/internal/domain/user"

	repo "github.com/alan-b-lima/almodon/internal/support/repository"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Map struct {
	uuid  repo.Index[uuid.UUID, int]
	siape repo.Index[string, int]

	repo []user.Entity
	mu   sync.RWMutex

	datapath string
}

func NewMap() user.Repository {
	repo := Map{
		uuid:  make(repo.Index[uuid.UUID, int]),
		siape: make(repo.Index[string, int]),
	}

	return &repo
}

func NewPersistantMap(datapath string) (user.Repository, error) {
	repo := Map{
		uuid:     make(repo.Index[uuid.UUID, int]),
		siape:    make(repo.Index[string, int]),
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

	var repo []entity
	if err := json.NewDecoder(f).Decode(&repo); err != nil {
		return err
	}

	m.repo = entity_from_json(repo)
	for i, record := range m.repo {
		m.uuid[record.UUID] = i
		m.siape[record.SIAPE] = i
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

func (m *Map) List(offset, limit int) (user.Entities, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	lo := clamp(0, offset, len(m.repo))
	hi := clamp(0, offset+limit, len(m.repo))

	if lo >= hi {
		return user.Entities{
			Records:      []user.Entity{},
			TotalRecords: len(m.repo),
		}, nil
	}

	res := make([]user.Entity, hi-lo)
	copy(res, m.repo[lo:hi])

	return user.Entities{
		Offset:       lo,
		Length:       len(res),
		Records:      res,
		TotalRecords: len(m.repo),
	}, nil
}

func (m *Map) Get(uuid uuid.UUID) (user.Entity, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.uuid[uuid]
	if !in {
		return user.Entity{}, user.ErrUserNotFound
	}

	return m.repo[index], nil
}

func (m *Map) GetBySIAPE(siape string) (user.Entity, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.siape[siape]
	if !in {
		return user.Entity{}, user.ErrUserNotFound
	}

	return m.repo[index], nil
}

func (m *Map) Create(user user.Entity) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	if _, in := m.siape[user.SIAPE]; in {
		return userpkg.ErrSiapeTaken
	}

	m.uuid.Set(user.UUID, len(m.repo))
	m.siape.Set(user.SIAPE, len(m.repo))
	m.repo = append(m.repo, user)

	return nil
}

func (m *Map) Patch(uuid uuid.UUID, user user.PartialEntity) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuid[uuid]
	if !in {
		return userpkg.ErrUserNotFound
	}

	u := &m.repo[index]

	if role, ok := user.Role.Unwrap(); ok {
		if role != u.Role && u.Role == auth.Chief && !enough_chiefs(m) {
			return userpkg.ErrNotEnoughChiefs
		} else {
			u.Role = role
		}
	}

	repo.SomeThen(&u.Name, user.Name)
	repo.SomeThen(&u.Email, user.Email)
	repo.SomeThen(&u.Password, user.Password)

	u.Updated = user.Updated
	return nil
}

func (m *Map) Delete(uuid uuid.UUID) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuid.Get(uuid)
	if !in {
		return nil
	}

	u := &m.repo[index]
	if u.Role == auth.Chief && !enough_chiefs(m) {
		return user.ErrNotEnoughChiefs
	}

	m.uuid.Del(u.UUID)
	m.siape.Del(u.SIAPE)

	last := len(m.repo) - 1
	if index != last {
		last := &m.repo[last]
		m.repo[index] = *last
		m.uuid.Set(last.UUID, index)
		m.siape.Set(last.SIAPE, index)
	}

	m.repo = m.repo[:last]
	return nil
}

func enough_chiefs(m *Map) bool {
	var found bool
	for _, user := range m.repo {
		if user.Role != auth.Chief {
			continue
		}

		if !found {
			found = true
		} else {
			return true
		}
	}

	return false
}

func clamp[T cmp.Ordered](mn, val, mx T) T {
	return min(max(mn, val), mx)
}

type entity struct {
	UUID     uuid.UUID `json:"uuid"`
	SIAPE    string    `json:"siape"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password pwd       `json:"password"`
	Role     auth.Role `json:"role"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

func json_for_entity(entities []user.Entity) []entity {
	return unsafe.Slice((*entity)(unsafe.Pointer(unsafe.SliceData(entities))), len(entities))
}

func entity_from_json(entities []entity) []user.Entity {
	return unsafe.Slice((*user.Entity)(unsafe.Pointer(unsafe.SliceData(entities))), len(entities))
}

type pwd [60]byte

func (v pwd) MarshalJSON() ([]byte, error) {
	buf := base64.StdEncoding.AppendEncode(nil, v[:])
	return fmt.Appendf(nil, "%+q", buf), nil
}

func (v *pwd) UnmarshalJSON(buf []byte) error {
	if len(buf) < 2 || buf[0] != '"' || buf[len(buf)-1] != '"' {
		return fmt.Errorf("invalid password encoding")
	}

	_, err := base64.StdEncoding.Decode(v[:], buf[1:len(buf)-1])
	return err
}
