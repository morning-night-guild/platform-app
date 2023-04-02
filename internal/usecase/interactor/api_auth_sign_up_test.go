package interactor_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestAPIAuthSignUpExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRPC func(t *testing.T) rpc.User
		authRPC func(t *testing.T) rpc.Auth
	}

	type args struct {
		ctx   context.Context
		input port.APIAuthSignUpInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    port.APIAuthSignUpOutput
		wantErr bool
	}{
		{
			name: "サインアップできる",
			args: args{
				ctx: context.Background(),
				input: port.APIAuthSignUpInput{
					EMail:    auth.EMail("test@example.com"),
					Password: auth.Password("password"),
				},
			},
			fields: fields{
				userRPC: func(t *testing.T) rpc.User {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockUser(ctrl)
					mock.EXPECT().Create(gomock.Any()).Return(model.User{
						UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					}, nil)
					return mock
				},
				authRPC: func(t *testing.T) rpc.Auth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockAuth(ctrl)
					mock.EXPECT().SignUp(
						gomock.Any(),
						user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.EMail("test@example.com"),
						auth.Password("password"),
					).Return(nil)
					return mock
				},
			},
			want:    port.APIAuthSignUpOutput{},
			wantErr: false,
		},
		{
			name: "UserRPC.CreateUser()でエラーが発生してサインアップできない",
			args: args{
				ctx: context.Background(),
				input: port.APIAuthSignUpInput{
					EMail:    auth.EMail("test@example.com"),
					Password: auth.Password("password"),
				},
			},
			fields: fields{
				userRPC: func(t *testing.T) rpc.User {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockUser(ctrl)
					mock.EXPECT().Create(gomock.Any()).Return(model.User{}, fmt.Errorf("test"))
					return mock
				},
				authRPC: func(t *testing.T) rpc.Auth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockAuth(ctrl)
					return mock
				},
			},
			want:    port.APIAuthSignUpOutput{},
			wantErr: true,
		},
		{
			name: "AuthRPC.SignUp()でエラーが発生してサインアップできない",
			args: args{
				ctx: context.Background(),
				input: port.APIAuthSignUpInput{
					EMail:    auth.EMail("test@example.com"),
					Password: auth.Password("password"),
				},
			},
			fields: fields{
				userRPC: func(t *testing.T) rpc.User {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockUser(ctrl)
					mock.EXPECT().Create(gomock.Any()).Return(model.User{
						UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					}, nil)
					return mock
				},
				authRPC: func(t *testing.T) rpc.Auth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockAuth(ctrl)
					mock.EXPECT().SignUp(
						gomock.Any(),
						user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.EMail("test@example.com"),
						auth.Password("password"),
					).Return(fmt.Errorf("test"))
					return mock
				},
			},
			want:    port.APIAuthSignUpOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			aas := interactor.NewAPIAuthSignUp(
				tt.fields.userRPC(t),
				tt.fields.authRPC(t),
			)
			got, err := aas.Execute(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuthSignUp.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIAuthSignUp.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
