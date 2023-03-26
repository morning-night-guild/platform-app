package interactor_test

import (
	"context"
	"crypto/rsa"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/mock"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestAPIAuthSignInExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		secret       auth.Secret
		authRPC      rpc.Auth
		authCache    cache.Cache[model.Auth]
		sessionCache cache.Cache[model.Session]
	}

	type args struct {
		ctx   context.Context
		input port.APIAuthSignInInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    port.APIAuthSignInOutput
		wantErr bool
	}{
		{
			name: "サインインできる",
			fields: fields{
				secret: auth.Secret("secret"),
				authRPC: &mock.AuthRPC{
					T: t,
					User: model.User{
						UserID: user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					},
					SignInAssert: func(t *testing.T, email auth.EMail, password auth.Password) {
						t.Helper()
						if !reflect.DeepEqual(email, auth.EMail("test@example.com")) {
							t.Errorf("email = %v, want %v", email, auth.EMail("test@example.com"))
						}
						if !reflect.DeepEqual(password, auth.Password("password")) {
							t.Errorf("password = %v, want %v", password, auth.Password("password"))
						}
					},
				},
				authCache: &mock.Cache[model.Auth]{
					T: t,
					SetAssert: func(t *testing.T, key string, value model.Auth, ttl time.Duration) {
						t.Helper()
						if !reflect.DeepEqual(ttl, model.DefaultAuthExpiresIn) {
							t.Errorf("ttl = %v, want %v", ttl, model.DefaultAuthExpiresIn)
						}
					},
				},
				sessionCache: &mock.Cache[model.Session]{
					T: t,
					SetAssert: func(t *testing.T, key string, value model.Session, ttl time.Duration) {
						t.Helper()
						if !reflect.DeepEqual(ttl, model.DefaultSessionExpiresIn) {
							t.Errorf("ttl = %v, want %v", ttl, model.DefaultSessionExpiresIn)
						}
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.APIAuthSignInInput{
					EMail:     auth.EMail("test@example.com"),
					Password:  auth.Password("password"),
					PublicKey: rsa.PublicKey{},
				},
			},
			want: port.APIAuthSignInOutput{
				AuthToken: auth.GenerateAuthToken(
					user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret(),
				),
				SessionToken: auth.GenerateSessionToken(
					auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					auth.Secret("secret"),
				),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			aas := interactor.NewAPIAuthSignIn(
				tt.fields.secret,
				tt.fields.authRPC,
				tt.fields.authCache,
				tt.fields.sessionCache,
			)
			_, err := aas.Execute(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuthSignIn.Execute() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("APIAuthSignIn.Execute() = %v, want %v", got, tt.want)
			// }
		})
	}
}
