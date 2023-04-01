package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/adapter/handler"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

func TestAPIV1ListArticles(t *testing.T) {
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
		params openapi.V1ArticleListParams
	}

	next := "next"

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "記事が一覧できる",
			fields: fields{
				key: "key",
				article: handler.NewArticle(
					port.APIArticleListMock{
						T:        t,
						Articles: []model.Article{},
					},
					port.APIArticleShareMock{
						T: t,
					},
				),
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
				},
				params: openapi.V1ArticleListParams{
					PageToken:   &next,
					MaxPageSize: 5,
				},
			},
			status: http.StatusOK,
		},
		{
			name: "sizeが不正な値で記事が一覧できない",
			fields: fields{
				key: "key",
				article: handler.NewArticle(
					port.APIArticleListMock{
						T:        t,
						Articles: []model.Article{},
					},
					port.APIArticleShareMock{
						T: t,
					},
				),
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
				},
				params: openapi.V1ArticleListParams{
					PageToken:   &next,
					MaxPageSize: -5,
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "coreにてerrorが発生して記事が一覧できない",
			fields: fields{
				key: "key",
				article: handler.NewArticle(
					port.APIArticleListMock{
						T:        t,
						Articles: []model.Article{},
						Err:      errors.New("error"),
					},
					port.APIArticleShareMock{
						T: t,
					},
				),
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
				},
				params: openapi.V1ArticleListParams{
					PageToken:   &next,
					MaxPageSize: 5,
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
			hdl.V1ArticleList(got, tt.args.r, tt.args.params)
			if got.Code != tt.status {
				t.Errorf("V1ListArticles() = %v, want %v", got.Code, tt.status)
			}
		})
	}
}

func TestAPIV1ShareArticle(t *testing.T) {
	t.Parallel()

	toPointer := func(s string) *string {
		return &s
	}

	type fields struct {
		key     string
		secret  auth.Secret
		auth    *handler.Auth
		article *handler.Article
		health  *handler.Health
	}

	type args struct {
		body openapi.V1ArticleShareRequestSchema
		r    *http.Request
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "記事が共有できる",
			fields: fields{
				key: "key",
				article: handler.NewArticle(
					port.APIArticleListMock{
						T: t,
					},
					port.APIArticleShareMock{
						T: t,
					},
				),
			},
			args: args{
				body: openapi.V1ArticleShareRequestSchema{
					Url:         "https://example.com",
					Title:       toPointer("title"),
					Description: toPointer("description"),
					Thumbnail:   toPointer("https://example.com/thumbnail"),
				},
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"Api-Key": []string{"key"},
					},
				},
			},
			status: http.StatusOK,
		},
		{
			name: "nil値が与えられても記事が共有できる",
			fields: fields{
				key: "key",
				article: handler.NewArticle(
					port.APIArticleListMock{
						T: t,
					},
					port.APIArticleShareMock{
						T: t,
					},
				),
			},
			args: args{
				body: openapi.V1ArticleShareRequestSchema{
					Url:         "https://example.com",
					Title:       nil,
					Description: nil,
					Thumbnail:   nil,
				},
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"Api-Key": []string{"key"},
					},
				},
			},
			status: http.StatusOK,
		},
		{
			name: "Api-Keyがなくて記事が共有できない",
			fields: fields{
				key: "key",
				article: handler.NewArticle(
					port.APIArticleListMock{
						T: t,
					},
					port.APIArticleShareMock{
						T: t,
					},
				),
			},
			args: args{
				body: openapi.V1ArticleShareRequestSchema{
					Url:         "https://example.com",
					Title:       toPointer("title"),
					Description: toPointer("description"),
					Thumbnail:   toPointer("https://example.com/thumbnail"),
				},
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"Api-Key": []string{""},
					},
				},
			},
			status: http.StatusUnauthorized,
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
			hdl.V1ArticleShare(got, tt.args.r)
			if got.Code != tt.status {
				t.Errorf("V1ShareArticle() = %v, want %v", got.Code, tt.status)
			}
		})
	}
}
