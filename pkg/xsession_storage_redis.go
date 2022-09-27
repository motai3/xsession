package xsession

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"sync"
	"time"
)

type StorageRedisHashTable struct {
	redis  *redis.Client
	prefix string
}

func NewStorageRedis(redis *redis.Client, prefix ...string) *StorageRedisHashTable {
	if redis == nil {
		panic("redis instance for storage cannot be empty")
		return nil
	}
	s := &StorageRedisHashTable{
		redis: redis,
	}
	if len(prefix) > 0 && prefix[0] != "" {
		s.prefix = prefix[0]
	}
	return s
}

func (s *StorageRedisHashTable) New(ctx context.Context, ttl time.Duration) (id string, err error) {
	return "", ErrorDisabled
}

func (s *StorageRedisHashTable) Get(ctx context.Context, id string, key string) (value interface{}, err error) {
	value, err = s.redis.HGet(s.key(id), key).Result()
	if value != nil {
		value = value.(string)
	}
	return
}

func (s *StorageRedisHashTable) GetMap(ctx context.Context, id string) (data map[string]interface{}, err error) {
	r, err := s.redis.HGetAll(s.key(id)).Result()
	if err != nil {
		return nil, err
	}
	data = make(map[string]interface{})
	for k, v := range r {
		data[k] = v
	}
	return data, err
}

func (s *StorageRedisHashTable) GetSize(ctx context.Context, id string) (size int, err error) {
	r, err := s.redis.HLen(s.key(id)).Result()
	if err != nil {
		return -1, err
	}
	strSize := strconv.FormatInt(r, 10)
	res, _ := strconv.Atoi(strSize)
	return res, nil
}

func (s *StorageRedisHashTable) Set(ctx context.Context, id string, key string, value interface{}, ttl time.Duration) error {
	err := s.redis.HSet(s.key(id), key, value).Err()
	return err
}

func (s *StorageRedisHashTable) SetMap(ctx context.Context, id string, data map[string]interface{}, ttl time.Duration) error {
	err := s.redis.HMSet(s.key(id), data).Err()
	return err
}

func (s *StorageRedisHashTable) Remove(ctx context.Context, id string, key string) error {
	err := s.redis.HDel(s.key(id), key).Err()
	return err
}

func (s *StorageRedisHashTable) RemoveAll(ctx context.Context, id string) error {
	err := s.redis.Del(s.key(id)).Err()
	return err
}

func (s *StorageRedisHashTable) GetSession(ctx context.Context, id string, ttl time.Duration, data *sync.Map) (*sync.Map, error) {
	fmt.Sprintf("StorageRedis.GetSession: %s, %v", id, ttl)
	redisData, err := s.GetMap(ctx, id)
	if err != nil {
		return nil, err
	}
	if redisData == nil {
		return nil, err
	}
	var newData sync.Map
	for k, v := range redisData {
		newData.Store(k, v)
	}

	return &newData, err
}

func (s *StorageRedisHashTable) SetSession(ctx context.Context, id string, data *sync.Map, ttl time.Duration) error {
	fmt.Sprintf("StorageRedis.SetSession: %s, %v", id, ttl)
	redisData := make(map[string]interface{})
	data.Range(func(key, value interface{}) bool {
		redisData[key.(string)] = value
		return true
	})
	err := s.SetMap(ctx, id, redisData, ttl)
	if err != nil {
		return err
	}
	err = s.redis.Expire(s.key(id), ttl).Err()
	return err
}

func (s *StorageRedisHashTable) UpdateTTL(ctx context.Context, id string, ttl time.Duration) error {
	err := s.redis.Expire(s.key(id), ttl).Err()
	return err
}

func (s *StorageRedisHashTable) key(id string) string {
	return s.prefix + id
}
