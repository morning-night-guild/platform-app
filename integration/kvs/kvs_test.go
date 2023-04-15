package kvs_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/integration/helper"
	"github.com/morning-night-guild/platform-app/internal/driver/redis"
)

func TestKVS(t *testing.T) {
	t.Parallel()

	url := helper.GetRedisURL(t)

	type KV struct {
		Key   string
		Value string
	}

	t.Run("保存と取得ができる", func(t *testing.T) {
		t.Parallel()

		rds, err := redis.NewRedis(url)
		if err != nil {
			t.Fatal("failed to connect to redis")
		}

		cache, err := redis.New[KV]().KVS("test", rds)
		if err != nil {
			t.Fatal("failed to create kvs")
		}

		ctx := context.Background()

		kv := KV{
			Key:   uuid.New().String(),
			Value: uuid.New().String(),
		}

		if err := cache.Set(ctx, kv.Key, kv, time.Hour); err != nil {
			t.Errorf("failed to set: %v", err)
		}

		got, err := cache.Get(ctx, kv.Key)
		if err != nil {
			t.Fatalf("failed to get: %v", err)
		}

		if !reflect.DeepEqual(got, kv) {
			t.Errorf("got = %v, want %v", got, kv)
		}
	})

	t.Run("削除ができる", func(t *testing.T) {
		t.Parallel()

		rds, err := redis.NewRedis(url)
		if err != nil {
			t.Fatal("failed to connect to redis")
		}

		cache, err := redis.New[KV]().KVS("test", rds)
		if err != nil {
			t.Fatal("failed to create kvs")
		}

		ctx := context.Background()

		kv := KV{
			Key:   uuid.New().String(),
			Value: uuid.New().String(),
		}

		if err := cache.Set(ctx, kv.Key, kv, time.Hour); err != nil {
			t.Fatalf("failed to set: %v", err)
		}

		if err := cache.Del(ctx, kv.Key); err != nil {
			t.Errorf("failed to delete: %v", err)
		}

		if _, err := cache.Get(ctx, kv.Key); err == nil {
			t.Error("failed to delete")
		}
	})

	t.Run("TTLが設定できる", func(t *testing.T) {
		t.Parallel()

		rds, err := redis.NewRedis(url)
		if err != nil {
			t.Fatal("failed to connect to redis")
		}

		cache, err := redis.New[KV]().KVS("test", rds)
		if err != nil {
			t.Fatal("failed to create kvs")
		}

		ctx := context.Background()

		kv := KV{
			Key:   uuid.New().String(),
			Value: uuid.New().String(),
		}

		if err := cache.Set(ctx, kv.Key, kv, time.Second); err != nil {
			t.Fatalf("failed to set: %v", err)
		}

		time.Sleep(1100 * time.Millisecond)

		if _, err := cache.Get(ctx, kv.Key); err == nil {
			t.Error("failed to delete")
		}
	})
}
