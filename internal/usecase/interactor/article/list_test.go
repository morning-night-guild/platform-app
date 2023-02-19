package article_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	interactor "github.com/morning-night-guild/platform-app/internal/usecase/interactor/article"
	"github.com/morning-night-guild/platform-app/internal/usecase/mock"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestListInteractorExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		articleRepository repository.Article
	}

	type args struct {
		ctx   context.Context
		input port.ListArticleInput
	}

	id := uuid.New()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    port.ListArticleOutput
		wantErr bool
	}{
		{
			name: "記事一覧を取得できる",
			fields: fields{
				articleRepository: &mock.Article{
					T: t,
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
					Err: nil,
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.ListArticleInput{
					Index: repository.Index(0),
					Size:  repository.Size(1),
				},
			},
			want: port.ListArticleOutput{
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
			l := interactor.NewListInteractor(tt.fields.articleRepository)
			got, err := l.Execute(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListInteractor.Execute() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListInteractor.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
