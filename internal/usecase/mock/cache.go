package mock

import (
	"context"
	"testing"
	"time"

	"github.com/morning-night-guild/platform-app/internal/domain/cache"
)

var _ cache.Cache[any] = (*Cache[any])(nil)

type Cache[V any] struct {
	T         *testing.T
	Value     V
	GetErr    error
	GetAssert func(t *testing.T, key string)
	SetErr    error
	SetAssert func(t *testing.T, key string, value V, ttl time.Duration)
	DelErr    error
	DelAssert func(t *testing.T, key string)
}

func (cc *Cache[V]) Get(ctx context.Context, key string) (V, error) {
	cc.T.Helper()

	cc.GetAssert(cc.T, key)

	return cc.Value, cc.GetErr
}

func (cc *Cache[V]) Set(ctx context.Context, key string, value V, ttl time.Duration) error {
	cc.T.Helper()

	cc.SetAssert(cc.T, key, value, ttl)

	return cc.SetErr
}

func (cc *Cache[V]) Del(ctx context.Context, key string) error {
	cc.T.Helper()

	cc.DelAssert(cc.T, key)

	return cc.DelErr
}
