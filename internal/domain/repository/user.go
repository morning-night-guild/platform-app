package repository

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

//go:generate mockgen -source user.go -destination user_mock.go -package repository

type User interface {
	Save(context.Context, model.User) error
	Find(context.Context, user.ID) (model.User, error)
}
