package controller_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/adapter/controller"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	userv1 "github.com/morning-night-guild/platform-app/pkg/connect/user/v1"
)

func TestUserCreate(t *testing.T) {
	t.Parallel()

	type fields struct {
		create port.CoreUserCreate
		update port.CoreUserUpdate
	}

	type args struct {
		ctx context.Context
		req *connect.Request[userv1.CreateRequest]
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *connect.Response[userv1.CreateResponse]
		wantErr bool
	}{
		{
			name: "ユーザーが作成できる",
			fields: fields{
				create: &port.CoreUserCreateMock{
					T: t,
					User: model.User{
						UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					},
				},
				update: &port.CoreUserUpdateMock{},
			},
			args: args{
				ctx: context.Background(),
				req: connect.NewRequest(&userv1.CreateRequest{}),
			},
			want: connect.NewResponse(&userv1.CreateResponse{
				User: &userv1.User{
					UserId: "01234567-0123-0123-0123-0123456789ab",
				},
			}),
			wantErr: false,
		},
		{
			name: "ユーザーが作成できない",
			fields: fields{
				create: &port.CoreUserCreateMock{
					T:    t,
					User: model.User{},
					Err:  fmt.Errorf("test"),
				},
				update: &port.CoreUserUpdateMock{},
			},
			args: args{
				ctx: context.Background(),
				req: connect.NewRequest(&userv1.CreateRequest{}),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			usr := controller.NewUser(
				controller.New(),
				tt.fields.create,
				tt.fields.update,
			)
			got, err := usr.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserUpdate(t *testing.T) {
	t.Parallel()

	type fields struct {
		create port.CoreUserCreate
		update port.CoreUserUpdate
	}

	type args struct {
		ctx context.Context
		req *connect.Request[userv1.UpdateRequest]
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *connect.Response[userv1.UpdateResponse]
		wantErr bool
	}{
		{
			name: "ユーザーが更新できる",
			fields: fields{
				create: &port.CoreUserCreateMock{},
				update: &port.CoreUserUpdateMock{
					T: t,
					User: model.User{
						UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: connect.NewRequest(&userv1.UpdateRequest{
					UserId: "01234567-0123-0123-0123-0123456789ab",
				}),
			},
			want: connect.NewResponse(&userv1.UpdateResponse{
				User: &userv1.User{
					UserId: "01234567-0123-0123-0123-0123456789ab",
				},
			}),
			wantErr: false,
		},
		{
			name: "ユーザーが更新できない",
			fields: fields{
				create: &port.CoreUserCreateMock{},
				update: &port.CoreUserUpdateMock{
					T:   t,
					Err: fmt.Errorf("test"),
				},
			},
			args: args{
				ctx: context.Background(),
				req: connect.NewRequest(&userv1.UpdateRequest{
					UserId: "01234567-0123-0123-0123-0123456789ab",
				}),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usr := controller.NewUser(
				controller.New(),
				tt.fields.create,
				tt.fields.update,
			)
			got, err := usr.Update(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
