package model_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
)

func TestNewArticle(t *testing.T) {
	t.Parallel()

	type args struct {
		id          article.ID
		url         article.URL
		title       article.Title
		description article.Description
		thumbnail   article.Thumbnail
		tags        article.TagList
	}

	tests := []struct {
		name    string
		args    args
		want    model.Article
		wantErr bool
	}{
		{
			name: "記事モデルが作成できる",
			args: args{
				id:          article.ID(uuid.MustParse("2f8e01fb-bf67-45cc-83b0-4cfa0548a9b2")),
				url:         article.URL("https://example.com"),
				title:       article.Title("タイトル"),
				description: article.Description("説明"),
				thumbnail:   article.Thumbnail("https://example.com"),
				tags:        article.TagList{},
			},
			want: model.Article{
				ArticleID:   article.ID(uuid.MustParse("2f8e01fb-bf67-45cc-83b0-4cfa0548a9b2")),
				URL:         article.URL("https://example.com"),
				Title:       article.Title("タイトル"),
				Description: article.Description("説明"),
				Thumbnail:   article.Thumbnail("https://example.com"),
				TagList:     article.TagList{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := model.NewArticle(
				tt.args.id,
				tt.args.url,
				tt.args.title,
				tt.args.description,
				tt.args.thumbnail,
				tt.args.tags,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArticle() = %v, want %v", got, tt.want)
			}
		})
	}
}
