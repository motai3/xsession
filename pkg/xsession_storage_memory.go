package xsession

import (
	"context"
	"sync"
	"time"
)

type StorageMemory struct{}

func NewStorageMemory() *StorageMemory {
	return &StorageMemory{}
}

func (s *StorageMemory) New(ctx context.Context, ttl time.Duration) (id string, err error) {
	return "", ErrorDisabled

}

func (s *StorageMemory) Get(ctx context.Context, id string, key string) (value interface{}, err error) {
	return nil, ErrorDisabled
}

func (s *StorageMemory) GetMap(ctx context.Context, id string) (data map[string]interface{}, err error) {
	return nil, ErrorDisabled
}

func (s *StorageMemory) GetSize(ctx context.Context, id string) (size int, err error) {
	return -1, ErrorDisabled
}

func (s *StorageMemory) Set(ctx context.Context, id string, key string, value interface{}, ttl time.Duration) error {
	return ErrorDisabled
}

func (s *StorageMemory) SetMap(ctx context.Context, id string, data map[string]interface{}, ttl time.Duration) error {
	return ErrorDisabled
}

func (s *StorageMemory) Remove(ctx context.Context, id string, key string) error {
	return ErrorDisabled
}

func (s *StorageMemory) RemoveAll(ctx context.Context, id string) error {
	return ErrorDisabled
}

func (s *StorageMemory) GetSession(ctx context.Context, id string, ttl time.Duration, data *sync.Map) (*sync.Map, error) {
	return data, nil
}

func (s *StorageMemory) SetSession(ctx context.Context, id string, data *sync.Map, ttl time.Duration) error {
	return nil
}

func (s *StorageMemory) UpdateTTL(ctx context.Context, id string, ttl time.Duration) error {
	return nil
}
