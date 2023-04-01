package port

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
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

var _ APIAuthSignIn = (*APIAuthSignInMock)(nil)

type APIAuthSignInMock struct {
	T            *testing.T
	AuthToken    auth.AuthToken
	SessionToken auth.SessionToken
	Err          error
}

func (mock *APIAuthSignInMock) Execute(
	ctx context.Context,
	input APIAuthSignInInput,
) (APIAuthSignInOutput, error) {
	mock.T.Helper()

	return APIAuthSignInOutput{
		AuthToken:    mock.AuthToken,
		SessionToken: mock.SessionToken,
	}, mock.Err
}

var _ APIAuthSignOut = (*APIAuthSignOutMock)(nil)

type APIAuthSignOutMock struct {
	T   *testing.T
	Err error
}

func (mock *APIAuthSignOutMock) Execute(
	ctx context.Context,
	input APIAuthSignOutInput,
) (APIAuthSignOutOutput, error) {
	mock.T.Helper()

	return APIAuthSignOutOutput{}, mock.Err
}

var _ APIAuthVerify = (*APIAuthVerifyMock)(nil)

type APIAuthVerifyMock struct {
	T   *testing.T
	Err error
}

func (mock *APIAuthVerifyMock) Execute(
	ctx context.Context,
	input APIAuthVerifyInput,
) (APIAuthVerifyOutput, error) {
	mock.T.Helper()

	return APIAuthVerifyOutput{}, mock.Err
}

var _ APIAuthRefresh = (*APIAuthRefreshMock)(nil)

type APIAuthRefreshMock struct {
	T         *testing.T
	AuthToken auth.AuthToken
	Err       error
}

func (mock *APIAuthRefreshMock) Execute(
	ctx context.Context,
	input APIAuthRefreshInput,
) (APIAuthRefreshOutput, error) {
	mock.T.Helper()

	return APIAuthRefreshOutput{
		AuthToken: mock.AuthToken,
	}, mock.Err
}

var _ APIAuthGenerateCode = (*APIAuthGenerateCodeMock)(nil)

type APIAuthGenerateCodeMock struct {
	T    *testing.T
	Code model.Code
	Err  error
}

func (mock *APIAuthGenerateCodeMock) Execute(
	ctx context.Context,
	input APIAuthGenerateCodeInput,
) (APIAuthGenerateCodeOutput, error) {
	mock.T.Helper()

	return APIAuthGenerateCodeOutput{
		Code: mock.Code,
	}, mock.Err
}
