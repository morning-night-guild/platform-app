package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ port.CoreUserUpdate = (*CoreUserUpdate)(nil)

type CoreUserUpdate struct {
	userRepository repository.User
}

func NewCoreUserUpdate(
	userRepository repository.User,
) *CoreUserUpdate {
	return &CoreUserUpdate{
		userRepository: userRepository,
	}
}

func (cuu *CoreUserUpdate) Execute(
	ctx context.Context,
	input port.CoreUserUpdateInput,
) (port.CoreUserUpdateOutput, error) {
	user, err := cuu.userRepository.Find(ctx, input.UserID)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to find user", log.ErrorField(err))

		return port.CoreUserUpdateOutput{}, err
	}

	if err := cuu.userRepository.Save(ctx, user); err != nil {
		log.GetLogCtx(ctx).Warn("failed to save user", log.ErrorField(err))

		return port.CoreUserUpdateOutput{}, err
	}

	return port.CoreUserUpdateOutput{
		User: user,
	}, nil
}
