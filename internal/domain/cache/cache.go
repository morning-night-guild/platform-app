package cache

import (
	"context"
	"time"
)

type Cache[T any] interface {
	Get(context.Context, string) (T, error)
	Set(context.Context, string, T, time.Duration) error
	Del(context.Context, string) error
}
