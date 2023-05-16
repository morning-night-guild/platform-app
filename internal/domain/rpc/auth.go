package rpc

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

//go:generate mockgen -source auth.go -destination auth_mock.go -package rpc

type Auth interface {
	SignUp(context.Context, user.ID, auth.Email, auth.Password) error
	SignIn(context.Context, auth.Email, auth.Password) (model.User, error)
	ChangePassword(context.Context, user.ID, auth.Password) error
	GetEmail(context.Context, user.ID) (auth.Email, error)
}
