package mock

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
)

var _ rpc.User = (*UserRPC)(nil)

type UserRPC struct {
	T         *testing.T
	User      model.User
	CreateErr error
	UpdateErr error
}

func (ur *UserRPC) Create(ctx context.Context) (model.User, error) {
	ur.T.Helper()

	return ur.User, ur.CreateErr
}

func (ur *UserRPC) Update(ctx context.Context, id user.ID) (model.User, error) {
	ur.T.Helper()

	return ur.User, ur.UpdateErr
}
