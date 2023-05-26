package interactor_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/application/interactor"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

func TestCoreArticleShare(t *testing.T) {
	t.Parallel()

	type fields struct {
		articleRepository func(t *testing.T) repository.Article
	}

	type args struct {
		ctx   context.Context
		input usecase.CoreArticleShareInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.CoreArticleShareOutput
		wantErr bool
	}{
		{
			name: "記事を共有できる",
			fields: fields{
				articleRepository: func(t *testing.T) repository.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockArticle(ctrl)
					mock.EXPECT().Save(
						gomock.Any(),
						gomock.Any(),
					).Return(nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CoreArticleShareInput{
					URL:         article.URL("https://example.com"),
					Title:       article.Title("title"),
					Description: article.Description("description"),
					Thumbnail:   article.Thumbnail("https://example.com"),
				},
			},
			want: usecase.CoreArticleShareOutput{
				Article: model.Article{
					ArticleID:   article.ID(uuid.New()),
					URL:         article.URL("https://example.com"),
					Title:       article.Title("title"),
					Description: article.Description("description"),
					Thumbnail:   article.Thumbnail("https://example.com"),
				},
			},
			wantErr: false,
		},
		{
			name: "記事Repositoryのerrorを握りつぶさない",
			fields: fields{
				articleRepository: func(t *testing.T) repository.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockArticle(ctrl)
					mock.EXPECT().Save(
						gomock.Any(),
						gomock.Any(),
					).Return(fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CoreArticleShareInput{
					URL:         article.URL("https://example.com"),
					Title:       article.Title("title"),
					Description: article.Description("description"),
					Thumbnail:   article.Thumbnail("https://example.com"),
				},
			},
			want:    usecase.CoreArticleShareOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			itr := interactor.NewCoreArticle(
				tt.fields.articleRepository(t),
				nil,
			)
			got, err := itr.Share(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CoreArticle.Share() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if _, err := uuid.Parse(got.Article.ArticleID.String()); err != nil {
				t.Errorf("CoreArticle.Share() got Article.ID = %v, err %v", got.Article.ArticleID, err)
			}
			if !reflect.DeepEqual(got.Article.URL, tt.want.Article.URL) {
				t.Errorf("CoreArticle.Share() got Article.URL = %v, want %v", got.Article.URL, tt.want.Article.URL)
			}
			if !reflect.DeepEqual(got.Article.Title, tt.want.Article.Title) {
				t.Errorf("CoreArticle.Share() got Article.Title = %v, want %v", got.Article.Title, tt.want.Article.Title)
			}
			if !reflect.DeepEqual(got.Article.Description, tt.want.Article.Description) {
				t.Errorf("CoreArticle.Share() got Article.Description = %v, want %v", got.Article.Description, tt.want.Article.Description)
			}
			if !reflect.DeepEqual(got.Article.Thumbnail, tt.want.Article.Thumbnail) {
				t.Errorf("CoreArticle.Share() got Article.Thumbnail = %v, want %v", got.Article.Thumbnail, tt.want.Article.Thumbnail)
			}
		})
	}
}

func TestCoreArticleList(t *testing.T) {
	t.Parallel()

	type fields struct {
		articleRepository func(t *testing.T) repository.Article
	}

	type args struct {
		ctx   context.Context
		input usecase.CoreArticleListInput
	}

	id := uuid.New()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.CoreArticleListOutput
		wantErr bool
	}{
		{
			name: "記事一覧を取得できる",
			fields: fields{
				articleRepository: func(t *testing.T) repository.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockArticle(ctrl)
					mock.EXPECT().FindAll(
						gomock.Any(),
						value.Index(0),
						value.Size(1),
					).Return([]model.Article{
						{
							ArticleID:   article.ID(id),
							Title:       article.Title("title"),
							URL:         article.URL("https://example.com"),
							Description: article.Description("description"),
							Thumbnail:   article.Thumbnail("https://example.com"),
							TagList:     article.TagList{},
						},
					}, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CoreArticleListInput{
					Index: value.Index(0),
					Size:  value.Size(1),
				},
			},
			want: usecase.CoreArticleListOutput{
				Articles: []model.Article{
					{
						ArticleID:   article.ID(id),
						Title:       article.Title("title"),
						URL:         article.URL("https://example.com"),
						Description: article.Description("description"),
						Thumbnail:   article.Thumbnail("https://example.com"),
						TagList:     article.TagList{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "タイトル部分一致で記事一覧を取得できる",
			fields: fields{
				articleRepository: func(t *testing.T) repository.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockArticle(ctrl)
					mock.EXPECT().FindAll(
						gomock.Any(),
						value.Index(0),
						value.Size(1),
						[]value.Filter{value.NewFilter("title", "title")},
					).Return([]model.Article{
						{
							ArticleID:   article.ID(id),
							Title:       article.Title("title"),
							URL:         article.URL("https://example.com"),
							Description: article.Description("description"),
							Thumbnail:   article.Thumbnail("https://example.com"),
							TagList:     article.TagList{},
						},
					}, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CoreArticleListInput{
					Index:  value.Index(0),
					Size:   value.Size(1),
					Filter: []value.Filter{value.NewFilter("title", "title")},
				},
			},
			want: usecase.CoreArticleListOutput{
				Articles: []model.Article{
					{
						ArticleID:   article.ID(id),
						Title:       article.Title("title"),
						URL:         article.URL("https://example.com"),
						Description: article.Description("description"),
						Thumbnail:   article.Thumbnail("https://example.com"),
						TagList:     article.TagList{},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			itr := interactor.NewCoreArticle(
				tt.fields.articleRepository(t),
				nil,
			)
			got, err := itr.List(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CoreArticle.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CoreArticle.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoreArticleDelete(t *testing.T) {
	t.Parallel()

	type fields struct {
		articleRepository func(*testing.T) repository.Article
	}

	type args struct {
		ctx   context.Context
		input usecase.CoreArticleDeleteInput
	}

	id := article.ID(uuid.New())

	article := model.Article{
		ArticleID:   id,
		Title:       article.Title("title"),
		URL:         article.URL("https://example.com"),
		Description: article.Description("description"),
		Thumbnail:   article.Thumbnail("https://example.com"),
		TagList:     article.TagList{},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.CoreArticleDeleteOutput
		wantErr bool
	}{
		{
			name: "記事を削除できる",
			fields: fields{
				articleRepository: func(t *testing.T) repository.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockArticle(ctrl)
					mock.EXPECT().Find(
						gomock.Any(),
						id,
					).Return(article, nil)
					mock.EXPECT().Delete(
						gomock.Any(),
						id,
					).Return(nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CoreArticleDeleteInput{
					ArticleID: id,
				},
			},
			want:    usecase.CoreArticleDeleteOutput{},
			wantErr: false,
		},
		{
			name: "存在しない記事を削除してもエラーにならない",
			fields: fields{
				articleRepository: func(t *testing.T) repository.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockArticle(ctrl)
					mock.EXPECT().Find(
						gomock.Any(),
						id,
					).Return(model.Article{}, errors.NewNotFoundError("article not found"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CoreArticleDeleteInput{
					ArticleID: id,
				},
			},
			want:    usecase.CoreArticleDeleteOutput{},
			wantErr: false,
		},
		{
			name: "記事削除のerrorを握りつぶさない",
			fields: fields{
				articleRepository: func(t *testing.T) repository.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockArticle(ctrl)
					mock.EXPECT().Find(
						gomock.Any(),
						id,
					).Return(article, nil)
					mock.EXPECT().Delete(
						gomock.Any(),
						id,
					).Return(fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CoreArticleDeleteInput{
					ArticleID: id,
				},
			},
			want:    usecase.CoreArticleDeleteOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			itr := interactor.NewCoreArticle(
				tt.fields.articleRepository(t),
				nil,
			)
			got, err := itr.Delete(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CoreArticle.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CoreArticle.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoreArticleAddToUser(t *testing.T) {
	t.Parallel()

	type fields struct {
		articleRepository func(*testing.T) repository.Article
		userRepository    func(*testing.T) repository.User
	}

	type args struct {
		ctx   context.Context
		input usecase.CoreArticleAddToUserInput
	}

	uid := user.ID(uuid.New())

	aid := article.ID(uuid.New())

	article := model.Article{
		ArticleID:   aid,
		Title:       article.Title("title"),
		URL:         article.URL("https://example.com"),
		Description: article.Description("description"),
		Thumbnail:   article.Thumbnail("https://example.com"),
		TagList:     article.TagList{},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.CoreArticleAddToUserOutput
		wantErr bool
	}{
		{
			name: "記事をユーザーに追加できる",
			fields: fields{
				articleRepository: func(t *testing.T) repository.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockArticle(ctrl)
					mock.EXPECT().Find(
						gomock.Any(),
						aid,
					).Return(article, nil)
					mock.EXPECT().AddToUser(
						gomock.Any(),
						aid,
						uid,
					).Return(nil)
					return mock
				},
				userRepository: func(t *testing.T) repository.User {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockUser(ctrl)
					mock.EXPECT().Find(
						gomock.Any(),
						uid,
					).Return(model.User{
						UserID: uid,
					}, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CoreArticleAddToUserInput{
					ArticleID: aid,
					UserID:    uid,
				},
			},
			want:    usecase.CoreArticleAddToUserOutput{},
			wantErr: false,
		},
		{
			name: "指定したユーザーが存在しない",
			fields: fields{
				articleRepository: func(t *testing.T) repository.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockArticle(ctrl)
					return mock
				},
				userRepository: func(t *testing.T) repository.User {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockUser(ctrl)
					mock.EXPECT().Find(
						gomock.Any(),
						uid,
					).Return(model.User{}, fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CoreArticleAddToUserInput{
					ArticleID: aid,
					UserID:    uid,
				},
			},
			want:    usecase.CoreArticleAddToUserOutput{},
			wantErr: true,
		},
		{
			name: "指定した記事が存在しない",
			fields: fields{
				articleRepository: func(t *testing.T) repository.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockArticle(ctrl)
					mock.EXPECT().Find(
						gomock.Any(),
						aid,
					).Return(model.Article{}, fmt.Errorf("error"))
					return mock
				},
				userRepository: func(t *testing.T) repository.User {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockUser(ctrl)
					mock.EXPECT().Find(
						gomock.Any(),
						uid,
					).Return(model.User{
						UserID: uid,
					}, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CoreArticleAddToUserInput{
					ArticleID: aid,
					UserID:    uid,
				},
			},
			want:    usecase.CoreArticleAddToUserOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			itr := interactor.NewCoreArticle(
				tt.fields.articleRepository(t),
				tt.fields.userRepository(t),
			)
			got, err := itr.AddToUser(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CoreArticle.AddToUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CoreArticle.AddToUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
