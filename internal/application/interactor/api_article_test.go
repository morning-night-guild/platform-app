package interactor_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/application/interactor"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

func TestAPIArticleShare(t *testing.T) {
	t.Parallel()

	type fields struct {
		authCache  cache.Cache[model.Auth]
		articleRPC func(t *testing.T) rpc.Article
	}

	type args struct {
		ctx   context.Context
		input usecase.APIArticleShareInput
	}

	id := article.GenerateID()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIArticleShareOutput
		wantErr bool
	}{
		{
			name: "記事が共有できる",
			fields: fields{
				articleRPC: func(t *testing.T) rpc.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockArticle(ctrl)
					mock.EXPECT().Share(
						gomock.Any(),
						article.URL("https://example.com"),
						article.Title("title"),
						article.Description("description"),
						article.Thumbnail("https://example.com"),
					).Return(model.Article{
						ArticleID:   id,
						Title:       article.Title("title"),
						URL:         article.URL("https://example.com"),
						Description: article.Description("description"),
						Thumbnail:   article.Thumbnail("https://example.com"),
					}, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIArticleShareInput{
					URL:         article.URL("https://example.com"),
					Title:       article.Title("title"),
					Description: article.Description("description"),
					Thumbnail:   article.Thumbnail("https://example.com"),
				},
			},
			want: usecase.APIArticleShareOutput{
				Article: model.Article{
					ArticleID:   id,
					Title:       article.Title("title"),
					URL:         article.URL("https://example.com"),
					Description: article.Description("description"),
					Thumbnail:   article.Thumbnail("https://example.com"),
				},
			},
			wantErr: false,
		},
		{
			name: "rpcでerrorが発生して記事の共有が共有できない",
			fields: fields{
				articleRPC: func(t *testing.T) rpc.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockArticle(ctrl)
					mock.EXPECT().Share(
						gomock.Any(),
						article.URL("https://example.com"),
						article.Title("title"),
						article.Description("description"),
						article.Thumbnail("https://example.com"),
					).Return(model.Article{}, fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIArticleShareInput{
					URL:         article.URL("https://example.com"),
					Title:       article.Title("title"),
					Description: article.Description("description"),
					Thumbnail:   article.Thumbnail("https://example.com"),
				},
			},
			want: usecase.APIArticleShareOutput{
				Article: model.Article{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			itr := interactor.NewAPIArticle(
				tt.fields.authCache,
				tt.fields.articleRPC(t),
			)
			got, err := itr.Share(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIArticle.Share() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIArticle.Share() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIArticleList(t *testing.T) {
	t.Parallel()

	type fields struct {
		authCache  cache.Cache[model.Auth]
		articleRPC func(t *testing.T) rpc.Article
	}

	type args struct {
		ctx   context.Context
		input usecase.APIArticleListInput
	}

	articles := []model.Article{
		{
			ArticleID:   article.GenerateID(),
			Title:       article.Title("title1"),
			URL:         article.URL("https://example.com/1"),
			Description: article.Description("description1"),
			Thumbnail:   article.Thumbnail("https://example.com/2"),
		},
		{
			ArticleID:   article.GenerateID(),
			Title:       article.Title("title2"),
			URL:         article.URL("https://example.com/2"),
			Description: article.Description("description2"),
			Thumbnail:   article.Thumbnail("https://example.com/2"),
		},
	}

	now := time.Now()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIArticleListOutput
		wantErr bool
	}{
		{
			name: "記事リストが取得できる",
			fields: fields{
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					Value: model.Auth{
						AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				articleRPC: func(t *testing.T) rpc.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockArticle(ctrl)
					mock.EXPECT().List(gomock.Any(), value.Index(0), value.Size(2)).Return(articles, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIArticleListInput{
					Scope:  article.All,
					UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					Index:  value.Index(0),
					Size:   value.Size(2),
				},
			},
			want: usecase.APIArticleListOutput{
				Articles: articles,
			},
			wantErr: false,
		},
		{
			name: "ユーザーに紐づく記事が一覧取得できる",
			fields: fields{
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					Value: model.Auth{
						AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				articleRPC: func(t *testing.T) rpc.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockArticle(ctrl)
					mock.EXPECT().ListByUser(
						gomock.Any(),
						user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						value.Index(0),
						value.Size(2),
					).Return(articles, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIArticleListInput{
					Scope:  article.Own,
					UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					Index:  value.Index(0),
					Size:   value.Size(2),
				},
			},
			want: usecase.APIArticleListOutput{
				Articles: articles,
			},
			wantErr: false,
		},
		{
			name: "タイトルの部分一致で記事リストが取得できる",
			fields: fields{
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					Value: model.Auth{
						AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				articleRPC: func(t *testing.T) rpc.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockArticle(ctrl)
					mock.EXPECT().List(gomock.Any(), value.Index(0), value.Size(2), []value.Filter{value.NewFilter("title", "title")}).Return(articles, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIArticleListInput{
					Scope:  article.All,
					UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					Index:  value.Index(0),
					Size:   value.Size(2),
					Filter: []value.Filter{value.NewFilter("title", "title")},
				},
			},
			want: usecase.APIArticleListOutput{
				Articles: articles,
			},
			wantErr: false,
		},
		{
			name: "ユーザーに紐づくタイトルの部分一致で記事一覧が取得できる",
			fields: fields{
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					Value: model.Auth{
						AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				articleRPC: func(t *testing.T) rpc.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockArticle(ctrl)
					mock.EXPECT().ListByUser(
						gomock.Any(),
						user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						value.Index(0),
						value.Size(2),
						[]value.Filter{value.NewFilter("title", "title")},
					).Return(articles, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIArticleListInput{
					Scope:  article.Own,
					UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					Index:  value.Index(0),
					Size:   value.Size(2),
					Filter: []value.Filter{value.NewFilter("title", "title")},
				},
			},
			want: usecase.APIArticleListOutput{
				Articles: articles,
			},
			wantErr: false,
		},
		{
			name: "認証に失敗する",
			fields: fields{
				authCache: &cache.CacheMock[model.Auth]{
					T:     t,
					Value: model.Auth{},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
					GetErr: fmt.Errorf("test"),
				},
				articleRPC: func(t *testing.T) rpc.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockArticle(ctrl)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIArticleListInput{
					Scope:  article.All,
					UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					Index:  value.Index(0),
					Size:   value.Size(2),
				},
			},
			want:    usecase.APIArticleListOutput{},
			wantErr: true,
		},
		{
			name: "rpcでerrorが発生して記事リストが取得できない",
			fields: fields{
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					Value: model.Auth{
						AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				articleRPC: func(t *testing.T) rpc.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockArticle(ctrl)
					mock.EXPECT().List(gomock.Any(), value.Index(0), value.Size(2)).Return(nil, fmt.Errorf("test"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIArticleListInput{
					Scope:  article.All,
					UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					Index:  value.Index(0),
					Size:   value.Size(2),
				},
			},
			want:    usecase.APIArticleListOutput{},
			wantErr: true,
		},
		{
			name: "rpcでerrorが発生してユーザーに紐づく記事一覧が取得できない",
			fields: fields{
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					Value: model.Auth{
						AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				articleRPC: func(t *testing.T) rpc.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockArticle(ctrl)
					mock.EXPECT().ListByUser(
						gomock.Any(),
						user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						value.Index(0),
						value.Size(2),
					).Return(nil, fmt.Errorf("test"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIArticleListInput{
					Scope:  article.Own,
					UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					Index:  value.Index(0),
					Size:   value.Size(2),
				},
			},
			want:    usecase.APIArticleListOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			itr := interactor.NewAPIArticle(
				tt.fields.authCache,
				tt.fields.articleRPC(t),
			)
			got, err := itr.List(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIArticle.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIArticle.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIArticleAddToUser(t *testing.T) {
	t.Parallel()

	type fields struct {
		authCache  cache.Cache[model.Auth]
		articleRPC func(*testing.T) rpc.Article
	}

	type args struct {
		ctx   context.Context
		input usecase.APIArticleAddToUserInput
	}

	now := time.Now()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIArticleAddToUserOutput
		wantErr bool
	}{
		{
			name: "ユーザーに記事を追加できる",
			fields: fields{
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					Value: model.Auth{
						AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				articleRPC: func(t *testing.T) rpc.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockArticle(ctrl)
					mock.EXPECT().AddToUser(
						gomock.Any(),
						article.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ac")),
					).Return(nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIArticleAddToUserInput{
					ArticleID: article.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ac")),
				},
			},
			want:    usecase.APIArticleAddToUserOutput{},
			wantErr: false,
		},
		{
			name: "ユーザーに記事を追加できない",
			fields: fields{
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					Value: model.Auth{
						AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				articleRPC: func(t *testing.T) rpc.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockArticle(ctrl)
					mock.EXPECT().AddToUser(
						gomock.Any(),
						article.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ac")),
					).Return(fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIArticleAddToUserInput{
					ArticleID: article.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ac")),
				},
			},
			want:    usecase.APIArticleAddToUserOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			itr := interactor.NewAPIArticle(
				tt.fields.authCache,
				tt.fields.articleRPC(t),
			)
			got, err := itr.AddToUser(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIArticle.AddToUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIArticle.AddToUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIArticleRemoveFromUser(t *testing.T) {
	t.Parallel()

	type fields struct {
		authCache  cache.Cache[model.Auth]
		articleRPC func(*testing.T) rpc.Article
	}

	type args struct {
		ctx   context.Context
		input usecase.APIArticleRemoveFromUserInput
	}

	now := time.Now()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIArticleRemoveFromUserOutput
		wantErr bool
	}{
		{
			name: "ユーザーの記事を削除できる",
			fields: fields{
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					Value: model.Auth{
						AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				articleRPC: func(t *testing.T) rpc.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockArticle(ctrl)
					mock.EXPECT().RemoveFromUser(
						gomock.Any(),
						article.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ac")),
					).Return(nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIArticleRemoveFromUserInput{
					ArticleID: article.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ac")),
				},
			},
			want:    usecase.APIArticleRemoveFromUserOutput{},
			wantErr: false,
		},
		{
			name: "ユーザーに記事を追加できない",
			fields: fields{
				authCache: &cache.CacheMock[model.Auth]{
					T: t,
					Value: model.Auth{
						AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(time.Hour * 24 * 30),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				articleRPC: func(t *testing.T) rpc.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockArticle(ctrl)
					mock.EXPECT().RemoveFromUser(
						gomock.Any(),
						article.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ac")),
					).Return(fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.APIArticleRemoveFromUserInput{
					ArticleID: article.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ac")),
				},
			},
			want:    usecase.APIArticleRemoveFromUserOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			itr := interactor.NewAPIArticle(
				tt.fields.authCache,
				tt.fields.articleRPC(t),
			)
			got, err := itr.RemoveFromUser(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIArticle.RemoveFromUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIArticle.RemoveFromUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
