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

	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/adapter/handler"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

func Cookie(t *testing.T) *handler.MockCookie {
	t.Helper()

	ctrl := gomock.NewController(t)

	cookie := handler.NewMockCookie(ctrl)
	cookie.EXPECT().Domain().Return("localhost").AnyTimes()
	cookie.EXPECT().Secure().Return(false).AnyTimes()
	cookie.EXPECT().HTTPOnly().Return(true).AnyTimes()
	cookie.EXPECT().SameSite().Return(http.SameSiteDefaultMode).AnyTimes()

	return cookie
}

func TestHandlerV1AuthRefresh(t *testing.T) {
	t.Parallel()

	type fields struct {
		key     string
		secret  auth.Secret
		auth    *handler.Auth
		article *handler.Article
		health  *handler.Health
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
				secret: auth.Secret("secret"),
				auth: handler.NewAuth(
					&port.APIAuthSignUpMock{},
					&port.APIAuthSignInMock{},
					&port.APIAuthSignOutMock{},
					&port.APIAuthVerifyMock{},
					&port.APIAuthRefreshMock{
						T:         t,
						AuthToken: auth.GenerateAuthToken(user.GenerateID(), "01234567-0123-0123-0123-0123456789ab"),
					},
					&port.APIAuthGenerateCodeMock{},
					Cookie(t),
				),
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookie: &http.Cookie{
					Name:  auth.SessionTokenKey,
					Value: string(auth.GenerateSessionToken(auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")), "secret")),
				},
				params: openapi.V1AuthRefreshParams{
					Signature: "signature",
					Code:      uuid.New().String(),
				},
			},
			status: http.StatusOK,
		},
		{
			name: "Codeが不正な値で更新できない",
			fields: fields{
				secret: auth.Secret("secret"),
				auth: handler.NewAuth(
					&port.APIAuthSignUpMock{},
					&port.APIAuthSignInMock{},
					&port.APIAuthSignOutMock{},
					&port.APIAuthVerifyMock{},
					&port.APIAuthRefreshMock{},
					&port.APIAuthGenerateCodeMock{},
					Cookie(t),
				),
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookie: &http.Cookie{
					Name:  auth.SessionTokenKey,
					Value: string(auth.GenerateSessionToken(auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")), "secret")),
				},
				params: openapi.V1AuthRefreshParams{
					Signature: "signature",
					Code:      "code",
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "usecaseにてエラーが発生して更新できない",
			fields: fields{
				secret: auth.Secret("secret"),
				auth: handler.NewAuth(
					&port.APIAuthSignUpMock{},
					&port.APIAuthSignInMock{},
					&port.APIAuthSignOutMock{},
					&port.APIAuthVerifyMock{},
					&port.APIAuthRefreshMock{
						T:   t,
						Err: fmt.Errorf("test"),
					},
					&port.APIAuthGenerateCodeMock{},
					Cookie(t),
				),
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookie: &http.Cookie{
					Name:  auth.SessionTokenKey,
					Value: string(auth.GenerateSessionToken(auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")), "secret")),
				},
				params: openapi.V1AuthRefreshParams{
					Signature: "signature",
					Code:      uuid.New().String(),
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
				tt.fields.key,
				tt.fields.secret,
				tt.fields.auth,
				tt.fields.article,
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			tt.args.r.AddCookie(&http.Cookie{
				Name:  auth.SessionTokenKey,
				Value: auth.GenerateSessionToken(auth.GenerateSessionID(), "secret").String(),
			})
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
		key     string
		secret  auth.Secret
		auth    *handler.Auth
		article *handler.Article
		health  *handler.Health
	}

	type args struct {
		r    *http.Request
		body openapi.V1AuthSignInRequestSchema
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "サインインできる",
			fields: fields{
				secret: auth.Secret("secret"),
				auth: handler.NewAuth(
					&port.APIAuthSignUpMock{},
					&port.APIAuthSignInMock{
						T:            t,
						AuthToken:    auth.GenerateAuthToken(user.GenerateID(), auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret()),
						SessionToken: auth.GenerateSessionToken(auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")), auth.Secret("secret")),
					},
					&port.APIAuthSignOutMock{},
					&port.APIAuthVerifyMock{},
					&port.APIAuthRefreshMock{},
					&port.APIAuthGenerateCodeMock{},
					Cookie(t),
				),
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
				},
				body: openapi.V1AuthSignInRequestSchema{
					Email:     "test@example.com",
					Password:  "password",
					PublicKey: GeneratePublicKey(t),
				},
			},
			status: http.StatusOK,
		},
		{
			name: "メールアドレスが不正な値でサインインできない",
			fields: fields{
				secret: auth.Secret("secret"),
				auth: handler.NewAuth(
					&port.APIAuthSignUpMock{},
					&port.APIAuthSignInMock{
						T:            t,
						AuthToken:    auth.GenerateAuthToken(user.GenerateID(), auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret()),
						SessionToken: auth.GenerateSessionToken(auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")), auth.Secret("secret")),
					},
					&port.APIAuthSignOutMock{},
					&port.APIAuthVerifyMock{},
					&port.APIAuthRefreshMock{},
					&port.APIAuthGenerateCodeMock{},
					Cookie(t),
				),
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
				},
				body: openapi.V1AuthSignInRequestSchema{
					Email:     "email",
					Password:  "password",
					PublicKey: GeneratePublicKey(t),
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "パスワードが不正な値でサインインできない",
			fields: fields{
				secret: auth.Secret("secret"),
				auth: handler.NewAuth(
					&port.APIAuthSignUpMock{},
					&port.APIAuthSignInMock{
						T:            t,
						AuthToken:    auth.GenerateAuthToken(user.GenerateID(), auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret()),
						SessionToken: auth.GenerateSessionToken(auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")), auth.Secret("secret")),
					},
					&port.APIAuthSignOutMock{},
					&port.APIAuthVerifyMock{},
					&port.APIAuthRefreshMock{},
					&port.APIAuthGenerateCodeMock{},
					Cookie(t),
				),
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
				},
				body: openapi.V1AuthSignInRequestSchema{
					Email:     "test@example.com",
					Password:  "",
					PublicKey: GeneratePublicKey(t),
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "公開鍵が不正な値でサインインできない",
			fields: fields{
				secret: auth.Secret("secret"),
				auth: handler.NewAuth(
					&port.APIAuthSignUpMock{},
					&port.APIAuthSignInMock{
						T:            t,
						AuthToken:    auth.GenerateAuthToken(user.GenerateID(), auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret()),
						SessionToken: auth.GenerateSessionToken(auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")), auth.Secret("secret")),
					},
					&port.APIAuthSignOutMock{},
					&port.APIAuthVerifyMock{},
					&port.APIAuthRefreshMock{},
					&port.APIAuthGenerateCodeMock{},
					Cookie(t),
				),
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hdl := handler.New(
				tt.fields.key,
				tt.fields.secret,
				tt.fields.auth,
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

func GeneratePublicKey(t *testing.T) string {
	t.Helper()

	prv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	b := new(bytes.Buffer)

	bt, err := x509.MarshalPKIXPublicKey(&prv.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	if err := pem.Encode(b, &pem.Block{
		Bytes: bt,
	}); err != nil {
		t.Fatal(err)
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

func TestHandlerV1AuthSignOut(t *testing.T) {
	t.Parallel()

	type fields struct {
		key     string
		secret  auth.Secret
		auth    *handler.Auth
		article *handler.Article
		health  *handler.Health
	}

	type args struct {
		r       *http.Request
		cookies []*http.Cookie
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "サインアウトできる",
			fields: fields{
				secret: auth.Secret("secret"),
				auth: handler.NewAuth(
					&port.APIAuthSignUpMock{},
					&port.APIAuthSignInMock{},
					&port.APIAuthSignOutMock{
						T: t,
					},
					&port.APIAuthVerifyMock{},
					&port.APIAuthRefreshMock{},
					&port.APIAuthGenerateCodeMock{},
					Cookie(t),
				),
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.AuthTokenKey,
						Value: auth.GenerateAuthToken(user.GenerateID(), auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret()).String(),
					},
					{
						Name:  auth.SessionTokenKey,
						Value: auth.GenerateSessionToken(auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")), auth.Secret("secret")).String(),
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
				tt.fields.key,
				tt.fields.secret,
				tt.fields.auth,
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

func TestHandlerV1AuthSignUp(t *testing.T) {
	t.Parallel()

	type fields struct {
		key     string
		secret  auth.Secret
		auth    *handler.Auth
		article *handler.Article
		health  *handler.Health
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
				secret: auth.Secret("secret"),
				auth: handler.NewAuth(
					&port.APIAuthSignUpMock{
						T: t,
					},
					&port.APIAuthSignInMock{},
					&port.APIAuthSignOutMock{},
					&port.APIAuthVerifyMock{},
					&port.APIAuthRefreshMock{},
					&port.APIAuthGenerateCodeMock{},
					Cookie(t),
				),
			},
			args: args{
				r: &http.Request{},
				body: openapi.V1AuthSignUpRequestSchema{
					Email:    "test@example.com",
					Password: "password",
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
				tt.fields.key,
				tt.fields.secret,
				tt.fields.auth,
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
		key     string
		secret  auth.Secret
		auth    *handler.Auth
		article *handler.Article
		health  *handler.Health
	}

	type args struct {
		r       *http.Request
		cookies []*http.Cookie
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "トークンを検証できる",
			fields: fields{
				secret: auth.Secret("secret"),
				auth: handler.NewAuth(
					&port.APIAuthSignUpMock{},
					&port.APIAuthSignInMock{},
					&port.APIAuthSignOutMock{},
					&port.APIAuthVerifyMock{
						T: t,
					},
					&port.APIAuthRefreshMock{},
					&port.APIAuthGenerateCodeMock{},
					Cookie(t),
				),
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.AuthTokenKey,
						Value: auth.GenerateAuthToken(user.GenerateID(), auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")).ToSecret()).String(),
					},
					{
						Name:  auth.SessionTokenKey,
						Value: auth.GenerateSessionToken(auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")), auth.Secret("secret")).String(),
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
				tt.fields.key,
				tt.fields.secret,
				tt.fields.auth,
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
