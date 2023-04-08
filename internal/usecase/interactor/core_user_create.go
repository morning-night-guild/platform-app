package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ port.CoreUserCreate = (*CoreUserCreate)(nil)

type CoreUserCreate struct {
	userRepository repository.User
}

func NewCoreUserCreate(
	userRepository repository.User,
) *CoreUserCreate {
	return &CoreUserCreate{
		userRepository: userRepository,
	}
}

func (cuc *CoreUserCreate) Execute(
	ctx context.Context,
	_ port.CoreUserCreateInput,
) (port.CoreUserCreateOutput, error) {
	user := model.CreateUser()

	if err := cuc.userRepository.Save(ctx, user); err != nil {
		log.GetLogCtx(ctx).Warn("failed to save user", log.ErrorField(err))

		return port.CoreUserCreateOutput{}, err
	}

	return port.CoreUserCreateOutput{
		User: user,
	}, nil
}
