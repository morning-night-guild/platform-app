package interactor_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/mock"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestAPIArticleShareExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		articleRPC rpc.Article
	}

	type args struct {
		ctx   context.Context
		input port.APIArticleShareInput
	}

	id := article.GenerateID()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    port.APIArticleShareOutput
		wantErr bool
	}{
		{
			name: "記事が共有できる",
			fields: fields{
				articleRPC: &mock.RPCArticle{
					T:   t,
					ID:  id,
					Err: nil,
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.APIArticleShareInput{
					URL:         article.URL("https://example.com"),
					Title:       article.Title("title"),
					Description: article.Description("description"),
					Thumbnail:   article.Thumbnail("https://example.com"),
				},
			},
			want: port.APIArticleShareOutput{
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
			name: "repositoryでerrorが発生して記事の共有が共有できない",
			fields: fields{
				articleRPC: &mock.RPCArticle{
					T:   t,
					ID:  id,
					Err: errors.New("error"),
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.APIArticleShareInput{
					URL:         article.URL("https://example.com"),
					Title:       article.Title("title"),
					Description: article.Description("description"),
					Thumbnail:   article.Thumbnail("https://example.com"),
				},
			},
			want: port.APIArticleShareOutput{
				Article: model.Article{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			aas := interactor.NewAPIArticleShare(tt.fields.articleRPC)
			got, err := aas.Execute(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIArticleShare.Execute() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIArticleShare.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
