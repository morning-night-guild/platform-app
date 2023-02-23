package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

var _ port.APIHealthCheck = (*APIHealthCheck)(nil)

type APIHealthCheck struct {
	healthRepository repository.Health
}

func NewAPIHealthCheck(
	healthRepository repository.Health,
) *APIHealthCheck {
	return &APIHealthCheck{
		healthRepository: healthRepository,
	}
}

func (l *APIHealthCheck) Execute(ctx context.Context, input port.APIHealthCheckInput) (port.APIHealthCheckOutput, error) {
	return port.APIHealthCheckOutput{}, nil
}
