package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/adapter/handler"
	"github.com/morning-night-guild/platform-app/internal/adapter/mock"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

func TestAPIV1HealthAPI(t *testing.T) {
	t.Parallel()

	type fields struct {
		key     string
		secret  auth.Secret
		auth    *handler.Auth
		article *handler.Article
		health  *handler.Health
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
				article: handler.NewArticle(mock.APIArticleList{
					T:        t,
					Articles: []model.Article{},
				}, mock.APIArticleShare{
					T: t,
				}),
				health: handler.NewHealth(&mock.APIHealthCheck{
					T: t,
				}),
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
			rest := handler.New(
				tt.fields.key,
				tt.fields.secret,
				tt.fields.auth,
				tt.fields.article,
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			rest.V1HealthAPI(got, tt.args.r)
			if got.Code != tt.status {
				t.Errorf("V1HealthAPI() = %v, want %v", got.Code, tt.status)
			}
		})
	}
}

func TestAPIV1HealthCore(t *testing.T) {
	t.Parallel()

	type fields struct {
		key     string
		secret  auth.Secret
		auth    *handler.Auth
		article *handler.Article
		health  *handler.Health
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
				article: handler.NewArticle(mock.APIArticleList{
					T:        t,
					Articles: []model.Article{},
				}, mock.APIArticleShare{
					T: t,
				}),
				health: handler.NewHealth(&mock.APIHealthCheck{
					T: t,
				}),
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
				key: "key",
				article: handler.NewArticle(mock.APIArticleList{
					T:        t,
					Articles: []model.Article{},
				}, mock.APIArticleShare{
					T: t,
				}),
				health: handler.NewHealth(&mock.APIHealthCheck{
					T:   t,
					Err: errors.New("error"),
				}),
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
			rest := handler.New(
				tt.fields.key,
				tt.fields.secret,
				tt.fields.auth,
				tt.fields.article,
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			rest.V1HealthCore(got, tt.args.r)
			if got.Code != tt.status {
				t.Errorf("V1HealthCore() = %v, want %v", got.Code, tt.status)
			}
		})
	}
}
