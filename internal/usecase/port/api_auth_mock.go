package port

import (
	"context"
	"testing"
)

var _ APIAuthSignUp = (*APIAuthSignUpMock)(nil)

type APIAuthSignUpMock struct {
	T   *testing.T
	Err error
}

func (mock *APIAuthSignUpMock) Execute(
	ctx context.Context,
	input APIAuthSignUpInput,
) (APIAuthSignUpOutput, error) {
	mock.T.Helper()

	return APIAuthSignUpOutput{}, mock.Err
}
