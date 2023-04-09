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
					ID:          article.ID(uuid.New()),
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
			ca := interactor.NewCoreArticle(tt.fields.articleRepository(t))
			got, err := ca.Share(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CoreArticle.Share() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if _, err := uuid.Parse(got.Article.ID.String()); err != nil {
				t.Errorf("CoreArticle.Share() got Article.ID = %v, err %v", got.Article.ID, err)
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
							ID:          article.ID(id),
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
						ID:          article.ID(id),
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
			ca := interactor.NewCoreArticle(tt.fields.articleRepository(t))
			got, err := ca.List(tt.args.ctx, tt.args.input)
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
					).Return(model.Article{}, nil)
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
					ID: id,
				},
			},
			want:    usecase.CoreArticleDeleteOutput{},
			wantErr: false,
		},
		{
			name: "存在しない記事の削除がエラーにならない",
			fields: fields{
				articleRepository: func(t *testing.T) repository.Article {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockArticle(ctrl)
					mock.EXPECT().Find(
						gomock.Any(),
						id,
					).Return(model.Article{}, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CoreArticleDeleteInput{
					ID: id,
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
					).Return(model.Article{}, nil)
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
					ID: id,
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
			ca := interactor.NewCoreArticle(tt.fields.articleRepository(t))
			got, err := ca.Delete(tt.args.ctx, tt.args.input)
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
