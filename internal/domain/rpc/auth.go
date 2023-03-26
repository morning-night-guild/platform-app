package rpc

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

type Auth interface {
	SignUp(context.Context, user.UserID, auth.EMail, auth.Password) error
	SignIn(context.Context, auth.EMail, auth.Password) (model.User, error)
}
