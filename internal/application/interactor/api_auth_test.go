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
	"github.com/morning-night-guild/platform-app/internal/domain/model/notice"
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

func TestAPIAuthInvite(t *testing.T) {
	t.Parallel()

	type fields struct {
		noticeRPC       func(t *testing.T) rpc.Notice
		authRPC         rpc.Auth
		userRPC         rpc.User
		invitationCache cache.Cache[model.Invitation]
		userCache       cache.Cache[model.User]
		authCache       cache.Cache[model.Auth]
		codeCache       cache.Cache[model.Code]
		sessionCache    cache.Cache[model.Session]
	}

	type args struct {
		ctx   context.Context
		input usecase.APIAuthInviteInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIAuthInviteOutput
		wantErr bool
	}{
		{
			name: "招待できる",
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthInviteInput{
					Email: auth.Email("test@example.com"),
				},
			},
			fields: fields{
				invitationCache: &cache.CacheMock[model.Invitation]{
					T: t,
					SetAssert: func(t *testing.T, key string, value model.Invitation, ttl time.Duration) {
						t.Helper()
						if !reflect.DeepEqual(ttl, 24*time.Hour) {
							t.Errorf("ttl = %v, want %v", ttl, 24*time.Hour)
						}
					},
				},
				noticeRPC: func(t *testing.T) rpc.Notice {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockNotice(ctrl)
					mock.EXPECT().Notify(
						gomock.Any(),
						auth.Email("test@example.com"),
						gomock.Any(),
						gomock.Any(),
					).Return(notice.ID(""), nil)
					return mock
				},
			},
			want:    usecase.APIAuthInviteOutput{},
			wantErr: false,
		},
		{
			name: "InvitationCache.Set()でエラーが発生して招待できない",
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthInviteInput{
					Email: auth.Email("test@example.com"),
				},
			},
			fields: fields{
				invitationCache: &cache.CacheMock[model.Invitation]{
					T:      t,
					SetErr: fmt.Errorf("error"),
					SetAssert: func(t *testing.T, key string, value model.Invitation, ttl time.Duration) {
						t.Helper()
						if !reflect.DeepEqual(ttl, 24*time.Hour) {
							t.Errorf("ttl = %v, want %v", ttl, 24*time.Hour)
						}
					},
				},
				noticeRPC: func(t *testing.T) rpc.Notice {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockNotice(ctrl)
					return mock
				},
			},
			want:    usecase.APIAuthInviteOutput{},
			wantErr: true,
		},
		{
			name: "NoticeRPC.Notify()でエラーが発生して招待できない",
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthInviteInput{
					Email: auth.Email("test@example.com"),
				},
			},
			fields: fields{
				invitationCache: &cache.CacheMock[model.Invitation]{
					T: t,
					SetAssert: func(t *testing.T, key string, value model.Invitation, ttl time.Duration) {
						t.Helper()
						if !reflect.DeepEqual(ttl, 24*time.Hour) {
							t.Errorf("ttl = %v, want %v", ttl, 24*time.Hour)
						}
					},
				},
				noticeRPC: func(t *testing.T) rpc.Notice {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockNotice(ctrl)
					mock.EXPECT().Notify(
						gomock.Any(),
						auth.Email("test@example.com"),
						gomock.Any(),
						gomock.Any(),
					).Return(notice.ID(""), fmt.Errorf("error"))
					return mock
				},
			},
			want:    usecase.APIAuthInviteOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			itr := interactor.NewAPIAuth(
				tt.fields.noticeRPC(t),
				tt.fields.authRPC,
				tt.fields.userRPC,
				tt.fields.invitationCache,
				tt.fields.userCache,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			got, err := itr.Invite(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuth.Invite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got.InvitationCode.String()) != 8 {
				t.Errorf("APIAuth.Invite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIAuthJoin(t *testing.T) {
	t.Parallel()

	type fields struct {
		noticeRPC       rpc.Notice
		authRPC         func(t *testing.T) rpc.Auth
		userRPC         func(t *testing.T) rpc.User
		invitationCache cache.Cache[model.Invitation]
		userCache       cache.Cache[model.User]
		authCache       cache.Cache[model.Auth]
		codeCache       cache.Cache[model.Code]
		sessionCache    cache.Cache[model.Session]
	}

	type args struct {
		ctx   context.Context
		input usecase.APIAuthJoinInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIAuthJoinOutput
		wantErr bool
	}{
		{
			name: "参加できる",
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthJoinInput{
					InvitationCode: auth.InvitationCode("01234567"),
					Password:       auth.Password("password"),
				},
			},
			fields: fields{
				invitationCache: &cache.CacheMock[model.Invitation]{
					T: t,
					Value: model.Invitation{
						Code:  auth.InvitationCode("01234567"),
						Email: auth.Email("test@example.com"),
					},
					GetDelAssert: func(t *testing.T, key string) {
						t.Helper()
						if key != "01234567" {
							t.Errorf("key = %v, want %v", key, "01234567")
						}
					},
				},
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
						auth.Email("test@example.com"),
						auth.Password("password"),
					).Return(nil)
					return mock
				},
			},
			want:    usecase.APIAuthJoinOutput{},
			wantErr: false,
		},
		{
			name: "InvitationCache.Get()でエラーが発生して参加できない",
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthJoinInput{
					InvitationCode: auth.InvitationCode("01234567"),
					Password:       auth.Password("password"),
				},
			},
			fields: fields{
				invitationCache: &cache.CacheMock[model.Invitation]{
					T: t,
					GetDelAssert: func(t *testing.T, key string) {
						t.Helper()
						if key != "01234567" {
							t.Errorf("key = %v, want %v", key, "01234567")
						}
					},
					GetDelErr: fmt.Errorf("error"),
				},
				userRPC: func(t *testing.T) rpc.User {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockUser(ctrl)
					return mock
				},
				authRPC: func(t *testing.T) rpc.Auth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockAuth(ctrl)
					return mock
				},
			},
			want:    usecase.APIAuthJoinOutput{},
			wantErr: true,
		},
		{
			name: "UserRPC.CreateUser()でエラーが発生して参加できない",
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthJoinInput{
					InvitationCode: auth.InvitationCode("01234567"),
					Password:       auth.Password("password"),
				},
			},
			fields: fields{
				invitationCache: &cache.CacheMock[model.Invitation]{
					T: t,
					Value: model.Invitation{
						Code:  auth.InvitationCode("01234567"),
						Email: auth.Email("test@example.com"),
					},
					GetDelAssert: func(t *testing.T, key string) {
						t.Helper()
						if key != "01234567" {
							t.Errorf("key = %v, want %v", key, "01234567")
						}
					},
				},
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
			want:    usecase.APIAuthJoinOutput{},
			wantErr: true,
		},
		{
			name: "AuthRPC.SignUp()でエラーが発生して参加できない",
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthJoinInput{
					InvitationCode: auth.InvitationCode("01234567"),
					Password:       auth.Password("password"),
				},
			},
			fields: fields{
				invitationCache: &cache.CacheMock[model.Invitation]{
					T: t,
					Value: model.Invitation{
						Code:  auth.InvitationCode("01234567"),
						Email: auth.Email("test@example.com"),
					},
					GetDelAssert: func(t *testing.T, key string) {
						t.Helper()
						if key != "01234567" {
							t.Errorf("key = %v, want %v", key, "01234567")
						}
					},
				},
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
						auth.Email("test@example.com"),
						auth.Password("password"),
					).Return(fmt.Errorf("test"))
					return mock
				},
			},
			want:    usecase.APIAuthJoinOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			itr := interactor.NewAPIAuth(
				tt.fields.noticeRPC,
				tt.fields.authRPC(t),
				tt.fields.userRPC(t),
				tt.fields.invitationCache,
				tt.fields.userCache,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			got, err := itr.Join(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuth.Join() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIAuth.Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIAuthSignUp(t *testing.T) {
	t.Parallel()

	type fields struct {
		noticeRPC       rpc.Notice
		authRPC         func(t *testing.T) rpc.Auth
		userRPC         func(t *testing.T) rpc.User
		invitationCache cache.Cache[model.Invitation]
		userCache       cache.Cache[model.User]
		authCache       cache.Cache[model.Auth]
		codeCache       cache.Cache[model.Code]
		sessionCache    cache.Cache[model.Session]
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
					Email:    auth.Email("test@example.com"),
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
						auth.Email("test@example.com"),
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
					Email:    auth.Email("test@example.com"),
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
					Email:    auth.Email("test@example.com"),
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
						auth.Email("test@example.com"),
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
			itr := interactor.NewAPIAuth(
				tt.fields.noticeRPC,
				tt.fields.authRPC(t),
				tt.fields.userRPC(t),
				tt.fields.invitationCache,
				tt.fields.userCache,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			got, err := itr.SignUp(tt.args.ctx, tt.args.input)
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
		noticeRPC       rpc.Notice
		authRPC         func(t *testing.T) rpc.Auth
		userRPC         rpc.User
		invitationCache cache.Cache[model.Invitation]
		userCache       cache.Cache[model.User]
		authCache       cache.Cache[model.Auth]
		codeCache       cache.Cache[model.Code]
		sessionCache    cache.Cache[model.Session]
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
				authRPC: func(t *testing.T) rpc.Auth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockAuth(ctrl)
					mock.EXPECT().SignIn(
						gomock.Any(),
						auth.Email("test@example.com"),
						auth.Password("password"),
					).Return(model.User{
						UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					}, nil)
					return mock
				},
				userCache: &cache.CacheMock[model.User]{
					T: t,
					CreateTxSetCmdAssert: func(t *testing.T, key string, value model.User, ttl time.Duration) {
						t.Helper()
					},
				},
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					CreateTxSetCmdAssert: func(t *testing.T, key string, value model.Auth, ttl time.Duration) {
						t.Helper()
						if !reflect.DeepEqual(ttl, model.DefaultAuthExpiresIn) {
							t.Errorf("ttl = %v, want %v", ttl, model.DefaultAuthExpiresIn)
						}
					},
				},
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					CreateTxSetCmdAssert: func(t *testing.T, key string, value model.Session, ttl time.Duration) {
						t.Helper()
						if !reflect.DeepEqual(ttl, model.DefaultSessionExpiresIn) {
							t.Errorf("ttl = %v, want %v", ttl, model.DefaultSessionExpiresIn)
						}
					},
					TxAssert: func(t *testing.T, setCmds []cache.TxSetCmd, delCmds []cache.TxDelCmd) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthSignInInput{
					Email:     auth.Email("test@example.com"),
					Password:  auth.Password("password"),
					PublicKey: rsa.PublicKey{},
					ExpiresIn: auth.DefaultExpiresIn,
				},
			},
			want: usecase.APIAuthSignInOutput{
				AuthToken: auth.GenerateAuthToken(
					user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret(),
					auth.DefaultExpiresIn,
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
			itr := interactor.NewAPIAuth(
				tt.fields.noticeRPC,
				tt.fields.authRPC(t),
				tt.fields.userRPC,
				tt.fields.invitationCache,
				tt.fields.userCache,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			_, err := itr.SignIn(tt.args.ctx, tt.args.input)
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
		noticeRPC       rpc.Notice
		authRPC         rpc.Auth
		userRPC         rpc.User
		invitationCache cache.Cache[model.Invitation]
		userCache       cache.Cache[model.User]
		authCache       cache.Cache[model.Auth]
		codeCache       cache.Cache[model.Code]
		sessionCache    cache.Cache[model.Session]
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
				userCache: &cache.CacheMock[model.User]{
					T: t,
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
					TxAssert: func(t *testing.T, setCmds []cache.TxSetCmd, delCmds []cache.TxDelCmd) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthSignOutInput{
					UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				},
			},
			want:    usecase.APIAuthSignOutOutput{},
			wantErr: false,
		},
		{
			name: "SessionCache.CreateTxDelCmd()でエラーが発生してもサインアウトできる",
			fields: fields{
				userCache: &cache.CacheMock[model.User]{
					T: t,
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
					CreateTxDelCmdErr: fmt.Errorf("error"),
					TxAssert: func(t *testing.T, setCmds []cache.TxSetCmd, delCmds []cache.TxDelCmd) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthSignOutInput{
					UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				},
			},
			want:    usecase.APIAuthSignOutOutput{},
			wantErr: false,
		},
		{
			name: "AuthCache.CreateTxDelCmd()でエラーが発生してもサインアウトできる",
			fields: fields{
				userCache: &cache.CacheMock[model.User]{
					T: t,
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
					CreateTxDelCmdErr: fmt.Errorf("error"),
				},
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
					TxAssert: func(t *testing.T, setCmds []cache.TxSetCmd, delCmds []cache.TxDelCmd) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthSignOutInput{
					UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				},
			},
			want:    usecase.APIAuthSignOutOutput{},
			wantErr: false,
		},
		{
			name: "UserCache.CreateTxDelCmd()でエラーが発生してもサインアウトできる",
			fields: fields{
				userCache: &cache.CacheMock[model.User]{
					T: t,
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
					CreateTxDelCmdErr: fmt.Errorf("error"),
				},
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
					TxAssert: func(t *testing.T, setCmds []cache.TxSetCmd, delCmds []cache.TxDelCmd) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthSignOutInput{
					UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
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
			itr := interactor.NewAPIAuth(
				tt.fields.noticeRPC,
				tt.fields.authRPC,
				tt.fields.userRPC,
				tt.fields.invitationCache,
				tt.fields.userCache,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			got, err := itr.SignOut(tt.args.ctx, tt.args.input)
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

func TestAPIAuthSignOutAll(t *testing.T) {
	t.Parallel()

	type fields struct {
		noticeRPC       rpc.Notice
		authRPC         rpc.Auth
		userRPC         rpc.User
		invitationCache cache.Cache[model.Invitation]
		userCache       cache.Cache[model.User]
		authCache       cache.Cache[model.Auth]
		codeCache       cache.Cache[model.Code]
		sessionCache    cache.Cache[model.Session]
	}

	type args struct {
		ctx   context.Context
		input usecase.APIAuthSignOutAllInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIAuthSignOutAllOutput
		wantErr bool
	}{
		{
			name: "サインアウトできる",
			fields: fields{
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					KeysValue: []string{
						"session:01234567-0123-0123-0123-0123456789ab:01234567-0123-0123-0123-0123456789ab",
						"session:01234567-0123-0123-0123-0123456789ab:01234567-0123-0123-0123-0123456789ac",
						"session:01234567-0123-0123-0123-0123456789ab:01234567-0123-0123-0123-0123456789ad",
					},
					KeysAssert: func(t *testing.T, pattern string, prefix cache.Prefix) {
						t.Helper()
					},
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
					TxAssert: func(t *testing.T, setCmds []cache.TxSetCmd, delCmds []cache.TxDelCmd) {
						t.Helper()
					},
				},
				userCache: &cache.CacheMock[model.User]{
					T: t,
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthSignOutAllInput{
					UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				},
			},
			want:    usecase.APIAuthSignOutAllOutput{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			itr := interactor.NewAPIAuth(
				tt.fields.noticeRPC,
				tt.fields.authRPC,
				tt.fields.userRPC,
				tt.fields.invitationCache,
				tt.fields.userCache,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			got, err := itr.SignOutAll(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuth.SignOutAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIAuth.SignOutAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIAuthVerify(t *testing.T) {
	t.Parallel()

	type fields struct {
		noticeRPC       rpc.Notice
		authRPC         rpc.Auth
		userRPC         rpc.User
		invitationCache cache.Cache[model.Invitation]
		userCache       cache.Cache[model.User]
		authCache       cache.Cache[model.Auth]
		codeCache       cache.Cache[model.Code]
		sessionCache    cache.Cache[model.Session]
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
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					Value: model.Session{
						SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
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
					UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				},
			},
			want:    usecase.APIAuthVerifyOutput{},
			wantErr: false,
		},
		{
			name: "AuthがCacheに存在せず検証に失敗する",
			fields: fields{
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					Value: model.Session{
						SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
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
					UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				},
			},
			want:    usecase.APIAuthVerifyOutput{},
			wantErr: true,
		},
		{
			name: "Authの有効期限が切れて検証に失敗する",
			fields: fields{
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					Value: model.Session{
						SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
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
					UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				},
			},
			want:    usecase.APIAuthVerifyOutput{},
			wantErr: true,
		},
		{
			name: "SessionがCacheに存在せず検証に失敗する",
			fields: fields{
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					Value: model.Session{
						SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
					GetErr: fmt.Errorf("test"),
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthVerifyInput{
					UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				},
			},
			want:    usecase.APIAuthVerifyOutput{},
			wantErr: true,
		},
		{
			name: "Sessionの有効期限が切れて検証に失敗する",
			fields: fields{
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					Value: model.Session{
						SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
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
					UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
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
			itr := interactor.NewAPIAuth(
				tt.fields.noticeRPC,
				tt.fields.authRPC,
				tt.fields.userRPC,
				tt.fields.invitationCache,
				tt.fields.userCache,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			got, err := itr.Verify(tt.args.ctx, tt.args.input)
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
		noticeRPC       rpc.Notice
		authRPC         rpc.Auth
		userRPC         rpc.User
		invitationCache cache.Cache[model.Invitation]
		userCache       cache.Cache[model.User]
		authCache       cache.Cache[model.Auth]
		codeCache       cache.Cache[model.Code]
		sessionCache    cache.Cache[model.Session]
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
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				userCache: &cache.CacheMock[model.User]{
					T: t,
					Value: model.User{
						UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					CreateTxSetCmdAssert: func(t *testing.T, key string, value model.Auth, ttl time.Duration) {
						t.Helper()
						if !reflect.DeepEqual(ttl, model.DefaultAuthExpiresIn) {
							t.Errorf("ttl = %v, want %v", ttl, model.DefaultAuthExpiresIn)
						}
					},
					TxAssert: func(t *testing.T, setCmds []cache.TxSetCmd, delCmds []cache.TxDelCmd) {
						t.Helper()
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
					SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					ExpiresIn: auth.DefaultExpiresIn,
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
			itr := interactor.NewAPIAuth(
				tt.fields.noticeRPC,
				tt.fields.authRPC,
				tt.fields.userRPC,
				tt.fields.invitationCache,
				tt.fields.userCache,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			_, err := itr.Refresh(tt.args.ctx, tt.args.input)
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
		noticeRPC       rpc.Notice
		authRPC         rpc.Auth
		userRPC         rpc.User
		invitationCache cache.Cache[model.Invitation]
		userCache       cache.Cache[model.User]
		authCache       cache.Cache[model.Auth]
		codeCache       cache.Cache[model.Code]
		sessionCache    cache.Cache[model.Session]
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
					SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
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
			itr := interactor.NewAPIAuth(
				tt.fields.noticeRPC,
				tt.fields.authRPC,
				tt.fields.userRPC,
				tt.fields.invitationCache,
				tt.fields.userCache,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			got, err := itr.GenerateCode(tt.args.ctx, tt.args.input)
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

func TestAPIAuthChangePassword(t *testing.T) {
	t.Parallel()

	type fields struct {
		noticeRPC       rpc.Notice
		authRPC         func(t *testing.T) rpc.Auth
		userRPC         rpc.User
		invitationCache cache.Cache[model.Invitation]
		userCache       cache.Cache[model.User]
		authCache       cache.Cache[model.Auth]
		codeCache       cache.Cache[model.Code]
		sessionCache    cache.Cache[model.Session]
	}

	type args struct {
		ctx   context.Context
		input usecase.APIAuthChangePasswordInput
	}

	now := time.Now()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIAuthChangePasswordOutput
		wantErr bool
	}{
		{
			name: "パスワードを変更できる",
			fields: fields{
				authRPC: func(t *testing.T) rpc.Auth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockAuth(ctrl)
					mock.EXPECT().GetEmail(
						gomock.Any(),
						user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					).Return(auth.Email("test@example.com"), nil)
					mock.EXPECT().SignIn(
						gomock.Any(),
						auth.Email("test@example.com"),
						gomock.Any(),
					).Return(model.User{
						UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					}, nil).Times(2)
					mock.EXPECT().ChangePassword(
						gomock.Any(),
						user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.Password("NewPassword"),
					).Return(nil)
					return mock
				},
				userCache: &cache.CacheMock[model.User]{
					T: t,
					CreateTxSetCmdAssert: func(t *testing.T, key string, value model.User, ttl time.Duration) {
						t.Helper()
					},
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					CreateTxSetCmdAssert: func(t *testing.T, key string, value model.Auth, ttl time.Duration) {
						t.Helper()
						if !reflect.DeepEqual(ttl, model.DefaultAuthExpiresIn) {
							t.Errorf("ttl = %v, want %v", ttl, model.DefaultAuthExpiresIn)
						}
					},
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				sessionCache: &cache.CacheMock[model.Session]{
					T: t,
					CreateTxSetCmdAssert: func(t *testing.T, key string, value model.Session, ttl time.Duration) {
						t.Helper()
						if !reflect.DeepEqual(ttl, model.DefaultSessionExpiresIn) {
							t.Errorf("ttl = %v, want %v", ttl, model.DefaultSessionExpiresIn)
						}
					},
					TxAssert: func(t *testing.T, setCmds []cache.TxSetCmd, delCmds []cache.TxDelCmd) {
						t.Helper()
					},
					KeysValue: []string{
						"session:01234567-0123-0123-0123-0123456789ab:01234567-0123-0123-0123-0123456789ab",
						"session:01234567-0123-0123-0123-0123456789ab:01234567-0123-0123-0123-0123456789ac",
						"session:01234567-0123-0123-0123-0123456789ab:01234567-0123-0123-0123-0123456789ad",
					},
					KeysAssert: func(t *testing.T, pattern string, prefix cache.Prefix) {
						t.Helper()
					},
					CreateTxDelCmdAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIAuthChangePasswordInput{
					Secret:      auth.Secret("secret"),
					UserID:      user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					OldPassword: auth.Password("OldPassword"),
					NewPassword: auth.Password("NewPassword"),
					PublicKey:   rsa.PublicKey{},
					ExpiresIn:   auth.DefaultExpiresIn,
				},
			},
			want: usecase.APIAuthChangePasswordOutput{
				Auth: model.Auth{
					AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					IssuedAt:  now,
					ExpiresAt: now,
				},
				AuthToken: auth.GenerateAuthToken(
					user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret(),
					auth.DefaultExpiresIn,
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
			itr := interactor.NewAPIAuth(
				tt.fields.noticeRPC,
				tt.fields.authRPC(t),
				tt.fields.userRPC,
				tt.fields.invitationCache,
				tt.fields.userCache,
				tt.fields.authCache,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			_, err := itr.ChangePassword(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuth.ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("APIAuth.ChangePassword() = %v, want %v", got, tt.want)
			// }
		})
	}
}
