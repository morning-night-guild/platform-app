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
	"github.com/morning-night-guild/platform-app/internal/adapter/mock"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

func TestAPIV1ListArticles(t *testing.T) {
	t.Parallel()

	type fields struct {
		key     string
		article *handler.Article
		health  *handler.Health
	}

	type args struct {
		r      *http.Request
		params openapi.V1ListArticlesParams
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
				params: openapi.V1ListArticlesParams{
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
				params: openapi.V1ListArticlesParams{
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
				article: handler.NewArticle(mock.APIArticleList{
					T:        t,
					Articles: []model.Article{},
					Err:      errors.New("error"),
				}, mock.APIArticleShare{
					T: t,
				}),
				health: handler.NewHealth(&mock.APIHealthCheck{}),
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
				},
				params: openapi.V1ListArticlesParams{
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
			rest := handler.New(tt.fields.key, tt.fields.article, tt.fields.health)
			got := httptest.NewRecorder()
			rest.V1ListArticles(got, tt.args.r, tt.args.params)
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
		article *handler.Article
		health  *handler.Health
	}

	type args struct {
		body openapi.V1ShareArticleJSONRequestBody
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
				article: handler.NewArticle(mock.APIArticleList{
					T: t,
				}, mock.APIArticleShare{
					T: t,
				}),
				health: handler.NewHealth(&mock.APIHealthCheck{
					T: t,
				}),
			},
			args: args{
				body: openapi.V1ShareArticleRequest{
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
				article: handler.NewArticle(mock.APIArticleList{
					T: t,
				}, mock.APIArticleShare{
					T: t,
				}),
				health: handler.NewHealth(&mock.APIHealthCheck{
					T: t,
				}),
			},
			args: args{
				body: openapi.V1ShareArticleRequest{
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
				article: handler.NewArticle(mock.APIArticleList{
					T: t,
				}, mock.APIArticleShare{
					T: t,
				}),
				health: handler.NewHealth(&mock.APIHealthCheck{
					T: t,
				}),
			},
			args: args{
				body: openapi.V1ShareArticleRequest{
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
			rest := handler.New(tt.fields.key, tt.fields.article, tt.fields.health)
			got := httptest.NewRecorder()
			buf, _ := json.Marshal(tt.args.body)
			tt.args.r.Body = io.NopCloser(bytes.NewBuffer(buf))
			rest.V1ShareArticle(got, tt.args.r)
			if got.Code != tt.status {
				t.Errorf("V1ShareArticle() = %v, want %v", got.Code, tt.status)
			}
		})
	}
}
