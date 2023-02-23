package gateway

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	healthv1 "github.com/morning-night-guild/platform-app/pkg/connect/health/v1"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ repository.APIHealth = (*APIHealth)(nil)

type APIHealth struct {
	connect *Connect
}

func NewAPIHealth(connect *Connect) *APIHealth {
	return &APIHealth{
		connect: connect,
	}
}

func (ch *APIHealth) Check(ctx context.Context) error {
	req := NewRequestWithTID(ctx, &healthv1.CheckRequest{})

	if _, err := ch.connect.Health.Check(ctx, req); err != nil {
		log.GetLogCtx(ctx).Warn("failed to check health core", log.ErrorField(err))

		return err
	}

	return nil
}
