package handler_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/adapter/handler"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

func TestHandlerV1AuthRefresh(t *testing.T) {
	t.Parallel()

	type fields struct {
		key     string
		secret  auth.Secret
		auth    *handler.Auth
		article *handler.Article
		health  *handler.Health
	}

	ctrl := gomock.NewController(t)
	cookie := handler.NewMockCookie(ctrl)
	cookie.EXPECT().Domain().Return("localhost").AnyTimes()
	cookie.EXPECT().Secure().Return(false).AnyTimes()
	cookie.EXPECT().HTTPOnly().Return(true).AnyTimes()
	cookie.EXPECT().SameSite().Return(http.SameSiteDefaultMode).AnyTimes()

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
					cookie,
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
			log.Printf("got: %v", got.Header())
		})
	}
}
