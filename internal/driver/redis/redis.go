package redis

import (
	"fmt"

	"github.com/morning-night-guild/platform-app/internal/adapter/kvs"
	"github.com/morning-night-guild/platform-app/internal/driver/env"
	"github.com/redis/go-redis/v9"
)

var _ kvs.Factory[any] = (*Redis[any])(nil)

type Redis[T any] struct{}

func New[T any]() *Redis[T] {
	return &Redis[T]{}
}

func (rds *Redis[T]) KVS(
	url string,
) (*kvs.KVS[T], error) {
	var opt *redis.Options

	if env.Get().IsProd() {
		var err error

		opt, err = redis.ParseURL(url)

		if err != nil {
			return nil, fmt.Errorf("failed to parse redis url: %w", err)
		}
	} else {
		opt = &redis.Options{
			Addr:     url,
			Password: "",
			DB:       0,
		}
	}

	return &kvs.KVS[T]{
		Client: redis.NewClient(opt),
	}, nil
}
