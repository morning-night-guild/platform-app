package interactor_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/mock"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestAPIArticleListExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		articleRPC rpc.Article
	}

	type args struct {
		ctx   context.Context
		input port.APIArticleListInput
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
		want    port.APIArticleListOutput
		wantErr bool
	}{
		{
			name: "記事リストが取得できる",
			fields: fields{
				articleRPC: &mock.RPCArticle{
					T:        t,
					Articles: articles,
					Err:      nil,
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.APIArticleListInput{
					Index: value.Index(0),
					Size:  value.Size(2),
				},
			},
			want: port.APIArticleListOutput{
				Articles: articles,
			},
			wantErr: false,
		},
		{
			name: "rpcでerrorが発生して記事リストが取得できない",
			fields: fields{
				articleRPC: &mock.RPCArticle{
					T:        t,
					Articles: nil,
					Err:      errors.New("error"),
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.APIArticleListInput{
					Index: value.Index(0),
					Size:  value.Size(2),
				},
			},
			want: port.APIArticleListOutput{
				Articles: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			aal := interactor.NewAPIArticleList(tt.fields.articleRPC)
			got, err := aal.Execute(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIArticleList.Execute() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIArticleList.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
