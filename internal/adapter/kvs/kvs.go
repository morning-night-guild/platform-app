package kvs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/redis/go-redis/v9"
)

const prefix = "%s:%s"

type Factory[T any] interface {
	KVS(string, *redis.Client) (*KVS[T], error)
}

var _ cache.Cache[any] = (*KVS[any])(nil)

type KVS[T any] struct {
	Prefix string
	Client *redis.Client
}

func (kvs *KVS[T]) Get(ctx context.Context, key string) (T, error) {
	var value T

	key = fmt.Sprintf(prefix, kvs.Prefix, key)

	str, err := kvs.Client.Get(ctx, key).Result()
	if err != nil {
		return value, errors.NewNotFoundError("failed to get cache", err)
	}

	if err := json.Unmarshal([]byte(str), &value); err != nil {
		return value, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return value, nil
}

func (kvs *KVS[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	key = fmt.Sprintf(prefix, kvs.Prefix, key)

	if err := kvs.Client.Set(ctx, key, val, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (kvs *KVS[T]) Del(ctx context.Context, key string) error {
	key = fmt.Sprintf(prefix, kvs.Prefix, key)

	if _, err := kvs.Client.Del(ctx, key).Result(); err != nil {
		return fmt.Errorf("failed to del cache: %w", err)
	}

	return nil
}
