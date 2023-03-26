package repository

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

type User interface {
	Save(context.Context, model.User) error
	Find(context.Context, user.UserID) (model.User, error)
}
