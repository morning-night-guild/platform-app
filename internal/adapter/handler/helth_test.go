package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/morning-night-guild/platform-app/internal/adapter/handler"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

func TestHandlerV1HealthAPI(t *testing.T) {
	t.Parallel()

	type fields struct {
		key     string
		secret  auth.Secret
		cookie  handler.Cookie
		auth    usecase.APIAuth
		article usecase.APIArticle
		health  usecase.APIHealth
	}

	type args struct {
		r *http.Request
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name:   "ヘルスチェックが成功する",
			fields: fields{},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
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
				tt.fields.cookie,
				tt.fields.auth,
				tt.fields.article,
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			hdl.V1HealthAPI(got, tt.args.r)
			if got.Code != tt.status {
				t.Errorf("V1HealthAPI() = %v, want %v", got.Code, tt.status)
			}
		})
	}
}

func TestHandlerV1HealthCore(t *testing.T) {
	t.Parallel()

	type fields struct {
		key     string
		secret  auth.Secret
		cookie  handler.Cookie
		auth    usecase.APIAuth
		article usecase.APIArticle
		health  func(t *testing.T) usecase.APIHealth
	}

	type args struct {
		r *http.Request
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "ヘルスチェックが成功する",
			fields: fields{
				health: func(t *testing.T) usecase.APIHealth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIHealth(ctrl)
					mock.EXPECT().Check(gomock.Any(), usecase.APIHealthCheckInput{}).Return(usecase.APIHealthCheckOutput{}, nil)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
				},
			},
			status: http.StatusOK,
		},
		{
			name: "Coreでエラーが発生してヘルスチェックが失敗する",
			fields: fields{
				health: func(t *testing.T) usecase.APIHealth {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIHealth(ctrl)
					mock.EXPECT().Check(gomock.Any(), usecase.APIHealthCheckInput{}).Return(usecase.APIHealthCheckOutput{}, fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
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
				tt.fields.cookie,
				tt.fields.auth,
				tt.fields.article,
				tt.fields.health(t),
			)
			got := httptest.NewRecorder()
			hdl.V1HealthCore(got, tt.args.r)
			if got.Code != tt.status {
				t.Errorf("V1HealthCore() = %v, want %v", got.Code, tt.status)
			}
		})
	}
}
