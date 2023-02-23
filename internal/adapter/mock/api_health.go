package mock

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

var _ port.APIHealthCheck = (*APIHealthCheck)(nil)

type APIHealthCheck struct {
	T   *testing.T
	Err error
}

func (ahc *APIHealthCheck) Execute(
	ctx context.Context,
	input port.APIHealthCheckInput,
) (port.APIHealthCheckOutput, error) {
	ahc.T.Helper()

	return port.APIHealthCheckOutput{}, ahc.Err
}
