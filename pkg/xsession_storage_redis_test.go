package xsession_test

import (
	"context"
	"github.com/go-redis/redis"
	xsession "github.com/motai3/xsession/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_StorageRedis(t *testing.T) {
	redis := redis.NewClient(&redis.Options{
		Addr:     "43.138.13.24:6379",
		Password: "",
		DB:       0,
	})
	sessionId := ""
	manager := xsession.New(time.Second*60, xsession.NewStorageRedis(redis))
	t.Run("RedisSet", func(t *testing.T) {
		s := manager.New(context.TODO())

		defer s.Close()
		s.Set("k1", "v1")
		s.Set("k2", "v2")
		m := make(map[string]interface{})
		m["k3"] = "v3"
		m["k4"] = "v4"
		s.Sets(m)
		//assert.True(t, s.IsDirty())
		sessionId = s.Id()
	})
	t.Run("RedisGet", func(t *testing.T) {
		s := manager.New(context.TODO(), sessionId)
		assert.Equal(t, s.Get("k1"), "v1")
		assert.Equal(t, s.Get("k2"), "v2")
		assert.Equal(t, s.Size(), 4)
	})
}

func Test_StorageMemory(t *testing.T) {
	sessionId := ""
	manager := xsession.New(time.Second*60, xsession.NewStorageMemory())
	t.Run("MemorySet", func(t *testing.T) {
		s := manager.New(context.TODO())
		defer s.Close()
		s.Set("k1", "v1")
		s.Set("k2", "v2")
		m := make(map[string]interface{})
		m["k3"] = "v3"
		m["k4"] = "v4"
		s.Sets(m)
		//assert.True(t, s.IsDirty())
		sessionId = s.Id()
	})
	t.Run("MemoryGet", func(t *testing.T) {
		s := manager.New(context.TODO(), sessionId)
		assert.Equal(t, s.Get("k1"), "v1")
		assert.Equal(t, s.Get("k4"), "v4")
		assert.Equal(t, s.Size(), 4)
	})
}
