package port

import (
	"github.com/morning-night-guild/platform-app/internal/domain/model"
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
