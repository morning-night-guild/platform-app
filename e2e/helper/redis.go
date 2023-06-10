package helper

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	T      *testing.T
	client *redis.Client
}

func NewRedis(
	t *testing.T,
	url string,
) *Redis {
	t.Helper()

	client := redis.NewClient(&redis.Options{
		Addr: url,
	})

	return &Redis{
		T:      t,
		client: client,
	}
}

func (rds *Redis) DeleteInvitationCode(
	code string,
) {
	rds.T.Helper()

	key := fmt.Sprintf("invitation:%s", code)

	if err := rds.client.Del(context.Background(), key).Err(); err != nil {
		rds.T.Errorf("failed to delete invitation code: %v", err)
	}
}

func (rds *Redis) Close() {
	rds.T.Helper()

	if err := rds.client.Close(); err != nil {
		rds.T.Errorf("failed to close redis: %v", err)
	}
}
