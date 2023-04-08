package usecase

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

//go:generate mockgen -source core_user.go -destination core_user_mock.go -package usecase

// CoreUser.
type CoreUser interface {
	Create(context.Context, CoreUserCreateInput) (CoreUserCreateOutput, error)
	Update(context.Context, CoreUserUpdateInput) (CoreUserUpdateOutput, error)
}

// CoreUserCreateInput.
type CoreUserCreateInput struct{}

// CoreUserCreateOutput.
type CoreUserCreateOutput struct {
	User model.User
}

// CoreUserUpdateInput.
type CoreUserUpdateInput struct {
	UserID user.ID
}

// CoreUserUpdateOutput.
type CoreUserUpdateOutput struct {
	User model.User
}
