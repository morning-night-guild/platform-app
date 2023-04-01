package port

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
)

var _ CoreUserCreate = (*CoreUserCreateMock)(nil)

type CoreUserCreateMock struct {
	T    *testing.T
	User model.User
	Err  error
}

func (mock *CoreUserCreateMock) Execute(
	ctx context.Context,
	input CoreUserCreateInput,
) (CoreUserCreateOutput, error) {
	mock.T.Helper()

	return CoreUserCreateOutput{
		User: mock.User,
	}, mock.Err
}

var _ CoreUserUpdate = (*CoreUserUpdateMock)(nil)

type CoreUserUpdateMock struct {
	T    *testing.T
	User model.User
	Err  error
}

func (mock *CoreUserUpdateMock) Execute(
	ctx context.Context,
	input CoreUserUpdateInput,
) (CoreUserUpdateOutput, error) {
	mock.T.Helper()

	mock.User.UserID = input.UserID

	return CoreUserUpdateOutput{
		User: mock.User,
	}, mock.Err
}
