package redis

import (
	"context"
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
	client *redis.Client,
) (*kvs.KVS[T], error) {
	return &kvs.KVS[T]{
		Client: client,
	}, nil
}

func NewRedis(url string) (*redis.Client, error) {
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

	client := redis.NewClient(opt)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
