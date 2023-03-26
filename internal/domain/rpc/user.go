package rpc

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

type User interface {
	Create(context.Context) (model.User, error)
	Update(context.Context, user.UserID) (model.User, error)
}
