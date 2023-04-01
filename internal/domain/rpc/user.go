package rpc

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

//go:generate mockgen -source user.go -destination user_mock.go -package rpc

type User interface {
	Create(context.Context) (model.User, error)
	Update(context.Context, user.ID) (model.User, error)
}
