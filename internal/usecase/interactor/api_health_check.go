package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

type APIHealthCheck struct {
	healthRepository repository.APIHealth
}

func NewAPIHealthCheck(
	healthRepository repository.APIHealth,
) port.APIHealthCheck {
	return &APIHealthCheck{
		healthRepository: healthRepository,
	}
}

func (ahc *APIHealthCheck) Execute(
	ctx context.Context,
	input port.APIHealthCheckInput,
) (port.APIHealthCheckOutput, error) {
	if err := ahc.healthRepository.Check(ctx); err != nil {
		return port.APIHealthCheckOutput{}, err
	}

	return port.APIHealthCheckOutput{}, nil
}
