package mock

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
)

var _ repository.User = (*UserRepository)(nil)

type UserRepository struct {
	T    *testing.T
	User model.User
	Err  error
}

func (ur *UserRepository) Save(ctx context.Context, User model.User) error {
	ur.T.Helper()

	return ur.Err
}

func (ur *UserRepository) Find(
	ctx context.Context,
	id user.UserID,
) (model.User, error) {
	ur.T.Helper()

	return ur.User, ur.Err
}
