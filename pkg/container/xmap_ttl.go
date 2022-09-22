package container

import (
	"sync"
	"time"
)

const (
	defaultTTl = time.Second * 60 * 60
)

type TTLMap struct {
	mu sync.Mutex

	dirty map[interface{}]interface{}

	ttlDirty map[interface{}]time.Time
}

func NewTTLMap() *TTLMap {
	return &TTLMap{
		dirty:    make(map[interface{}]interface{}),
		ttlDirty: make(map[interface{}]time.Time),
	}
}

func (m *TTLMap) load(key interface{}) (value interface{}, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.ttlDirty[key]
	if ok {
		if t.Before(time.Now()) {
			delete(m.ttlDirty, key)
			delete(m.dirty, key)
		}
		i, ok := m.dirty[key]
		return i, ok
	}
	return t, ok
}

func (m *TTLMap) Store(key, value interface{}, ttl time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ttlDirty[key] = time.Now().Add(ttl)
	m.dirty[key] = value
}
