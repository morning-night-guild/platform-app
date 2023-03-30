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
	T          *testing.T
	User       model.User
	SaveErr    error
	SaveAssert func(t *testing.T, item model.User)
	FindErr    error
	FindAssert func(t *testing.T, id user.ID)
}

func (ur *UserRepository) Save(ctx context.Context, item model.User) error {
	ur.T.Helper()

	ur.SaveAssert(ur.T, item)

	return ur.SaveErr
}

func (ur *UserRepository) Find(
	ctx context.Context,
	id user.ID,
) (model.User, error) {
	ur.T.Helper()

	ur.User.UserID = id

	return ur.User, ur.FindErr
}
