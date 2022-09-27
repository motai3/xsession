package pkg

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Session struct {
	id      string
	ctx     context.Context
	data    *sync.Map
	dirty   bool
	start   bool
	manager *Manager
	idFunc  func(ttl time.Duration) (id string)
}

func (s *Session) init() {
	if s.start {
		return
	}
	var err error
	if s.id != "" {
		if r, _ := s.manager.sessionData.Get(s.id); r != nil {
			s.data = r.(*sync.Map)
			fmt.Printf("session init data : %s", s.data)
		}
		if s.manager.storage != nil {
			if s.data, err = s.manager.storage.GetSession(s.ctx, s.id, s.manager.ttl, s.data); err != nil && err != ErrorDisabled {
				fmt.Errorf("session restoring failed for id '%s': %v", s.id, err)
			}
		}
	}
	if s.id == "" && s.idFunc != nil {
		s.id = s.idFunc(s.manager.ttl)
	}
	if s.id == "" {
		s.id, err = s.manager.storage.New(s.ctx, s.manager.ttl)
		if err != nil && err != ErrorDisabled {
			fmt.Errorf("create session id failed")
		}
	}
	if s.id == "" {
		s.id = NewSessionId()
	}
	s.start = true
}

func (s *Session) Close() {
	if s.start && s.id != "" {
		size := 0
		s.data.Range(func(key, value interface{}) bool {
			size++
			return true
		})
		if s.dirty {
			if err := s.manager.storage.SetSession(s.ctx, s.id, s.data, s.manager.ttl); err != nil {
				return
			}
		} else if size > 0 {
			if err := s.manager.storage.UpdateTTL(s.ctx, s.id, s.manager.ttl); err != nil {
				return
			}
		}
		if s.dirty || size > 0 {
			s.manager.UpdateSessionTTL(s.id, s.data)
		}
	}

}

func (s *Session) Set(key string, value interface{}) error {
	s.init()
	if err := s.manager.storage.Set(s.ctx, s.id, key, value, s.manager.ttl); err != nil {
		if err == ErrorDisabled {
			s.data.Store(key, value)
		} else {
			return err
		}
	}
	s.dirty = true
	return nil
}

func (s *Session) Sets(data map[string]interface{}) error {
	return s.SetMap(data)
}

func (s *Session) SetMap(data map[string]interface{}) error {
	s.init()
	if err := s.manager.storage.SetMap(s.ctx, s.id, data, s.manager.ttl); err != nil {
		if err == ErrorDisabled {
			for k, v := range data {
				s.data.Store(k, v)
			}
		} else {
			return err
		}
	}
	s.dirty = true
	return nil
}

func (s *Session) Remove(keys ...string) error {
	if s.id == "" {
		return nil
	}
	s.init()
	for _, key := range keys {
		if err := s.manager.storage.Remove(s.ctx, s.id, key); err != nil {
			if err == ErrorDisabled {
				s.data.Delete(key)
			} else {
				return err
			}
		}
	}
	s.dirty = true
	return nil
}

func (s *Session) Clear() error {
	return s.RemoveAll()
}

// RemoveAll deletes all key-value pairs from this session.
func (s *Session) RemoveAll() error {
	if s.id == "" {
		return nil
	}
	s.init()
	if err := s.manager.storage.RemoveAll(s.ctx, s.id); err != nil {
		if err == ErrorDisabled {
			s.data.Range(func(key, value interface{}) bool {
				s.data.Delete(key)
				return true
			})
		} else {
			return err
		}
	}
	s.dirty = true
	return nil
}

func (s *Session) Id() string {
	s.init()
	return s.id
}

// SetId sets custom session before session starts.
// It returns error if it is called after session starts.
func (s *Session) SetId(id string) error {
	if s.start {
		return errors.New("session already started")
	}
	s.id = id
	return nil
}

// SetIdFunc sets custom session id creating function before session starts.
// It returns error if it is called after session starts.
func (s *Session) SetIdFunc(f func(ttl time.Duration) string) error {
	if s.start {
		return errors.New("session already started")
	}
	s.idFunc = f
	return nil
}

func (s *Session) Map() map[string]interface{} {
	if s.id != "" {
		s.init()
		data, err := s.manager.storage.GetMap(s.ctx, s.id)
		if err != nil && err != ErrorDisabled {
			fmt.Errorf("error is : %s", err)
		}
		if data != nil {
			return data
		}
		dataMap := make(map[string]interface{})
		s.data.Range(func(key, value interface{}) bool {
			dataMap[key.(string)] = value
			return true
		})
		return dataMap
	}
	return nil
}

func (s *Session) Size() int {
	if s.id != "" {
		s.init()
		size, err := s.manager.storage.GetSize(s.ctx, s.id)
		if err != nil && err != ErrorDisabled {
			fmt.Errorf("error is : %s", err)
		}
		if size >= 0 {
			return size
		}
		len := 0
		s.data.Range(func(key, value interface{}) bool {
			len++
			return true
		})
		return len
	}
	return 0
}

func (s *Session) Contains(key string) bool {
	s.init()
	return s.Get(key) != nil
}

func (s *Session) IsDirty() bool {
	return s.dirty
}

func (s *Session) Get(key string, def ...interface{}) interface{} {
	if s.id == "" {
		return nil
	}
	s.init()
	v, err := s.manager.storage.Get(s.ctx, s.id, key)
	if err != nil && err != ErrorDisabled {
		fmt.Errorf("error is : %s", err)
	}
	if v != nil {
		return v
	}
	if v, _ := s.data.Load(key); v != nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}
