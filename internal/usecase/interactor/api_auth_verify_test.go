package interactor_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/mock"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestAPIAuthVerifyExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		secret    auth.Secret
		authCache cache.Cache[model.Auth]
	}

	type args struct {
		ctx   context.Context
		input port.APIAuthVerifyInput
	}

	now := time.Now()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    port.APIAuthVerifyOutput
		wantErr bool
	}{
		{
			name: "検証できる",
			fields: fields{
				secret: auth.Secret("secret"),
				authCache: &mock.Cache[model.Auth]{
					T: t,
					Value: model.Auth{
						AuthID:    user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.APIAuthVerifyInput{
					AuthToken: auth.GenerateAuthToken(
						user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret(),
					),
					SessionToken: auth.GenerateSessionToken(
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.Secret("secret"),
					),
				},
			},
			want:    port.APIAuthVerifyOutput{},
			wantErr: false,
		},
		{
			name: "AuthがCacheに存在せず検証に失敗する",
			fields: fields{
				secret: auth.Secret("secret"),
				authCache: &mock.Cache[model.Auth]{
					T:      t,
					Value:  model.Auth{},
					GetErr: fmt.Errorf("test"),
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.APIAuthVerifyInput{
					AuthToken: auth.GenerateAuthToken(
						user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret(),
					),
					SessionToken: auth.GenerateSessionToken(
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.Secret("secret"),
					),
				},
			},
			want:    port.APIAuthVerifyOutput{},
			wantErr: true,
		},
		{
			name: "Authの有効期限が切れて検証に失敗する",
			fields: fields{
				secret: auth.Secret("secret"),
				authCache: &mock.Cache[model.Auth]{
					T: t,
					Value: model.Auth{
						AuthID:    user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(-time.Hour),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.APIAuthVerifyInput{
					AuthToken: auth.GenerateAuthToken(
						user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret(),
					),
					SessionToken: auth.GenerateSessionToken(
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.Secret("secret"),
					),
				},
			},
			want:    port.APIAuthVerifyOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			aav := interactor.NewAPIAuthVerify(
				tt.fields.secret,
				tt.fields.authCache,
			)
			got, err := aav.Execute(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuthVerify.Execute() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIAuthVerify.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
