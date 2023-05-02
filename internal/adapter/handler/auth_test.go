package handler_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/adapter/handler"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

const (
	sid = "01234567-0123-0123-0123-0123456789ab"
	cid = "01234567-0123-0123-0123-0123456789ac"
)

type Public struct {
	T   *testing.T
	Key rsa.PublicKey
}

func NewPublicKey(t *testing.T) Public {
	t.Helper()

	prv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	return Public{
		T:   t,
		Key: prv.PublicKey,
	}
}

func (pub Public) String() string {
	pub.T.Helper()

	b := new(bytes.Buffer)

	bt, err := x509.MarshalPKIXPublicKey(&pub.Key)
	if err != nil {
		pub.T.Fatal(err)
	}

	if err := pem.Encode(b, &pem.Block{
		Bytes: bt,
	}); err != nil {
		pub.T.Fatal(err)
	}

	remove := func(arr []string, i int) []string {
		return append(arr[:i], arr[i+1:]...)
	}

	pems := strings.Split(b.String(), "\n")

	pems = remove(pems, len(pems)-1)

	pems = remove(pems, len(pems)-1)

	pems = remove(pems, 0)

	return strings.Join(pems, "")
}

func TestHandlerV1AuthRefresh(t *testing.T) {
	t.Parallel()

	type fields struct {
		auth    func(*testing.T) usecase.APIAuth
		article usecase.APIArticle
		health  usecase.APIHealth
	}

	type args struct {
		r      *http.Request
		cookie *http.Cookie
		params openapi.V1AuthRefreshParams
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "トークンが更新できる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().Refresh(gomock.Any(), usecase.APIAuthRefreshInput{
						CodeID:    auth.CodeID(uuid.MustParse(cid)),
						Signature: auth.Signature("signature"),
						SessionID: auth.SessionID(uuid.MustParse(sid)),
						ExpiresIn: auth.DefaultExpiresIn,
					}).Return(usecase.APIAuthRefreshOutput{}, nil)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookie: &http.Cookie{
					Name:  auth.SessionTokenKey,
					Value: auth.GenerateSessionToken(auth.SessionID(uuid.MustParse(sid)), "secret").String(),
				},
				params: openapi.V1AuthRefreshParams{
					Signature: "signature",
					Code:      uuid.MustParse(cid).String(),
				},
			},
			status: http.StatusOK,
		},
		{
			name: "Codeが不正な値で更新できない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookie: &http.Cookie{
					Name:  auth.SessionTokenKey,
					Value: auth.GenerateSessionToken(auth.SessionID(uuid.MustParse(sid)), "secret").String(),
				},
				params: openapi.V1AuthRefreshParams{
					Signature: "signature",
					Code:      "code",
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "usecaseでエラーが発生して更新できない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().Refresh(gomock.Any(), usecase.APIAuthRefreshInput{
						CodeID:    auth.CodeID(uuid.MustParse(cid)),
						Signature: auth.Signature("signature"),
						SessionID: auth.SessionID(uuid.MustParse(sid)),
						ExpiresIn: auth.DefaultExpiresIn,
					}).Return(usecase.APIAuthRefreshOutput{}, fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookie: &http.Cookie{
					Name:  auth.SessionTokenKey,
					Value: auth.GenerateSessionToken(auth.SessionID(uuid.MustParse(sid)), "secret").String(),
				},
				params: openapi.V1AuthRefreshParams{
					Signature: "signature",
					Code:      uuid.MustParse(cid).String(),
				},
			},
			status: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hdl := handler.New(
				"key",
				auth.Secret("secret"),
				Cookie(t),
				tt.fields.auth(t),
				tt.fields.article,
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			tt.args.r.AddCookie(tt.args.cookie)
			hdl.V1AuthRefresh(got, tt.args.r, tt.args.params)
			if got.Code != tt.status {
				t.Errorf("got %v, want %v", got.Code, tt.status)
			}
		})
	}
}

func TestHandlerV1AuthSignIn(t *testing.T) {
	t.Parallel()

	type fields struct {
		auth    func(*testing.T) usecase.APIAuth
		article usecase.APIArticle
		health  usecase.APIHealth
	}

	type args struct {
		r    *http.Request
		body openapi.V1AuthSignInRequestSchema
	}

	pubkey := NewPublicKey(t)

	expiresIn := 600

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "サインインできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().SignIn(gomock.Any(), usecase.APIAuthSignInInput{
						Secret:    auth.Secret("secret"),
						EMail:     auth.EMail("test@example.com"),
						Password:  auth.Password("password"),
						PublicKey: pubkey.Key,
						ExpiresIn: auth.DefaultExpiresIn,
					}).Return(usecase.APIAuthSignInOutput{}, nil)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
				},
				body: openapi.V1AuthSignInRequestSchema{
					Email:     "test@example.com",
					Password:  "password",
					PublicKey: pubkey.String(),
				},
			},
			status: http.StatusOK,
		},
		{
			name: "有効期限を指定してサインインできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().SignIn(gomock.Any(), usecase.APIAuthSignInInput{
						Secret:    auth.Secret("secret"),
						EMail:     auth.EMail("test@example.com"),
						Password:  auth.Password("password"),
						PublicKey: pubkey.Key,
						ExpiresIn: auth.ExpiresIn(600),
					}).Return(usecase.APIAuthSignInOutput{}, nil)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
				},
				body: openapi.V1AuthSignInRequestSchema{
					Email:     "test@example.com",
					Password:  "password",
					PublicKey: pubkey.String(),
					ExpiresIn: &expiresIn,
				},
			},
			status: http.StatusOK,
		},
		{
			name: "メールアドレスが不正な値でサインインできない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
				},
				body: openapi.V1AuthSignInRequestSchema{
					Email:     "email",
					Password:  "password",
					PublicKey: pubkey.String(),
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "パスワードが不正な値でサインインできない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
				},
				body: openapi.V1AuthSignInRequestSchema{
					Email:     "test@example.com",
					Password:  "",
					PublicKey: pubkey.String(),
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "公開鍵が不正な値でサインインできない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
				},
				body: openapi.V1AuthSignInRequestSchema{
					Email:     "test@example.com",
					Password:  "password",
					PublicKey: "key",
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "usecaseでエラーが発生してサインインできない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().SignIn(gomock.Any(), usecase.APIAuthSignInInput{
						Secret:    auth.Secret("secret"),
						EMail:     auth.EMail("test@example.com"),
						Password:  auth.Password("password"),
						PublicKey: pubkey.Key,
						ExpiresIn: auth.DefaultExpiresIn,
					}).Return(usecase.APIAuthSignInOutput{}, fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
				},
				body: openapi.V1AuthSignInRequestSchema{
					Email:     "test@example.com",
					Password:  "password",
					PublicKey: pubkey.String(),
				},
			},
			status: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hdl := handler.New(
				"key",
				auth.Secret("secret"),
				Cookie(t),
				tt.fields.auth(t),
				tt.fields.article,
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			buf, _ := json.Marshal(tt.args.body)
			tt.args.r.Body = io.NopCloser(bytes.NewBuffer(buf))
			hdl.V1AuthSignIn(got, tt.args.r)
			if got.Code != tt.status {
				t.Errorf("got %v, want %v", got.Code, tt.status)
			}
		})
	}
}

func TestHandlerV1AuthSignOut(t *testing.T) {
	t.Parallel()

	type fields struct {
		auth    func(*testing.T) usecase.APIAuth
		article usecase.APIArticle
		health  usecase.APIHealth
	}

	type args struct {
		r       *http.Request
		cookies []*http.Cookie
	}

	token := GenerateToken(t)

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "サインアウトできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().SignOut(gomock.Any(), usecase.APIAuthSignOutInput{
						UserID:    token.UserID,
						SessionID: token.SessionToken.ID(auth.Secret("secret")),
					}).Return(usecase.APIAuthSignOutOutput{}, nil)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.AuthTokenKey,
						Value: token.AuthTokenString,
					},
					{
						Name:  auth.SessionTokenKey,
						Value: token.SessionTokenString,
					},
				},
			},
			status: http.StatusOK,
		},
		{
			name: "AuthCookieがなくてもサインアウトできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.SessionTokenKey,
						Value: token.SessionTokenString,
					},
				},
			},
			status: http.StatusOK,
		},
		{
			name: "AuthCookieが不正な値でもサインアウトできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.SessionTokenKey,
						Value: token.SessionTokenString,
					},
					{
						Name:  auth.AuthTokenKey,
						Value: "token",
					},
				},
			},
			status: http.StatusOK,
		},
		{
			name: "SessionCookieがなくてもサインアウトできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.AuthTokenKey,
						Value: token.AuthTokenString,
					},
				},
			},
			status: http.StatusOK,
		},
		{
			name: "SessionCookieが不正な値でもサインアウトできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.SessionTokenKey,
						Value: "token",
					},
					{
						Name:  auth.AuthTokenKey,
						Value: token.AuthTokenString,
					},
				},
			},
			status: http.StatusOK,
		},
		{
			name: "usecaseでエラーが発生してもサインアウトできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().SignOut(gomock.Any(), usecase.APIAuthSignOutInput{
						UserID:    token.UserID,
						SessionID: token.SessionToken.ID(auth.Secret("secret")),
					}).Return(usecase.APIAuthSignOutOutput{}, fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.AuthTokenKey,
						Value: token.AuthTokenString,
					},
					{
						Name:  auth.SessionTokenKey,
						Value: token.SessionTokenString,
					},
				},
			},
			status: http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hdl := handler.New(
				"key",
				auth.Secret("secret"),
				Cookie(t),
				tt.fields.auth(t),
				tt.fields.article,
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			for _, cookie := range tt.args.cookies {
				tt.args.r.AddCookie(cookie)
			}
			hdl.V1AuthSignOut(got, tt.args.r)
			if got.Code != tt.status {
				t.Errorf("got %v, want %v", got.Code, tt.status)
			}
		})
	}
}

func TestHandlerV1AuthSignOutAll(t *testing.T) {
	t.Parallel()

	type fields struct {
		auth    func(*testing.T) usecase.APIAuth
		article usecase.APIArticle
		health  usecase.APIHealth
	}

	type args struct {
		r       *http.Request
		cookies []*http.Cookie
	}

	token := GenerateToken(t)

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "サインアウトできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().SignOut(gomock.Any(), usecase.APIAuthSignOutInput{
						UserID:    token.UserID,
						SessionID: token.SessionToken.ID(auth.Secret("secret")),
					}).Return(usecase.APIAuthSignOutOutput{}, nil)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.AuthTokenKey,
						Value: token.AuthTokenString,
					},
					{
						Name:  auth.SessionTokenKey,
						Value: token.SessionTokenString,
					},
				},
			},
			status: http.StatusOK,
		},
		{
			name: "AuthCookieがなくてもサインアウトできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.SessionTokenKey,
						Value: token.SessionTokenString,
					},
				},
			},
			status: http.StatusOK,
		},
		{
			name: "AuthCookieが不正な値でもサインアウトできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.SessionTokenKey,
						Value: token.SessionTokenString,
					},
					{
						Name:  auth.AuthTokenKey,
						Value: "token",
					},
				},
			},
			status: http.StatusOK,
		},
		{
			name: "SessionCookieがなくてもサインアウトできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.AuthTokenKey,
						Value: token.AuthTokenString,
					},
				},
			},
			status: http.StatusOK,
		},
		{
			name: "SessionCookieが不正な値でもサインアウトできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.SessionTokenKey,
						Value: "token",
					},
					{
						Name:  auth.AuthTokenKey,
						Value: token.AuthTokenString,
					},
				},
			},
			status: http.StatusOK,
		},
		{
			name: "usecaseでエラーが発生してもサインアウトできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().SignOut(gomock.Any(), usecase.APIAuthSignOutInput{
						UserID:    token.UserID,
						SessionID: token.SessionToken.ID(auth.Secret("secret")),
					}).Return(usecase.APIAuthSignOutOutput{}, fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.AuthTokenKey,
						Value: token.AuthTokenString,
					},
					{
						Name:  auth.SessionTokenKey,
						Value: token.SessionTokenString,
					},
				},
			},
			status: http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hdl := handler.New(
				"key",
				auth.Secret("secret"),
				Cookie(t),
				tt.fields.auth(t),
				tt.fields.article,
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			for _, cookie := range tt.args.cookies {
				tt.args.r.AddCookie(cookie)
			}
			hdl.V1AuthSignOutAll(got, tt.args.r)
			if got.Code != tt.status {
				t.Errorf("got %v, want %v", got.Code, tt.status)
			}
		})
	}
}

func TestHandlerV1AuthSignUp(t *testing.T) {
	t.Parallel()

	type fields struct {
		auth    func(*testing.T) usecase.APIAuth
		article usecase.APIArticle
		health  usecase.APIHealth
	}

	type args struct {
		r    *http.Request
		body openapi.V1AuthSignUpRequestSchema
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "サインアップできる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().SignUp(gomock.Any(), usecase.APIAuthSignUpInput{
						EMail:    auth.EMail("test@example.com"),
						Password: auth.Password("password"),
					}).Return(usecase.APIAuthSignUpOutput{}, nil)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"Api-Key": []string{"key"},
					},
				},
				body: openapi.V1AuthSignUpRequestSchema{
					Email:    "test@example.com",
					Password: "password",
				},
			},
			status: http.StatusOK,
		},
		{
			name: "Api-Keyがなくてサインアップできない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{},
				body: openapi.V1AuthSignUpRequestSchema{
					Email:    "test@example.com",
					Password: "password",
				},
			},
			status: http.StatusUnauthorized,
		},
		{
			name: "メールアドレスが不正な値でサインアップできない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"Api-Key": []string{"key"},
					},
				},
				body: openapi.V1AuthSignUpRequestSchema{
					Email:    "email",
					Password: "password",
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "パスワードが不正な値でサインアップできない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"Api-Key": []string{"key"},
					},
				},
				body: openapi.V1AuthSignUpRequestSchema{
					Email:    "test@example.com",
					Password: "",
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "usecaseでエラーが発生してサインアップできない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().SignUp(gomock.Any(), usecase.APIAuthSignUpInput{
						EMail:    auth.EMail("test@example.com"),
						Password: auth.Password("password"),
					}).Return(usecase.APIAuthSignUpOutput{}, fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"Api-Key": []string{"key"},
					},
				},
				body: openapi.V1AuthSignUpRequestSchema{
					Email:    "test@example.com",
					Password: "password",
				},
			},
			status: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hdl := handler.New(
				"key",
				auth.Secret("secret"),
				Cookie(t),
				tt.fields.auth(t),
				tt.fields.article,
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			buf, _ := json.Marshal(tt.args.body)
			tt.args.r.Body = io.NopCloser(bytes.NewBuffer(buf))
			hdl.V1AuthSignUp(got, tt.args.r)
			if got.Code != tt.status {
				t.Errorf("got %v, want %v", got.Code, tt.status)
			}
		})
	}
}

func TestHandlerV1AuthVerify(t *testing.T) {
	t.Parallel()

	type fields struct {
		auth    func(*testing.T) usecase.APIAuth
		article usecase.APIArticle
		health  usecase.APIHealth
	}

	type args struct {
		r       *http.Request
		cookies []*http.Cookie
	}

	token := GenerateToken(t)

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "トークンを検証できる",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().Verify(gomock.Any(), usecase.APIAuthVerifyInput{
						UserID: token.UserID,
					}).Return(usecase.APIAuthVerifyOutput{}, nil)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.AuthTokenKey,
						Value: token.AuthTokenString,
					},
					{
						Name:  auth.SessionTokenKey,
						Value: token.SessionTokenString,
					},
				},
			},
			status: http.StatusOK,
		},
		{
			name: "Sessionトークンが存在せず検証できない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.AuthTokenKey,
						Value: token.AuthTokenString,
					},
				},
			},
			status: http.StatusUnauthorized,
		},
		{
			name: "Sessionトークンは存在するが不正な値で検証できない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.SessionTokenKey,
						Value: "token",
					},
				},
			},
			status: http.StatusUnauthorized,
		},
		{
			name: "Authトークンが存在せず検証できない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().GenerateCode(gomock.Any(), usecase.APIAuthGenerateCodeInput{
						SessionID: token.SessionToken.ID(auth.Secret("secret")),
					}).Return(usecase.APIAuthGenerateCodeOutput{
						Code: model.Code{
							CodeID:    auth.CodeID(uuid.MustParse(cid)),
							SessionID: token.SessionToken.ID(auth.Secret("secret")),
							IssuedAt:  time.Now(),
							ExpiresAt: time.Now().Add(time.Minute * 10),
						},
					}, nil)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.SessionTokenKey,
						Value: token.SessionTokenString,
					},
				},
			},
			status: http.StatusUnauthorized,
		},
		{
			name: "Authトークンは存在するが不正な値で検証できない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().GenerateCode(gomock.Any(), usecase.APIAuthGenerateCodeInput{
						SessionID: token.SessionToken.ID(auth.Secret("secret")),
					}).Return(usecase.APIAuthGenerateCodeOutput{
						Code: model.Code{
							CodeID:    auth.CodeID(uuid.MustParse(cid)),
							SessionID: token.SessionToken.ID(auth.Secret("secret")),
							IssuedAt:  time.Now(),
							ExpiresAt: time.Now().Add(time.Minute * 10),
						},
					}, nil)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.SessionTokenKey,
						Value: token.SessionTokenString,
					},
					{
						Name:  auth.AuthTokenKey,
						Value: "token",
					},
				},
			},
			status: http.StatusUnauthorized,
		},
		{
			name: "usecaseでエラーが発生して検証できない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().Verify(gomock.Any(), usecase.APIAuthVerifyInput{
						UserID: token.UserID,
					}).Return(usecase.APIAuthVerifyOutput{}, fmt.Errorf("error"))
					mock.EXPECT().GenerateCode(gomock.Any(), usecase.APIAuthGenerateCodeInput{
						SessionID: token.SessionToken.ID(auth.Secret("secret")),
					}).Return(usecase.APIAuthGenerateCodeOutput{
						Code: model.Code{
							CodeID:    auth.CodeID(uuid.MustParse(cid)),
							SessionID: token.SessionToken.ID(auth.Secret("secret")),
							IssuedAt:  time.Now(),
							ExpiresAt: time.Now().Add(time.Minute * 10),
						},
					}, nil)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.AuthTokenKey,
						Value: token.AuthTokenString,
					},
					{
						Name:  auth.SessionTokenKey,
						Value: token.SessionTokenString,
					},
				},
			},
			status: http.StatusUnauthorized,
		},
		{
			name: "code生成でエラーが発生して検証できない",
			fields: fields{
				auth: func(t *testing.T) usecase.APIAuth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIAuth(ctrl)
					mock.EXPECT().Verify(gomock.Any(), usecase.APIAuthVerifyInput{
						UserID: token.UserID,
					}).Return(usecase.APIAuthVerifyOutput{}, fmt.Errorf("error"))
					mock.EXPECT().GenerateCode(gomock.Any(), usecase.APIAuthGenerateCodeInput{
						SessionID: token.SessionToken.ID(auth.Secret("secret")),
					}).Return(usecase.APIAuthGenerateCodeOutput{}, fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.AuthTokenKey,
						Value: token.AuthTokenString,
					},
					{
						Name:  auth.SessionTokenKey,
						Value: token.SessionTokenString,
					},
				},
			},
			status: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hdl := handler.New(
				"key",
				auth.Secret("secret"),
				Cookie(t),
				tt.fields.auth(t),
				tt.fields.article,
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			for _, cookie := range tt.args.cookies {
				tt.args.r.AddCookie(cookie)
			}
			hdl.V1AuthVerify(got, tt.args.r)
			if got.Code != tt.status {
				t.Errorf("got %v, want %v", got.Code, tt.status)
			}
		})
	}
}
