package controller

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	userv1 "github.com/morning-night-guild/platform-app/pkg/connect/user/v1"
	"github.com/morning-night-guild/platform-app/pkg/connect/user/v1/userv1connect"
)

var _ userv1connect.UserServiceHandler = (*User)(nil)

// User.
type User struct {
	ctl    *Controller
	create port.CoreUserCreate
	update port.CoreUserUpdate
}

// NewUser ユーザーコントローラを新規作成する関数.
func NewUser(
	ctl *Controller,
	create port.CoreUserCreate,
	update port.CoreUserUpdate,
) *User {
	return &User{
		ctl:    ctl,
		create: create,
		update: update,
	}
}

func (usr *User) Create(
	ctx context.Context,
	req *connect.Request[userv1.CreateRequest],
) (*connect.Response[userv1.CreateResponse], error) {
	input := port.CoreUserCreateInput{}

	output, err := usr.create.Execute(ctx, input)
	if err != nil {
		return nil, usr.ctl.HandleConnectError(ctx, err)
	}

	res := &userv1.CreateResponse{
		User: &userv1.User{
			UserId: output.User.UserID.String(),
		},
	}

	return connect.NewResponse(res), nil
}

func (usr *User) Update(
	ctx context.Context,
	req *connect.Request[userv1.UpdateRequest],
) (*connect.Response[userv1.UpdateResponse], error) {
	uid, err := user.NewID(req.Msg.UserId)
	if err != nil {
		return nil, usr.ctl.HandleConnectError(ctx, err)
	}

	input := port.CoreUserUpdateInput{
		UserID: uid,
	}

	output, err := usr.update.Execute(ctx, input)
	if err != nil {
		return nil, usr.ctl.HandleConnectError(ctx, err)
	}

	res := &userv1.UpdateResponse{
		User: &userv1.User{
			UserId: output.User.UserID.String(),
		},
	}

	return connect.NewResponse(res), nil
}
