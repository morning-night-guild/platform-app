package handler_test

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/golang/mock/gomock"
	"github.com/morning-night-guild/platform-app/internal/adapter/handler"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

func Cookie(t *testing.T) *handler.MockCookie {
	t.Helper()

	ctrl := gomock.NewController(t)

	cookie := handler.NewMockCookie(ctrl)
	cookie.EXPECT().Domain().Return("localhost").AnyTimes()
	cookie.EXPECT().Secure().Return(false).AnyTimes()
	cookie.EXPECT().SameSite().Return(http.SameSiteDefaultMode).AnyTimes()

	return cookie
}

func GenerateToken(t *testing.T) struct {
	UserID             user.ID
	AuthToken          auth.AuthToken
	AuthTokenString    string
	SessionToken       auth.SessionToken
	SessionTokenString string
} {
	t.Helper()

	sid := auth.GenerateSessionID()

	st := auth.GenerateSessionToken(sid, auth.Secret("secret"))

	uid := user.GenerateID()

	at := auth.GenerateAuthToken(uid, sid.ToSecret())

	return struct {
		UserID             user.ID
		AuthToken          auth.AuthToken
		AuthTokenString    string
		SessionToken       auth.SessionToken
		SessionTokenString string
	}{
		UserID:             uid,
		AuthToken:          at,
		AuthTokenString:    at.String(),
		SessionToken:       st,
		SessionTokenString: st.String(),
	}
}

func TestHandleConnectError(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		err error
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "400",
			args: args{
				ctx: context.Background(),
				err: connect.NewError(connect.CodeInvalidArgument, errors.New("invalid argument")),
			},
			want: http.StatusBadRequest,
		},
		{
			name: "500",
			args: args{
				ctx: context.Background(),
				err: errors.New("unknown error"),
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rest := handler.New("", auth.Secret(""), nil, nil, nil, nil)
			if got := rest.HandleConnectError(tt.args.ctx, tt.args.err); got != tt.want {
				t.Errorf("API.HandleConnectError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlerPointerToString(t *testing.T) {
	t.Parallel()

	type args struct {
		s *string
	}

	test := "test"

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil",
			args: args{
				s: nil,
			},
			want: "",
		},
		{
			name: "not nil",
			args: args{
				s: &test,
			},
			want: "test",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hdl := handler.New(
				"",
				auth.Secret(""),
				nil,
				nil,
				nil,
				nil,
			)
			if got := hdl.PointerToString(tt.args.s); got != tt.want {
				t.Errorf("API.PointerToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlerExtractUserID(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx     context.Context
		r       *http.Request
		cookies []*http.Cookie
	}

	token := GenerateToken(t)

	tests := []struct {
		name    string
		args    args
		want    user.ID
		wantErr bool
	}{
		{
			name: "cookieからUserIDを取り出せる",
			args: args{
				ctx: context.Background(),
				r: &http.Request{
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
			want:    token.UserID,
			wantErr: false,
		},
		{
			name: "Sessionトークンが存在せずUserIDが取り出せない",
			args: args{
				ctx: context.Background(),
				r: &http.Request{
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.AuthTokenKey,
						Value: token.AuthTokenString,
					},
				},
			},
			want:    user.GenerateZeroID(),
			wantErr: true,
		},
		{
			name: "Sessionトークンは存在するが不正な値でUserIDが取り出せない",
			args: args{
				ctx: context.Background(),
				r: &http.Request{
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.SessionTokenKey,
						Value: "token",
					},
				},
			},
			want:    user.GenerateZeroID(),
			wantErr: true,
		},
		{
			name: "Authトークンが存在せずUserIDが取り出せない",
			args: args{
				ctx: context.Background(),
				r: &http.Request{
					Header: http.Header{},
				},
				cookies: []*http.Cookie{
					{
						Name:  auth.SessionTokenKey,
						Value: token.SessionTokenString,
					},
				},
			},
			want:    user.GenerateZeroID(),
			wantErr: true,
		},
		{
			name: "Authトークンは存在するが不正な値でUserIDが取り出せない",
			args: args{
				ctx: context.Background(),
				r: &http.Request{
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
			want:    user.GenerateZeroID(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hdl := handler.New(
				"",
				auth.Secret("secret"),
				nil,
				nil,
				nil,
				nil,
			)
			for _, cookie := range tt.args.cookies {
				tt.args.r.AddCookie(cookie)
			}
			got, err := hdl.ExtractUserID(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.ExtractUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.ExtractUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}
