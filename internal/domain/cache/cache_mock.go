package cache

import (
	"context"
	"testing"
	"time"
)

var _ Cache[any] = (*CacheMock[any])(nil)

type CacheMock[V any] struct {
	T         *testing.T
	Value     V
	GetErr    error
	GetAssert func(t *testing.T, key string)
	SetErr    error
	SetAssert func(t *testing.T, key string, value V, ttl time.Duration)
	DelErr    error
	DelAssert func(t *testing.T, key string)
}

func (mock *CacheMock[V]) Get(ctx context.Context, key string) (V, error) {
	mock.T.Helper()

	mock.GetAssert(mock.T, key)

	return mock.Value, mock.GetErr
}

func (mock *CacheMock[V]) Set(ctx context.Context, key string, value V, ttl time.Duration) error {
	mock.T.Helper()

	mock.SetAssert(mock.T, key, value, ttl)

	return mock.SetErr
}

func (mock *CacheMock[V]) Del(ctx context.Context, key string) error {
	mock.T.Helper()

	mock.DelAssert(mock.T, key)

	return mock.DelErr
}