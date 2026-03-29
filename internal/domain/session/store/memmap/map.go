package sessionrepo

import (
	"context"
	"sync"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/support"
	"github.com/alan-b-lima/almodon/internal/support/store"

	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/pkg/timeout"
)

type Map struct {
	uuid  store.Index[uuid.UUID]
	user  store.Index[uuid.UUID]
	sched *timeout.Timeout

	repo []session.Record

	mu sync.RWMutex
}

func NewMemMap(t *timeout.Timeout) (session.Store, error) {
	if t == nil {
		return nil, support.ErrNilPointer
	}

	repo := Map{
		uuid:  make(store.Index[uuid.UUID]),
		user:  make(store.Index[uuid.UUID]),
		sched: t,
	}

	return &repo, nil
}

func (m *Map) Get(ctx context.Context, uuid uuid.UUID) (session.Record, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	return m.get(ctx, uuid)
}

func (m *Map) get(_ context.Context, uuid uuid.UUID) (session.Record, error) {
	index, in := m.uuid.Get(uuid)
	if !in {
		return session.Record{}, session.ErrNotFound
	}

	s := m.repo[index]
	if time.Now().After(s.Expires) {
		return session.Record{}, session.ErrNotFound
	}

	return s, nil
}

func (m *Map) Create(_ context.Context, s session.CreateRecord) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	if index, in := m.user.Get(s.User); in {
		s := m.repo[index]
		m.delete(s.UUID)
	}

	m.uuid.Set(s.UUID, len(m.repo))
	m.user.Set(s.User, len(m.repo))
	m.repo = append(m.repo, session.Record(s))

	m.post(s.UUID, s.Expires)
	return nil
}

func (m *Map) Update(_ context.Context, uuid uuid.UUID, req session.UpdateRecord) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuid.Get(uuid)
	if !in {
		return session.ErrNotFound
	}

	s := &m.repo[index]
	s.Expires = req.Expires

	m.post(s.UUID, s.Expires)
	return nil
}

func (m *Map) Delete(_ context.Context, uuid uuid.UUID) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	return m.delete(uuid)
}

func (m *Map) delete(uuid uuid.UUID) error {
	index, in := m.uuid.Get(uuid)
	if !in {
		return nil
	}

	s := &m.repo[index]

	m.uuid.Del(s.UUID)
	m.user.Del(s.User)

	last := len(m.repo) - 1
	if index != last {
		last := &m.repo[last]
		m.repo[index] = *last
		m.uuid.Set(last.UUID, index)
		m.user.Set(last.User, index)
	}

	m.repo = m.repo[:last]

	return nil
}

func (m *Map) delete_expired(ctx context.Context, uuid uuid.UUID) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	s, err := m.get(ctx, uuid)
	if err != nil {
		if err == session.ErrNotFound {
			return nil
		}

		return err
	}

	if time.Now().Before(s.Expires) {
		return nil // too early
	}

	return m.delete(uuid)
}

func (m *Map) RunTx(ctx context.Context, proc func(session.Store) error) error {
	return proc(m)
}

func (m *Map) post(session uuid.UUID, expires time.Time) {
	m.sched.Post(
		func() { m.delete_expired(context.Background(), session) },
		expires,
	)
}
