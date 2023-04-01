package port

import (
	"context"
	"testing"
)

var _ APIHealthCheck = (*APIHealthCheckMock)(nil)

type APIHealthCheckMock struct {
	T   *testing.T
	Err error
}

func (mock *APIHealthCheckMock) Execute(
	ctx context.Context,
	input APIHealthCheckInput,
) (APIHealthCheckOutput, error) {
	mock.T.Helper()

	return APIHealthCheckOutput{}, mock.Err
}
