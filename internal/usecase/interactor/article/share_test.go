package article_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	interactor "github.com/morning-night-guild/platform-app/internal/usecase/interactor/article"
	"github.com/morning-night-guild/platform-app/internal/usecase/mock"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestShareInteractorExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		articleRepository repository.Article
	}

	type args struct {
		ctx   context.Context
		input port.ShareArticleInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    port.ShareArticleOutput
		wantErr bool
	}{
		{
			name: "記事を共有できる",
			fields: fields{
				articleRepository: &mock.Article{
					T:   t,
					Err: nil,
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.ShareArticleInput{
					URL:         article.URL("https://example.com"),
					Title:       article.Title("title"),
					Description: article.Description("description"),
					Thumbnail:   article.Thumbnail("https://example.com"),
				},
			},
			want: port.ShareArticleOutput{
				Article: model.Article{
					ID:          article.ID{},
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
				articleRepository: &mock.Article{
					T:   t,
					Err: errors.New("article repository error"),
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.ShareArticleInput{
					URL:         article.URL("https://example.com"),
					Title:       article.Title("title"),
					Description: article.Description("description"),
					Thumbnail:   article.Thumbnail("https://example.com"),
				},
			},
			want:    port.ShareArticleOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := interactor.NewShareInteractor(tt.fields.articleRepository)
			got, err := s.Execute(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShareInteractor.Execute() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got.Article.URL, tt.want.Article.URL) {
				t.Errorf("ShareInteractor.Execute() got Article.URL = %v, want %v", got.Article.URL, tt.want.Article.URL)
			}
			if !reflect.DeepEqual(got.Article.Title, tt.want.Article.Title) {
				t.Errorf("ShareInteractor.Execute() got Article.Title = %v, want %v", got.Article.Title, tt.want.Article.Title)
			}
			if !reflect.DeepEqual(got.Article.Description, tt.want.Article.Description) {
				t.Errorf("ShareInteractor.Execute() got Article.Description = %v, want %v", got.Article.Description, tt.want.Article.Description)
			}
			if !reflect.DeepEqual(got.Article.Thumbnail, tt.want.Article.Thumbnail) {
				t.Errorf("ShareInteractor.Execute() got Article.Thumbnail = %v, want %v", got.Article.Thumbnail, tt.want.Article.Thumbnail)
			}
		})
	}
}
