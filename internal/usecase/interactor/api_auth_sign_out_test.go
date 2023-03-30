package interactor_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/mock"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestAPIAuthSignOutExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		secret       auth.Secret
		authCache    cache.Cache[model.Auth]
		sessionCache cache.Cache[model.Session]
	}

	type args struct {
		ctx   context.Context
		input port.APIAuthSignOutInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    port.APIAuthSignOutOutput
		wantErr bool
	}{
		{
			name: "サインアウトできる",
			fields: fields{
				secret: auth.Secret("secret"),
				authCache: &mock.Cache[model.Auth]{
					T: t,
					DelAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				sessionCache: &mock.Cache[model.Session]{
					T: t,
					DelAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.APIAuthSignOutInput{
					AuthToken: auth.GenerateAuthToken(
						user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret(),
					),
					SessionToken: auth.GenerateSessionToken(
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.Secret("secret"),
					),
				},
			},
			want:    port.APIAuthSignOutOutput{},
			wantErr: false,
		},
		{
			name: "AuthCache.Del()でエラーが発生してもサインアウトできる",
			fields: fields{
				secret: auth.Secret("secret"),
				authCache: &mock.Cache[model.Auth]{
					T:      t,
					DelErr: fmt.Errorf("test"),
					DelAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				sessionCache: &mock.Cache[model.Session]{
					T: t,
					DelAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.APIAuthSignOutInput{
					AuthToken: auth.GenerateAuthToken(
						user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret(),
					),
					SessionToken: auth.GenerateSessionToken(
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.Secret("secret"),
					),
				},
			},
			want:    port.APIAuthSignOutOutput{},
			wantErr: false,
		},
		{
			name: "SessionCache.Del()でエラーが発生してもサインアウトできる",
			fields: fields{
				secret: auth.Secret("secret"),
				authCache: &mock.Cache[model.Auth]{
					T: t,
					DelAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				sessionCache: &mock.Cache[model.Session]{
					T:      t,
					DelErr: fmt.Errorf("test"),
					DelAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.APIAuthSignOutInput{
					AuthToken: auth.GenerateAuthToken(
						user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret(),
					),
					SessionToken: auth.GenerateSessionToken(
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.Secret("secret"),
					),
				},
			},
			want:    port.APIAuthSignOutOutput{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			aas := interactor.NewAPIAuthSignOut(
				tt.fields.secret,
				tt.fields.authCache,
				tt.fields.sessionCache,
			)
			got, err := aas.Execute(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuthSignOut.Execute() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIAuthSignOut.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
