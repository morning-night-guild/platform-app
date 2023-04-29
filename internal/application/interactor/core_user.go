package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ usecase.CoreUser = (*CoreUser)(nil)

type CoreUser struct {
	userRepository repository.User
}

func NewCoreUser(
	userRepository repository.User,
) *CoreUser {
	return &CoreUser{
		userRepository: userRepository,
	}
}

func (itr *CoreUser) Create(
	ctx context.Context,
	_ usecase.CoreUserCreateInput,
) (usecase.CoreUserCreateOutput, error) {
	user := model.CreateUser()

	if err := itr.userRepository.Save(ctx, user); err != nil {
		log.GetLogCtx(ctx).Warn("failed to save user", log.ErrorField(err))

		return usecase.CoreUserCreateOutput{}, err
	}

	return usecase.CoreUserCreateOutput{
		User: user,
	}, nil
}

func (itr *CoreUser) Update(
	ctx context.Context,
	input usecase.CoreUserUpdateInput,
) (usecase.CoreUserUpdateOutput, error) {
	user, err := itr.userRepository.Find(ctx, input.UserID)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to find user", log.ErrorField(err))

		return usecase.CoreUserUpdateOutput{}, err
	}

	if err := itr.userRepository.Save(ctx, user); err != nil {
		log.GetLogCtx(ctx).Warn("failed to save user", log.ErrorField(err))

		return usecase.CoreUserUpdateOutput{}, err
	}

	return usecase.CoreUserUpdateOutput{
		User: user,
	}, nil
}
