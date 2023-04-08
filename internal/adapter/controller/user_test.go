package controller_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/adapter/controller"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
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
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreUser(ctrl)
					mock.EXPECT().Create(gomock.Any(), usecase.CoreUserCreateInput{}).Return(usecase.CoreUserCreateOutput{
						User: model.User{
							UserID: user.ID(uuid.MustParse(uid)),
						},
					}, nil)
					return mock
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
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreUser(ctrl)
					mock.EXPECT().Create(gomock.Any(), usecase.CoreUserCreateInput{}).Return(usecase.CoreUserCreateOutput{}, fmt.Errorf("error"))
					return mock
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
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreUser(ctrl)
					mock.EXPECT().Update(gomock.Any(), usecase.CoreUserUpdateInput{
						UserID: user.ID(uuid.MustParse(uid)),
					}).Return(usecase.CoreUserUpdateOutput{
						User: model.User{
							UserID: user.ID(uuid.MustParse(uid)),
						},
					}, nil)
					return mock
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
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreUser(ctrl)
					mock.EXPECT().Update(gomock.Any(), usecase.CoreUserUpdateInput{
						UserID: user.ID(uuid.MustParse(uid)),
					}).Return(usecase.CoreUserUpdateOutput{}, fmt.Errorf("error"))
					return mock
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
