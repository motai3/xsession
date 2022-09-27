package pkg

import (
	"context"
	"sync"
	"time"
	"xsession/pkg/container"
)

type Manager struct {
	ttl     time.Duration
	storage Storage

	// 存储Seervice session信息
	sessionData *container.TTLMap
}

func Mew(ttl time.Duration, storage ...Storage) *Manager {
	m := &Manager{
		ttl:         ttl,
		sessionData: container.NewTTLMap(),
	}
	if len(storage) > 0 && storage != nil {
		m.storage = storage[0]
	} else {
		m.storage = NewStorageMemory()
	}
	return m
}

func (m *Manager) New(ctx context.Context, sessionId ...string) *Session {
	var id string
	if len(sessionId) > 0 && sessionId[0] != "" {
		id = sessionId[0]
	}
	return &Session{
		id:      id,
		ctx:     ctx,
		manager: m,
	}
}

func (m *Manager) SetStorage(storage Storage) {
	m.storage = storage
}

func (m *Manager) SetTTL(ttl time.Duration) {
	m.ttl = ttl
}

func (m *Manager) TTL() time.Duration {
	return m.ttl
}

func (m *Manager) UpdateSessionTTL(sessionId string, data *sync.Map) {
	m.sessionData.Store(sessionId, data, m.ttl)
}
