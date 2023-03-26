package mock

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
)

var _ rpc.Auth = (*AuthRPC)(nil)

type AuthRPC struct {
	T            *testing.T
	User         model.User
	SignUpAssert func(t *testing.T, userID user.UserID, email auth.EMail, password auth.Password)
	SignUpErr    error
	SignInAssert func(t *testing.T, email auth.EMail, password auth.Password)
	SignInErr    error
}

func (ar *AuthRPC) SignUp(
	ctx context.Context,
	userID user.UserID,
	email auth.EMail,
	password auth.Password,
) error {
	ar.T.Helper()

	ar.SignUpAssert(ar.T, userID, email, password)

	return ar.SignUpErr
}

func (ar *AuthRPC) SignIn(
	ctx context.Context,
	email auth.EMail,
	password auth.Password,
) (model.User, error) {
	ar.T.Helper()

	ar.SignInAssert(ar.T, email, password)

	return ar.User, ar.SignInErr
}
