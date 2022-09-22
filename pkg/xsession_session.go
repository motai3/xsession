package pkg

import (
	"context"
	"sync"
	"time"
)

type Session struct {
	id      string          // Session id.
	ctx     context.Context // Context for current session, note that: one session one context.
	data    *sync.Map       // Session data.
	dirty   bool            // Used to mark session is modified.
	start   bool            // Used to mark session is started.
	manager *Manager        // Parent manager.

	// idFunc is a callback function used for creating custom session id.
	// This is called if session id is empty ever when session starts.
	idFunc func(ttl time.Duration) (id string)
}
