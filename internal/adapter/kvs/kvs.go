package kvs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/redis/go-redis/v9"
)

const (
	format    = "%s:%s"
	keyFormat = "%s:%s*"
)

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

	key = fmt.Sprintf(format, kvs.Prefix, key)

	str, err := kvs.Client.Get(ctx, key).Result()
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get cache", log.ErrorField(err))

		return value, errors.NewNotFoundError("failed to get cache", err)
	}

	dec, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return value, fmt.Errorf("failed to decode base64: %w", err)
	}

	if err := json.Unmarshal(dec, &value); err != nil {
		return value, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return value, nil
}

func (kvs *KVS[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	key = fmt.Sprintf(format, kvs.Prefix, key)

	enc := base64.StdEncoding.EncodeToString(val)

	if err := kvs.Client.Set(ctx, key, enc, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (kvs *KVS[T]) Del(ctx context.Context, key string) error {
	key = fmt.Sprintf(format, kvs.Prefix, key)

	if _, err := kvs.Client.Del(ctx, key).Result(); err != nil {
		return fmt.Errorf("failed to del cache: %w", err)
	}

	return nil
}

func (kvs *KVS[T]) GetDel(ctx context.Context, key string) (T, error) {
	var value T

	key = fmt.Sprintf(format, kvs.Prefix, key)

	str, err := kvs.Client.GetDel(ctx, key).Result()
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get cache", log.ErrorField(err))

		return value, errors.NewNotFoundError("failed to get cache", err)
	}

	dec, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return value, fmt.Errorf("failed to decode base64: %w", err)
	}

	if err := json.Unmarshal(dec, &value); err != nil {
		return value, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return value, nil
}

func (kvs *KVS[T]) CreateTxSetCmd(_ context.Context, key string, value T, ttl time.Duration) (cache.TxSetCmd, error) {
	val, err := json.Marshal(value)
	if err != nil {
		return cache.TxSetCmd{}, fmt.Errorf("failed to marshal json: %w", err)
	}

	key = fmt.Sprintf(format, kvs.Prefix, key)

	enc := base64.StdEncoding.EncodeToString(val)

	return cache.TxSetCmd{
		Key:   key,
		Value: enc,
		TTL:   ttl,
	}, nil
}

func (kvs *KVS[T]) CreateTxDelCmd(_ context.Context, key string) (cache.TxDelCmd, error) {
	key = fmt.Sprintf(format, kvs.Prefix, key)

	return cache.TxDelCmd{
		Key: key,
	}, nil
}

func (kvs *KVS[T]) Tx(
	ctx context.Context,
	sets []cache.TxSetCmd,
	dels []cache.TxDelCmd,
) error {
	if _, err := kvs.Client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		if len(sets) != 0 {
			for _, set := range sets {
				if err := pipe.Set(ctx, set.Key, set.Value, set.TTL).Err(); err != nil {
					return fmt.Errorf("failed to set cache: %w", err)
				}
			}
		}

		if len(dels) != 0 {
			for _, del := range dels {
				if err := pipe.Del(ctx, del.Key).Err(); err != nil {
					return fmt.Errorf("failed to del cache: %w", err)
				}
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to transaction: %w", err)
	}

	return nil
}

func (kvs *KVS[T]) Keys(
	ctx context.Context,
	pattern string,
	prefix cache.Prefix,
) ([]string, error) {
	ptn := fmt.Sprintf(keyFormat, kvs.Prefix, pattern)

	keys, err := kvs.Client.Keys(ctx, ptn).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}

	if prefix {
		return keys, nil
	}

	res := make([]string, len(keys))

	for i, key := range keys {
		res[i] = strings.ReplaceAll(key, fmt.Sprintf("%s:", kvs.Prefix), "")
	}

	return res, nil
}
