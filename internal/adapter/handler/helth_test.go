package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
		auth    func(*testing.T) usecase.APIAuth
		article func(*testing.T) usecase.APIArticle
		health  func(*testing.T) usecase.APIHealth
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
				key: "key",
				health: func(t *testing.T) usecase.APIHealth {
					t.Helper()
					return nil
				},
			},
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
				tt.fields.auth(t),
				tt.fields.article(t),
				tt.fields.health(t),
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
		auth    func(*testing.T) usecase.APIAuth
		article func(*testing.T) usecase.APIArticle
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
					return nil
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
					return nil
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
				tt.fields.auth(t),
				tt.fields.article(t),
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
