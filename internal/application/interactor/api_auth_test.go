package interactor_test

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/application/interactor"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
)

func sign(t *testing.T, prv *rsa.PrivateKey, code string) auth.Signature {
	t.Helper()

	h := crypto.Hash.New(crypto.SHA256)

	h.Write([]byte(code))

	hashed := h.Sum(nil)

	signed, err := rsa.SignPSS(rand.Reader, prv, crypto.SHA256, hashed, nil)
	if err != nil {
		t.Fatal(err)
	}

	return auth.Signature(base64.StdEncoding.EncodeToString(signed))
}

func generateKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	prv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	return prv
}

func TestAPIAuthSignUp(t *testing.T) {
	t.Parallel()

	type fields struct {
		secret       auth.Secret
		authRPC      func(t *testing.T) rpc.Auth
		userRPC      func(t *testing.T) rpc.User
		authCache    cache.Cache[model.Auth]
		codeCache    cache.Cache[model.Code]
		sessionCache cache.Cache[model.Session]
	}

	type args struct {
		ctx   context.Context
		input usecase.APIAuthSignUpInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIAuthSignUpOutput
		wantErr bool
	}{
		{
			name: "サインアップできる",
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthSignUpInput{
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
			want:    usecase.APIAuthSignUpOutput{},
			wantErr: false,
		},
		{
			name: "UserRPC.CreateUser()でエラーが発生してサインアップできない",
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthSignUpInput{
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
			want:    usecase.APIAuthSignUpOutput{},
			wantErr: true,
		},
		{
			name: "AuthRPC.SignUp()でエラーが発生してサインアップできない",
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthSignUpInput{
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
			want:    usecase.APIAuthSignUpOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			aa := interactor.NewAPIAuth(
				tt.fields.secret,
				tt.fields.authRPC(t),
				tt.fields.userRPC(t),
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			got, err := aa.SignUp(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuth.SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIAuth.SignUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIAuthSignIn(t *testing.T) {
	t.Parallel()

	type fields struct {
		secret       auth.Secret
		authRPC      func(t *testing.T) rpc.Auth
		userRPC      rpc.User
		authCache    cache.Cache[model.Auth]
		codeCache    cache.Cache[model.Code]
		sessionCache cache.Cache[model.Session]
	}

	type args struct {
		ctx   context.Context
		input usecase.APIAuthSignInInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIAuthSignInOutput
		wantErr bool
	}{
		{
			name: "サインインできる",
			fields: fields{
				secret: auth.Secret("secret"),
				authRPC: func(t *testing.T) rpc.Auth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockAuth(ctrl)
					mock.EXPECT().SignIn(
						gomock.Any(),
						auth.EMail("test@example.com"),
						auth.Password("password"),
					).Return(model.User{
						UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					}, nil)
					return mock
				},
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					SetAssert: func(t *testing.T, key string, value model.Auth, ttl time.Duration) {
						t.Helper()
						if !reflect.DeepEqual(ttl, model.DefaultAuthExpiresIn) {
							t.Errorf("ttl = %v, want %v", ttl, model.DefaultAuthExpiresIn)
						}
					},
				},
				sessionCache: &cache.CacheMock[model.Session]{
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
				input: usecase.APIAuthSignInInput{
					EMail:     auth.EMail("test@example.com"),
					Password:  auth.Password("password"),
					PublicKey: rsa.PublicKey{},
				},
			},
			want: usecase.APIAuthSignInOutput{
				AuthToken: auth.GenerateAuthToken(
					user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
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
			aa := interactor.NewAPIAuth(
				tt.fields.secret,
				tt.fields.authRPC(t),
				tt.fields.userRPC,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			_, err := aa.SignIn(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuth.SignIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("APIAuth.SignIn() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestAPIAuthSignOut(t *testing.T) {
	t.Parallel()

	type fields struct {
		secret       auth.Secret
		authRPC      rpc.Auth
		userRPC      rpc.User
		authCache    cache.Cache[model.Auth]
		codeCache    cache.Cache[model.Code]
		sessionCache cache.Cache[model.Session]
	}

	type args struct {
		ctx   context.Context
		input usecase.APIAuthSignOutInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIAuthSignOutOutput
		wantErr bool
	}{
		{
			name: "サインアウトできる",
			fields: fields{
				secret: auth.Secret("secret"),
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					DelAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					DelAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthSignOutInput{
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
			want:    usecase.APIAuthSignOutOutput{},
			wantErr: false,
		},
		{
			name: "AuthCache.Del()でエラーが発生してもサインアウトできる",
			fields: fields{
				secret: auth.Secret("secret"),
				authCache: &cache.CacheMock[model.Auth]{
					T:      t,
					DelErr: fmt.Errorf("test"),
					DelAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					DelAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthSignOutInput{
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
			want:    usecase.APIAuthSignOutOutput{},
			wantErr: false,
		},
		{
			name: "SessionCache.Del()でエラーが発生してもサインアウトできる",
			fields: fields{
				secret: auth.Secret("secret"),
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					DelAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				sessionCache: &cache.CacheMock[model.Session]{
					T:      t,
					DelErr: fmt.Errorf("test"),
					DelAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthSignOutInput{
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
			want:    usecase.APIAuthSignOutOutput{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			aa := interactor.NewAPIAuth(
				tt.fields.secret,
				tt.fields.authRPC,
				tt.fields.userRPC,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			got, err := aa.SignOut(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuth.SignOut() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIAuth.SignOut() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIAuthVerify(t *testing.T) {
	t.Parallel()

	type fields struct {
		secret       auth.Secret
		authRPC      rpc.Auth
		userRPC      rpc.User
		authCache    cache.Cache[model.Auth]
		codeCache    cache.Cache[model.Code]
		sessionCache cache.Cache[model.Session]
	}

	type args struct {
		ctx   context.Context
		input usecase.APIAuthVerifyInput
	}

	now := time.Now()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIAuthVerifyOutput
		wantErr bool
	}{
		{
			name: "検証できる",
			fields: fields{
				secret: auth.Secret("secret"),
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					Value: model.Auth{
						AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
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
				input: usecase.APIAuthVerifyInput{
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
			want:    usecase.APIAuthVerifyOutput{},
			wantErr: false,
		},
		{
			name: "AuthがCacheに存在せず検証に失敗する",
			fields: fields{
				secret: auth.Secret("secret"),
				authCache: &cache.CacheMock[model.Auth]{
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
				input: usecase.APIAuthVerifyInput{
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
			want:    usecase.APIAuthVerifyOutput{},
			wantErr: true,
		},
		{
			name: "Authの有効期限が切れて検証に失敗する",
			fields: fields{
				secret: auth.Secret("secret"),
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					Value: model.Auth{
						AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
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
				input: usecase.APIAuthVerifyInput{
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
			want:    usecase.APIAuthVerifyOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			aa := interactor.NewAPIAuth(
				tt.fields.secret,
				tt.fields.authRPC,
				tt.fields.userRPC,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			got, err := aa.Verify(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuth.Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIAuth.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIAuthRefresh(t *testing.T) {
	t.Parallel()

	type fields struct {
		secret       auth.Secret
		authRPC      rpc.Auth
		userRPC      rpc.User
		authCache    cache.Cache[model.Auth]
		codeCache    cache.Cache[model.Code]
		sessionCache cache.Cache[model.Session]
	}

	type args struct {
		ctx   context.Context
		input usecase.APIAuthRefreshInput
	}

	now := time.Now()

	key := generateKey(t)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIAuthRefreshOutput
		wantErr bool
	}{
		{
			name: "リフレッシュできる",
			fields: fields{
				secret: auth.Secret("secret"),
				codeCache: &cache.CacheMock[model.Code]{
					T: t,
					Value: model.Code{
						CodeID:    auth.CodeID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(model.DefaultCodeExpiresIn),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
					DelAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					SetAssert: func(t *testing.T, key string, value model.Auth, ttl time.Duration) {
						t.Helper()
						if !reflect.DeepEqual(ttl, model.DefaultAuthExpiresIn) {
							t.Errorf("ttl = %v, want %v", ttl, model.DefaultAuthExpiresIn)
						}
					},
				},
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					Value: model.Session{
						SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						PublicKey: key.PublicKey,
						IssuedAt:  now,
						ExpiresAt: now.Add(model.DefaultSessionExpiresIn),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthRefreshInput{
					CodeID:    auth.CodeID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					Signature: sign(t, key, "01234567-0123-0123-0123-0123456789ab"),
					SessionToken: auth.GenerateSessionToken(
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.Secret("secret"),
					),
				},
			},
			want:    usecase.APIAuthRefreshOutput{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			aa := interactor.NewAPIAuth(
				tt.fields.secret,
				tt.fields.authRPC,
				tt.fields.userRPC,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			_, err := aa.Refresh(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuth.Refresh() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("APIAuth.Refresh() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestAPIAuthGenerateCode(t *testing.T) {
	t.Parallel()

	type fields struct {
		secret       auth.Secret
		authRPC      rpc.Auth
		userRPC      rpc.User
		authCache    cache.Cache[model.Auth]
		codeCache    cache.Cache[model.Code]
		sessionCache cache.Cache[model.Session]
	}

	type args struct {
		ctx   context.Context
		input usecase.APIAuthGenerateCodeInput
	}

	now := time.Now()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIAuthGenerateCodeOutput
		wantErr bool
	}{
		{
			name: "コードが生成できる",
			fields: fields{
				secret: auth.Secret("secret"),
				codeCache: &cache.CacheMock[model.Code]{
					T: t,
					SetAssert: func(t *testing.T, key string, value model.Code, ttl time.Duration) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthGenerateCodeInput{
					SessionToken: auth.GenerateSessionToken(
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.Secret("secret"),
					),
				},
			},
			want: usecase.APIAuthGenerateCodeOutput{
				Code: model.Code{
					CodeID:    auth.CodeID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					IssuedAt:  now,
					ExpiresAt: now.Add(model.DefaultCodeExpiresIn),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			aa := interactor.NewAPIAuth(
				tt.fields.secret,
				tt.fields.authRPC,
				tt.fields.userRPC,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			got, err := aa.GenerateCode(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuth.GenerateCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Code.SessionID, tt.want.Code.SessionID) {
				t.Errorf("APIAuth.GenerateCode() = %v, want %v", got, tt.want)
			}
		})
	}
}