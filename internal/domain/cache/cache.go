package cache

import (
	"context"
	"time"
)

type Prefix bool

const (
	WithPrefix    Prefix = true
	WithoutPrefix Prefix = false
)

type Cache[T any] interface {
	Get(context.Context, string) (T, error)
	Set(context.Context, string, T, time.Duration) error
	Del(context.Context, string) error
	GetDel(context.Context, string) (T, error)
	CreateTxSetCmd(context.Context, string, T, time.Duration) (TxSetCmd, error)
	CreateTxDelCmd(context.Context, string) (TxDelCmd, error)
	Tx(context.Context, []TxSetCmd, []TxDelCmd) error
	Keys(context.Context, string, Prefix) ([]string, error)
}

type TxSetCmd struct {
	Key   string
	Value string
	TTL   time.Duration
}

type TxDelCmd struct {
	Key string
}
