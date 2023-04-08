package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

var _ port.APIHealthCheck = (*APIHealthCheck)(nil)

type APIHealthCheck struct {
	healthRPC rpc.Health
}

func NewAPIHealthCheck(
	healthRPC rpc.Health,
) *APIHealthCheck {
	return &APIHealthCheck{
		healthRPC: healthRPC,
	}
}

func (ahc *APIHealthCheck) Execute(
	ctx context.Context,
	_ port.APIHealthCheckInput,
) (port.APIHealthCheckOutput, error) {
	if err := ahc.healthRPC.Check(ctx); err != nil {
		return port.APIHealthCheckOutput{}, err
	}

	return port.APIHealthCheckOutput{}, nil
}
