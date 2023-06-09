package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/adapter/handler"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

const aid = "01234567-0123-0123-0123-0123456789ab"

func TestHandlerV1ArticleList(t *testing.T) {
	t.Parallel()

	toIntPointer := func(v int) *int {
		return &v
	}

	toStringPointer := func(v string) *string {
		return &v
	}

	type fields struct {
		cookie  handler.Cookie
		auth    usecase.APIAuth
		article func(*testing.T) usecase.APIArticle
		health  usecase.APIHealth
	}

	type args struct {
		r       *http.Request
		cookies []*http.Cookie
		params  openapi.V1ArticleListParams
	}

	next := "next"

	token := GenerateToken(t)

	all := openapi.V1ArticleListParamsScope("all")

	own := openapi.V1ArticleListParamsScope("own")

	invalid := openapi.V1ArticleListParamsScope("invalid")

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "記事が一覧できる",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
					mock.EXPECT().List(gomock.Any(), usecase.APIArticleListInput{
						UserID: token.UserID,
						Scope:  article.All,
						Index:  value.Index(0),
						Size:   value.Size(5),
					}).Return(usecase.APIArticleListOutput{
						Articles: []model.Article{},
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
				params: openapi.V1ArticleListParams{
					Scope:       &all,
					PageToken:   &next,
					MaxPageSize: toIntPointer(5),
				},
			},
			status: http.StatusOK,
		},
		{
			name: "ユーザーに紐づく記事が一覧できる",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
					mock.EXPECT().List(gomock.Any(), usecase.APIArticleListInput{
						UserID: token.UserID,
						Scope:  article.Own,
						Index:  value.Index(0),
						Size:   value.Size(5),
					}).Return(usecase.APIArticleListOutput{
						Articles: []model.Article{},
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
				params: openapi.V1ArticleListParams{
					Scope:       &own,
					PageToken:   &next,
					MaxPageSize: toIntPointer(5),
				},
			},
			status: http.StatusOK,
		},
		{
			name: "Scopeを指定しなくても記事が一覧できる",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
					mock.EXPECT().List(gomock.Any(), usecase.APIArticleListInput{
						UserID: token.UserID,
						Scope:  article.All,
						Index:  value.Index(0),
						Size:   value.Size(5),
					}).Return(usecase.APIArticleListOutput{
						Articles: []model.Article{},
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
				params: openapi.V1ArticleListParams{
					Scope:       &all,
					PageToken:   &next,
					MaxPageSize: toIntPointer(5),
				},
			},
			status: http.StatusOK,
		},
		{
			name: "タイトルで部分一致検索して記事が一覧できる",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
					mock.EXPECT().List(gomock.Any(), usecase.APIArticleListInput{
						UserID: token.UserID,
						Scope:  article.All,
						Index:  value.Index(0),
						Size:   value.Size(5),
						Filter: []value.Filter{value.NewFilter("title", "title")},
					}).Return(usecase.APIArticleListOutput{
						Articles: []model.Article{},
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
				params: openapi.V1ArticleListParams{
					Scope:       &all,
					PageToken:   &next,
					MaxPageSize: toIntPointer(5),
					Title:       toStringPointer("title"),
				},
			},
			status: http.StatusOK,
		},
		{
			name: "sizeが不正な値で記事が一覧できない",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
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
				params: openapi.V1ArticleListParams{
					Scope:       &all,
					PageToken:   &next,
					MaxPageSize: toIntPointer(-5),
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "scopeが不正な値で記事が一覧できない",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
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
				params: openapi.V1ArticleListParams{
					Scope:       &invalid,
					PageToken:   &next,
					MaxPageSize: toIntPointer(5),
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "coreにてerrorが発生して記事が一覧できない",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
					mock.EXPECT().List(gomock.Any(), usecase.APIArticleListInput{
						UserID: token.UserID,
						Scope:  article.All,
						Index:  value.Index(0),
						Size:   value.Size(5),
					}).Return(usecase.APIArticleListOutput{}, fmt.Errorf("error"))
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
				params: openapi.V1ArticleListParams{
					Scope:       &all,
					PageToken:   &next,
					MaxPageSize: toIntPointer(5),
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
				tt.fields.cookie,
				tt.fields.auth,
				tt.fields.article(t),
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			for _, cookie := range tt.args.cookies {
				tt.args.r.AddCookie(cookie)
			}
			hdl.V1ArticleList(got, tt.args.r, tt.args.params)
			if got.Code != tt.status {
				t.Errorf("V1ListArticles() = %v, want %v", got.Code, tt.status)
			}
		})
	}
}

func TestHandlerV1V1InternalArticleShare(t *testing.T) {
	t.Parallel()

	toPointer := func(s string) *string {
		return &s
	}

	type fields struct {
		secret  auth.Secret
		cookie  handler.Cookie
		auth    usecase.APIAuth
		article func(*testing.T) usecase.APIArticle
		health  usecase.APIHealth
	}

	type args struct {
		r    *http.Request
		body openapi.V1ArticleShareRequestSchema
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
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
					mock.EXPECT().Share(gomock.Any(), usecase.APIArticleShareInput{
						URL:         article.URL("https://example.com"),
						Title:       article.Title("title"),
						Description: article.Description("description"),
						Thumbnail:   article.Thumbnail("https://example.com"),
					}).Return(usecase.APIArticleShareOutput{
						Article: model.Article{
							ArticleID:   article.ID(uuid.MustParse(aid)),
							URL:         article.URL("https://example.com"),
							Title:       article.Title("title"),
							Description: article.Description("description"),
							Thumbnail:   article.Thumbnail("https://example.com"),
							TagList:     []article.Tag{},
						},
					}, nil)
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
				body: openapi.V1ArticleShareRequestSchema{
					Url:         "https://example.com",
					Title:       toPointer("title"),
					Description: toPointer("description"),
					Thumbnail:   toPointer("https://example.com"),
				},
			},
			status: http.StatusOK,
		},
		{
			name: "nil値が与えられても記事が共有できる",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
					mock.EXPECT().Share(gomock.Any(), usecase.APIArticleShareInput{
						URL:         article.URL("https://example.com"),
						Title:       article.Title(""),
						Description: article.Description(""),
						Thumbnail:   article.Thumbnail(""),
					}).Return(usecase.APIArticleShareOutput{
						Article: model.Article{
							ArticleID:   article.ID(uuid.MustParse(aid)),
							URL:         article.URL("https://example.com"),
							Title:       article.Title(""),
							Description: article.Description(""),
							Thumbnail:   article.Thumbnail(""),
						},
					}, nil)
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
				body: openapi.V1ArticleShareRequestSchema{
					Url:         "https://example.com",
					Title:       nil,
					Description: nil,
					Thumbnail:   nil,
				},
			},
			status: http.StatusOK,
		},
		{
			name: "Api-Keyがなくて記事が共有できない",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"Api-Key": []string{""},
					},
				},
				body: openapi.V1ArticleShareRequestSchema{
					Url:         "https://example.com",
					Title:       toPointer("title"),
					Description: toPointer("description"),
					Thumbnail:   toPointer("https://example.com/thumbnail"),
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
				"key",
				tt.fields.secret,
				tt.fields.cookie,
				tt.fields.auth,
				tt.fields.article(t),
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			buf, _ := json.Marshal(tt.args.body)
			tt.args.r.Body = io.NopCloser(bytes.NewBuffer(buf))
			hdl.V1InternalArticleShare(got, tt.args.r)
			if got.Code != tt.status {
				t.Errorf("V1ShareArticle() = %v, want %v", got.Code, tt.status)
			}
		})
	}
}

func TestHandlerV1InternalArticleDelete(t *testing.T) {
	t.Parallel()

	type fields struct {
		secret  auth.Secret
		cookie  handler.Cookie
		auth    usecase.APIAuth
		article func(*testing.T) usecase.APIArticle
		health  usecase.APIHealth
	}

	type args struct {
		r         *http.Request
		articleID uuid.UUID
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "記事が削除できる",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
					mock.EXPECT().Delete(gomock.Any(), usecase.APIArticleDeleteInput{
						ArticleID: article.ID(uuid.MustParse(aid)),
					}).Return(usecase.APIArticleDeleteOutput{}, nil)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodDelete,
					Header: http.Header{
						"Api-Key": []string{"key"},
					},
					URL: &url.URL{
						Path: "/v1/articles/" + aid,
					},
				},
				articleID: uuid.MustParse(aid),
			},
			status: http.StatusOK,
		},
		{
			name: "Api-Keyがなくて記事が削除できない",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodDelete,
					Header: http.Header{
						"Api-Key": []string{""},
					},
					URL: &url.URL{
						Path: "/v1/articles/" + aid,
					},
				},
				articleID: uuid.MustParse(aid),
			},
			status: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hdl := handler.New(
				"key",
				tt.fields.secret,
				tt.fields.cookie,
				tt.fields.auth,
				tt.fields.article(t),
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			hdl.V1InternalArticleDelete(got, tt.args.r, tt.args.articleID)
			if got.Code != tt.status {
				t.Errorf("V1DeleteArticle() = %v, want %v", got.Code, tt.status)
			}
		})
	}
}

func TestHandlerV1ArticleAddOwn(t *testing.T) {
	t.Parallel()

	type fields struct {
		cookie  handler.Cookie
		auth    usecase.APIAuth
		article func(*testing.T) usecase.APIArticle
		health  usecase.APIHealth
	}

	type args struct {
		r         *http.Request
		cookies   []*http.Cookie
		articleID types.UUID
	}

	token := GenerateToken(t)

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "自身に記事が追加できる",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
					mock.EXPECT().AddToUser(
						gomock.Any(),
						usecase.APIArticleAddToUserInput{
							ArticleID: article.ID(uuid.MustParse(aid)),
							UserID:    token.UserID,
						},
					).Return(usecase.APIArticleAddToUserOutput{}, nil)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
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
				articleID: uuid.MustParse(aid),
			},
			status: http.StatusOK,
		},
		{
			name: "認証に失敗して記事が追加できない",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{},
				},
				cookies:   []*http.Cookie{},
				articleID: uuid.MustParse(aid),
			},
			status: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hdl := handler.New(
				"key",
				auth.Secret("secret"),
				tt.fields.cookie,
				tt.fields.auth,
				tt.fields.article(t),
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			for _, cookie := range tt.args.cookies {
				tt.args.r.AddCookie(cookie)
			}
			hdl.V1ArticleAddOwn(got, tt.args.r, tt.args.articleID)
			if got.Code != tt.status {
				t.Errorf("V1ArticleAddOwn() = %v, want %v", got.Code, tt.status)
			}
		})
	}
}

func TestHandlerV1ArticleRemoveOwn(t *testing.T) {
	t.Parallel()

	type fields struct {
		cookie  handler.Cookie
		auth    usecase.APIAuth
		article func(*testing.T) usecase.APIArticle
		health  usecase.APIHealth
	}

	type args struct {
		r         *http.Request
		cookies   []*http.Cookie
		articleID types.UUID
	}

	token := GenerateToken(t)

	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "自身の記事が削除できる",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
					mock.EXPECT().RemoveFromUser(
						gomock.Any(),
						usecase.APIArticleRemoveFromUserInput{
							ArticleID: article.ID(uuid.MustParse(aid)),
							UserID:    token.UserID,
						},
					).Return(usecase.APIArticleRemoveFromUserOutput{}, nil)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
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
				articleID: uuid.MustParse(aid),
			},
			status: http.StatusOK,
		},
		{
			name: "認証に失敗して記事が追加できない",
			fields: fields{
				article: func(t *testing.T) usecase.APIArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockAPIArticle(ctrl)
					return mock
				},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{},
				},
				cookies:   []*http.Cookie{},
				articleID: uuid.MustParse(aid),
			},
			status: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hdl := handler.New(
				"key",
				auth.Secret("secret"),
				tt.fields.cookie,
				tt.fields.auth,
				tt.fields.article(t),
				tt.fields.health,
			)
			got := httptest.NewRecorder()
			for _, cookie := range tt.args.cookies {
				tt.args.r.AddCookie(cookie)
			}
			hdl.V1ArticleRemoveOwn(got, tt.args.r, tt.args.articleID)
			if got.Code != tt.status {
				t.Errorf("V1ArticleRemoveOwn() = %v, want %v", got.Code, tt.status)
			}
		})
	}
}
