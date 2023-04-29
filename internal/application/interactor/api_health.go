package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
)

var _ usecase.APIHealth = (*APIHealth)(nil)

type APIHealth struct {
	healthRPC rpc.Health
}

func NewAPIHealth(
	healthRPC rpc.Health,
) *APIHealth {
	return &APIHealth{
		healthRPC: healthRPC,
	}
}

func (itr *APIHealth) Check(
	ctx context.Context,
	_ usecase.APIHealthCheckInput,
) (usecase.APIHealthCheckOutput, error) {
	if err := itr.healthRPC.Check(ctx); err != nil {
		return usecase.APIHealthCheckOutput{}, err
	}

	return usecase.APIHealthCheckOutput{}, nil
}
