package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ port.APIAuthSignUp = (*APIAuthSignUp)(nil)

type APIAuthSignUp struct {
	userRPC rpc.User
	authRPC rpc.Auth
}

func NewAPIAuthSignUp(
	userRPC rpc.User,
	authRPC rpc.Auth,
) *APIAuthSignUp {
	return &APIAuthSignUp{
		userRPC: userRPC,
		authRPC: authRPC,
	}
}

func (aas *APIAuthSignUp) Execute(
	ctx context.Context,
	input port.APIAuthSignUpInput,
) (port.APIAuthSignUpOutput, error) {
	user, err := aas.userRPC.Create(ctx)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to create user", log.ErrorField(err))

		return port.APIAuthSignUpOutput{}, err
	}

	if err := aas.authRPC.SignUp(ctx, user.UserID, input.EMail, input.Password); err != nil {
		log.GetLogCtx(ctx).Warn("failed to sign up", log.ErrorField(err))

		return port.APIAuthSignUpOutput{}, err
	}

	return port.APIAuthSignUpOutput{}, nil
}
