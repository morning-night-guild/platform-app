package cache

import (
	"context"
	"time"
)

type Cache[T any] interface {
	Get(context.Context, string) (T, error)
	Set(context.Context, string, T, time.Duration) error
	Del(context.Context, string) error
	CreateTxSetCmd(context.Context, string, T, time.Duration) (TxSetCmd, error)
	CreateTxDelCmd(context.Context, string) (TxDelCmd, error)
	Tx(context.Context, []TxSetCmd, []TxDelCmd) error
}

type TxSetCmd struct {
	Key   string
	Value string
	TTL   time.Duration
}

type TxDelCmd struct {
	Key string
}
