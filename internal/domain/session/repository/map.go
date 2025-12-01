package sessionrepo

import (
	"sync"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/session"
	repo "github.com/alan-b-lima/almodon/internal/support/repository"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/heap"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Map struct {
	uuid        repo.Index[uuid.UUID, int]
	user        repo.Index[uuid.UUID, int]
	expiresHeap sleepqueue

	repo []session.Entity
	mu   sync.RWMutex
}

func NewMap() session.Repository {
	repo := Map{
		uuid: make(repo.Index[uuid.UUID, int]),
		user: make(repo.Index[uuid.UUID, int]),
		expiresHeap: sleepqueue{
			new:    make(chan ess, 64),
			cancel: make(chan struct{}, 1),
		},
	}

	go flush(&repo)

	return &repo
}

func (m *Map) Get(uuid uuid.UUID) (session.Entity, error) {
	defer m.mu.RUnlock()
	m.mu.RLock()

	index, in := m.uuid.Get(uuid)
	if !in {
		return session.Entity{}, xerrors.ErrSessionNotFound
	}

	s := m.repo[index]
	if time.Now().After(s.Expires) {
		return session.Entity{}, xerrors.ErrSessionNotFound
	}

	return m.repo[index], nil
}

func (m *Map) Create(session session.Entity) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	if index, in := m.user.Get(session.User); in {
		s := m.repo[index]
		m.delete(s.UUID)
	}

	m.uuid.Set(session.UUID, len(m.repo))
	m.user.Set(session.User, len(m.repo))
	m.repo = append(m.repo, session)

	m.expiresHeap.new <- ess{session.UUID, session.Expires}

	return nil
}

func (m *Map) Update(uuid uuid.UUID, expires time.Time) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	index, in := m.uuid.Get(uuid)
	if !in {
		return xerrors.ErrSessionNotFound
	}

	s := &m.repo[index]
	s.Expires = expires

	m.expiresHeap.new <- ess{s.UUID, expires}

	return nil
}

func (m *Map) Delete(uuid uuid.UUID) error {
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

func flush(m *Map) {
	h := m.expiresHeap

	for {
		var after <-chan time.Time
		if h.heap.Len() > 0 {
			delay := time.Until(h.heap.Peek().expires)
			after = time.After(delay)
		}

		select {
		case <-h.cancel:
			return

		case es := <-h.new:
			h.heap.Push(es)

		case <-after:
			es := h.heap.Pop()
			m.Delete(es.session)
		}
	}
}

type sleepqueue struct {
	heap   heap.Heap[ess]
	new    chan ess
	cancel chan struct{}
}

type ess struct {
	session uuid.UUID
	expires time.Time
}

func (o0 ess) Less(o1 ess) bool { return o0.expires.Before(o1.expires) }
