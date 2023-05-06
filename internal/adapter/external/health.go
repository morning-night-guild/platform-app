package external

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	healthv1 "github.com/morning-night-guild/platform-app/pkg/connect/health/v1"
	"github.com/morning-night-guild/platform-app/pkg/connect/health/v1/healthv1connect"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

type HealthFactory interface {
	Health(string) (*Health, error)
}

var _ rpc.Health = (*Health)(nil)

type Health struct {
	connect  healthv1connect.HealthServiceClient
	external *External
}

func NewHealth(
	connect healthv1connect.HealthServiceClient,
) *Health {
	return &Health{
		connect:  connect,
		external: New(),
	}
}

func (ext *Health) Check(ctx context.Context) error {
	req := NewRequest(ctx, &healthv1.CheckRequest{})

	if _, err := ext.connect.Check(ctx, req); err != nil {
		log.GetLogCtx(ctx).Warn("failed to check health core", log.ErrorField(err))

		return ext.external.HandleError(ctx, err)
	}

	return nil
}
