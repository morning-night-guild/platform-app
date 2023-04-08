package interactor_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/morning-night-guild/platform-app/internal/application/interactor"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

func TestAPIArticleShare(t *testing.T) {
	t.Parallel()

	type fields struct {
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
						ID:          id,
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
					ID:          id,
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
			aa := interactor.NewAPIArticle(tt.fields.articleRPC(t))
			got, err := aa.Share(tt.args.ctx, tt.args.input)
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
		articleRPC func(t *testing.T) rpc.Article
	}

	type args struct {
		ctx   context.Context
		input usecase.APIArticleListInput
	}

	articles := []model.Article{
		{
			ID:          article.GenerateID(),
			Title:       article.Title("title1"),
			URL:         article.URL("https://example.com/1"),
			Description: article.Description("description1"),
			Thumbnail:   article.Thumbnail("https://example.com/2"),
		},
		{
			ID:          article.GenerateID(),
			Title:       article.Title("title2"),
			URL:         article.URL("https://example.com/2"),
			Description: article.Description("description2"),
			Thumbnail:   article.Thumbnail("https://example.com/2"),
		},
	}

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
					Index: value.Index(0),
					Size:  value.Size(2),
				},
			},
			want: usecase.APIArticleListOutput{
				Articles: articles,
			},
			wantErr: false,
		},
		{
			name: "rpcでerrorが発生して記事リストが取得できない",
			fields: fields{
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
					Index: value.Index(0),
					Size:  value.Size(2),
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
			aa := interactor.NewAPIArticle(tt.fields.articleRPC(t))
			got, err := aa.List(tt.args.ctx, tt.args.input)
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
