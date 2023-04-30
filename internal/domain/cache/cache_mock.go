package cache

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"
)

var _ Cache[any] = (*CacheMock[any])(nil)

type CacheMock[V any] struct {
	T                    *testing.T
	Value                V
	GetAssert            func(t *testing.T, key string)
	GetErr               error
	SetAssert            func(t *testing.T, key string, value V, ttl time.Duration)
	SetErr               error
	DelAssert            func(t *testing.T, key string)
	DelErr               error
	CreateTxSetCmdAssert func(t *testing.T, key string, value V, ttl time.Duration)
	CreateTxSetCmdErr    error
	CreateTxDelCmdAssert func(t *testing.T, key string)
	CreateTxDelCmdErr    error
	TxAssert             func(t *testing.T, setCmds []TxSetCmd, delCmds []TxDelCmd)
	TxErr                error
	KeysValue            []string
	KeysAssert           func(t *testing.T, pattern string)
	KeysErr              error
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

func (mock *CacheMock[V]) CreateTxSetCmd(ctx context.Context, key string, value V, ttl time.Duration) (TxSetCmd, error) {
	mock.T.Helper()

	mock.CreateTxSetCmdAssert(mock.T, key, value, ttl)

	if mock.CreateTxSetCmdErr != nil {
		return TxSetCmd{}, mock.CreateTxSetCmdErr
	}

	val, err := json.Marshal(value)
	if err != nil {
		mock.T.Fatalf("failed to marshal json: %v", err)
	}

	return TxSetCmd{
		Key:   key,
		Value: base64.StdEncoding.EncodeToString(val),
		TTL:   ttl,
	}, nil
}

func (mock *CacheMock[V]) CreateTxDelCmd(ctx context.Context, key string) (TxDelCmd, error) {
	mock.T.Helper()

	mock.CreateTxDelCmdAssert(mock.T, key)

	return TxDelCmd{
		Key: key,
	}, mock.CreateTxDelCmdErr
}

func (mock *CacheMock[V]) Tx(ctx context.Context, setCmds []TxSetCmd, delCmds []TxDelCmd) error {
	mock.T.Helper()

	mock.TxAssert(mock.T, setCmds, delCmds)

	return mock.TxErr
}

func (mock *CacheMock[V]) Keys(ctx context.Context, pattern string) ([]string, error) {
	mock.T.Helper()

	mock.KeysAssert(mock.T, pattern)

	return mock.KeysValue, mock.KeysErr
}
