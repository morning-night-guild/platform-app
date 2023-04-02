package external

import (
	"context"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	userv1 "github.com/morning-night-guild/platform-app/pkg/connect/user/v1"
	"github.com/morning-night-guild/platform-app/pkg/connect/user/v1/userv1connect"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

type UserFactory interface {
	User(string) (*User, error)
}

var _ rpc.User = (*User)(nil)

type User struct {
	connect userv1connect.UserServiceClient
}

func NewUser(
	connect userv1connect.UserServiceClient,
) *User {
	return &User{
		connect: connect,
	}
}

func (usr *User) Create(ctx context.Context) (model.User, error) {
	req := NewRequestWithTID(ctx, &userv1.CreateRequest{})

	res, err := usr.connect.Create(ctx, req)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to create user core", log.ErrorField(err))

		return model.User{}, err
	}

	return model.User{
		UserID: user.ID(uuid.MustParse(res.Msg.User.UserId)),
	}, nil
}

func (usr *User) Update(ctx context.Context, uid user.ID) (model.User, error) {
	req := NewRequestWithTID(ctx, &userv1.UpdateRequest{
		UserId: uid.String(),
	})

	res, err := usr.connect.Update(ctx, req)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to update user core", log.ErrorField(err))

		return model.User{}, err
	}

	return model.User{
		UserID: user.ID(uuid.MustParse(res.Msg.User.UserId)),
	}, nil
}
