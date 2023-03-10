package interactor_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/mock"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestCoreArticleListExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		articleRepository repository.Article
	}

	type args struct {
		ctx   context.Context
		input port.CoreArticleListInput
	}

	id := uuid.New()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    port.CoreArticleListOutput
		wantErr bool
	}{
		{
			name: "記事一覧を取得できる",
			fields: fields{
				articleRepository: &mock.RepositoryArticle{
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
				input: port.CoreArticleListInput{
					Index: value.Index(0),
					Size:  value.Size(1),
				},
			},
			want: port.CoreArticleListOutput{
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
			cal := interactor.NewCoreArticleList(tt.fields.articleRepository)
			got, err := cal.Execute(tt.args.ctx, tt.args.input)
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
