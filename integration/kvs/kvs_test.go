package kvs_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/integration/helper"
	"github.com/morning-night-guild/platform-app/internal/domain/cache"
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

		kvs, err := redis.New[KV]().KVS("test", rds)
		if err != nil {
			t.Fatal("failed to create kvs")
		}

		ctx := context.Background()

		kv := KV{
			Key:   uuid.New().String(),
			Value: uuid.New().String(),
		}

		if err := kvs.Set(ctx, kv.Key, kv, time.Hour); err != nil {
			t.Errorf("failed to set: %v", err)
		}

		got, err := kvs.Get(ctx, kv.Key)
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

		kvs, err := redis.New[KV]().KVS("test", rds)
		if err != nil {
			t.Fatal("failed to create kvs")
		}

		ctx := context.Background()

		kv := KV{
			Key:   uuid.New().String(),
			Value: uuid.New().String(),
		}

		if err := kvs.Set(ctx, kv.Key, kv, time.Hour); err != nil {
			t.Fatalf("failed to set: %v", err)
		}

		if err := kvs.Del(ctx, kv.Key); err != nil {
			t.Errorf("failed to delete: %v", err)
		}

		if _, err := kvs.Get(ctx, kv.Key); err == nil {
			t.Error("failed to delete")
		}
	})

	t.Run("TTLが設定できる", func(t *testing.T) {
		t.Parallel()

		rds, err := redis.NewRedis(url)
		if err != nil {
			t.Fatal("failed to connect to redis")
		}

		kvs, err := redis.New[KV]().KVS("test", rds)
		if err != nil {
			t.Fatal("failed to create kvs")
		}

		ctx := context.Background()

		kv := KV{
			Key:   uuid.New().String(),
			Value: uuid.New().String(),
		}

		if err := kvs.Set(ctx, kv.Key, kv, time.Second); err != nil {
			t.Fatalf("failed to set: %v", err)
		}

		time.Sleep(1100 * time.Millisecond)

		if _, err := kvs.Get(ctx, kv.Key); err == nil {
			t.Error("failed to delete")
		}
	})

	t.Run("トランザクション保存ができる", func(t *testing.T) {
		t.Parallel()

		rds, err := redis.NewRedis(url)
		if err != nil {
			t.Fatal("failed to connect to redis")
		}

		kvs, err := redis.New[KV]().KVS("test", rds)
		if err != nil {
			t.Fatal("failed to create kvs")
		}

		ctx := context.Background()

		kv0 := KV{
			Key:   uuid.New().String(),
			Value: uuid.New().String(),
		}

		if err := kvs.Set(ctx, kv0.Key, kv0, time.Hour); err != nil {
			t.Fatalf("failed to set: %v", err)
		}

		kv1 := KV{
			Key:   uuid.New().String(),
			Value: uuid.New().String(),
		}

		kv2 := KV{
			Key:   uuid.New().String(),
			Value: uuid.New().String(),
		}

		cmd0, err := kvs.CreateTxDelCmd(ctx, kv0.Key)
		if err != nil {
			t.Fatalf("failed to create tx del cmd: %v", err)
		}

		cmd1, err := kvs.CreateTxSetCmd(ctx, kv1.Key, kv1, time.Hour)
		if err != nil {
			t.Fatalf("failed to create tx set cmd: %v", err)
		}

		cmd2, err := kvs.CreateTxSetCmd(ctx, kv2.Key, kv2, time.Hour)
		if err != nil {
			t.Fatalf("failed to create tx set cmd: %v", err)
		}

		if err := kvs.Tx(ctx, []cache.TxSetCmd{cmd1, cmd2}, []cache.TxDelCmd{cmd0}); err != nil {
			t.Fatalf("failed to transaction: %v", err)
		}

		if _, err := kvs.Get(ctx, kv0.Key); err == nil {
			t.Error("failed to delete")
		}

		got1, err := kvs.Get(ctx, kv1.Key)
		if err != nil {
			t.Fatalf("failed to get: %v", err)
		}

		if !reflect.DeepEqual(got1, kv1) {
			t.Errorf("got = %v, want %v", got1, kv1)
		}

		got2, err := kvs.Get(ctx, kv2.Key)
		if err != nil {
			t.Fatalf("failed to get: %v", err)
		}

		if !reflect.DeepEqual(got2, kv2) {
			t.Errorf("got = %v, want %v", got2, kv2)
		}
	})
}
