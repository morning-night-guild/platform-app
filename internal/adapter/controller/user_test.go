package controller_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/adapter/controller"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	userv1 "github.com/morning-night-guild/platform-app/pkg/connect/user/v1"
)

const uid = "01234567-0123-0123-0123-0123456789ab"

func TestUserCreate(t *testing.T) {
	t.Parallel()

	type fields struct {
		usecase func(*testing.T) usecase.CoreUser
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
				usecase: func(t *testing.T) usecase.CoreUser {
					t.Helper()
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				req: connect.NewRequest(&userv1.CreateRequest{}),
			},
			want: connect.NewResponse(&userv1.CreateResponse{
				User: &userv1.User{
					UserId: uid,
				},
			}),
			wantErr: false,
		},
		{
			name: "ユーザーが作成できない",
			fields: fields{
				usecase: func(t *testing.T) usecase.CoreUser {
					t.Helper()
					return nil
				},
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
				tt.fields.usecase(t),
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
		usecase func(*testing.T) usecase.CoreUser
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
				usecase: func(t *testing.T) usecase.CoreUser {
					t.Helper()
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				req: connect.NewRequest(&userv1.UpdateRequest{
					UserId: uid,
				}),
			},
			want: connect.NewResponse(&userv1.UpdateResponse{
				User: &userv1.User{
					UserId: uid,
				},
			}),
			wantErr: false,
		},
		{
			name: "ユーザーが更新できない",
			fields: fields{
				usecase: func(t *testing.T) usecase.CoreUser {
					t.Helper()
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				req: connect.NewRequest(&userv1.UpdateRequest{
					UserId: uid,
				}),
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
				tt.fields.usecase(t),
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
