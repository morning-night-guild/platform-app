package controller

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	userv1 "github.com/morning-night-guild/platform-app/pkg/connect/user/v1"
	"github.com/morning-night-guild/platform-app/pkg/connect/user/v1/userv1connect"
)

var _ userv1connect.UserServiceHandler = (*User)(nil)

// User.
type User struct {
	ctrl    *Controller
	usecase usecase.CoreUser
}

// NewUser ユーザーコントローラを新規作成する関数.
func NewUser(
	ctrl *Controller,
	usecase usecase.CoreUser,
) *User {
	return &User{
		ctrl:    ctrl,
		usecase: usecase,
	}
}

func (ctrl *User) Create(
	ctx context.Context,
	_ *connect.Request[userv1.CreateRequest],
) (*connect.Response[userv1.CreateResponse], error) {
	input := usecase.CoreUserCreateInput{}

	output, err := ctrl.usecase.Create(ctx, input)
	if err != nil {
		return nil, ctrl.ctrl.HandleConnectError(ctx, err)
	}

	res := &userv1.CreateResponse{
		User: &userv1.User{
			UserId: output.User.UserID.String(),
		},
	}

	return connect.NewResponse(res), nil
}

func (ctrl *User) Update(
	ctx context.Context,
	req *connect.Request[userv1.UpdateRequest],
) (*connect.Response[userv1.UpdateResponse], error) {
	uid, err := user.NewID(req.Msg.UserId)
	if err != nil {
		return nil, ctrl.ctrl.HandleConnectError(ctx, err)
	}

	input := usecase.CoreUserUpdateInput{
		UserID: uid,
	}

	output, err := ctrl.usecase.Update(ctx, input)
	if err != nil {
		return nil, ctrl.ctrl.HandleConnectError(ctx, err)
	}

	res := &userv1.UpdateResponse{
		User: &userv1.User{
			UserId: output.User.UserID.String(),
		},
	}

	return connect.NewResponse(res), nil
}
