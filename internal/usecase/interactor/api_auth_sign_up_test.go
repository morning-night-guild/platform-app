package interactor_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/mock"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestAPIAuthSignUpExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRPC rpc.User
		authRPC rpc.Auth
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
				userRPC: &mock.UserRPC{
					T: t,
					User: model.User{
						UserID: user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					},
				},
				authRPC: &mock.AuthRPC{
					T: t,
					SignUpAssert: func(
						t *testing.T,
						userID user.UserID,
						email auth.EMail,
						password auth.Password,
					) {
						t.Helper()

						if !reflect.DeepEqual(userID, user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab"))) {
							t.Errorf("userID = %v, want %v", userID, user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")))
						}
						if !reflect.DeepEqual(email, auth.EMail("test@example.com")) {
							t.Errorf("email = %v, want %v", email, auth.EMail("test@example.com"))
						}
						if !reflect.DeepEqual(password, auth.Password("password")) {
							t.Errorf("password = %v, want %v", password, auth.Password("password"))
						}
					},
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
				userRPC: &mock.UserRPC{
					T:         t,
					CreateErr: fmt.Errorf("test"),
				},
				authRPC: &mock.AuthRPC{},
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
				userRPC: &mock.UserRPC{
					T: t,
					User: model.User{
						UserID: user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					},
				},
				authRPC: &mock.AuthRPC{
					T: t,
					SignUpAssert: func(
						t *testing.T,
						userID user.UserID,
						email auth.EMail,
						password auth.Password,
					) {
						t.Helper()
					},
					SignUpErr: fmt.Errorf("test"),
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
				tt.fields.userRPC,
				tt.fields.authRPC,
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
