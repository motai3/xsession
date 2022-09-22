package pkg

import (
	"context"
	"sync"
	"time"
)

type Storage interface {
	New(ctx context.Context, ttl time.Duration) (id string, err error)

	Get(ctx context.Context, id string, key string) (value interface{}, err error)

	GetMap(ctx context.Context, id string) (data map[string]interface{}, err error)

	GetSize(ctx context.Context, id string) (size int, err error)

	Set(ctx context.Context, id string, key string, value interface{}, ttl time.Duration) error

	SetMap(ctx context.Context, id string, data map[string]interface{}, ttl time.Duration) error

	Remove(ctx context.Context, id string, key string) error

	RemoveAll(ctx context.Context, id string) error

	GetSession(ctx context.Context, id string, ttl time.Duration, data *sync.Map) (*sync.Map, error)

	SetSession(ctx context.Context, id string, data *sync.Map, ttl time.Duration) error

	UpdateTTL(ctx context.Context, id string, ttl time.Duration) error
}
