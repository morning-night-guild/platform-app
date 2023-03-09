package external

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	healthv1 "github.com/morning-night-guild/platform-app/pkg/connect/health/v1"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ rpc.Health = (*Health)(nil)

type Health struct {
	connect *Connect
}

func NewHealth(connect *Connect) *Health {
	return &Health{
		connect: connect,
	}
}

func (ch *Health) Check(ctx context.Context) error {
	req := NewRequestWithTID(ctx, &healthv1.CheckRequest{})

	if _, err := ch.connect.Health.Check(ctx, req); err != nil {
		log.GetLogCtx(ctx).Warn("failed to check health core", log.ErrorField(err))

		return err
	}

	return nil
}
