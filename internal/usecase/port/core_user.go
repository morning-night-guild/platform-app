package port

import (
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/usecase"
)

type CoreUserCreateInput struct {
	usecase.Input
}

type CoreUserCreateOutput struct {
	usecase.Output
	User model.User
}

type CoreUserCreate interface {
	usecase.Usecase[CoreUserCreateInput, CoreUserCreateOutput]
}

type CoreUserUpdateInput struct {
	usecase.Input
	UserID user.ID
}

type CoreUserUpdateOutput struct {
	usecase.Output
	User model.User
}

type CoreUserUpdate interface {
	usecase.Usecase[CoreUserUpdateInput, CoreUserUpdateOutput]
}
