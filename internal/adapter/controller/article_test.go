package controller_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/adapter/controller"
	"github.com/morning-night-guild/platform-app/internal/adapter/mock"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	dme "github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/usecase"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	articlev1 "github.com/morning-night-guild/platform-app/pkg/connect/article/v1"
)

func TestArticleShare(t *testing.T) {
	t.Parallel()

	type fields struct {
		share usecase.Usecase[port.ShareArticleInput, port.ShareArticleOutput]
		list  usecase.Usecase[port.ListArticleInput, port.ListArticleOutput]
	}

	type args struct {
		ctx context.Context
		req *connect.Request[articlev1.ShareRequest]
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *connect.Response[articlev1.ShareResponse]
		wantErr error
	}{
		{
			name: "記事の共有ができる",
			fields: fields{
				share: mock.ShareUsecase{
					T:   t,
					Err: nil,
				},
				list: mock.ListUsecase{},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.ShareRequest]{
					Msg: &articlev1.ShareRequest{
						Url:         "https://example.com",
						Title:       "title",
						Description: "description",
						Thumbnail:   "https://example.com",
					},
				},
			},
			want: connect.NewResponse(&articlev1.ShareResponse{
				Article: &articlev1.Article{
					Id:          mock.ID,
					Url:         "https://example.com",
					Title:       "title",
					Description: "description",
					Thumbnail:   "https://example.com",
				},
			}),
			wantErr: nil,
		},
		{
			name: "URLが不正の時、バッドリクエストエラーになる",
			fields: fields{
				share: mock.ShareUsecase{
					T:   t,
					Err: nil,
				},
				list: mock.ListUsecase{},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.ShareRequest]{
					Msg: &articlev1.ShareRequest{
						Url:         "http://example.com",
						Title:       "title",
						Description: "description",
						Thumbnail:   "https://example.com",
					},
				},
			},
			want:    nil,
			wantErr: controller.ErrInvalidArgument,
		},
		{
			name: "Thumbnailが不正の時、バッドリクエストエラーになる",
			fields: fields{
				share: mock.ShareUsecase{
					T:   t,
					Err: nil,
				},
				list: mock.ListUsecase{},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.ShareRequest]{
					Msg: &articlev1.ShareRequest{
						Url:         "https://example.com",
						Title:       "title",
						Description: "description",
						Thumbnail:   "http://example.com",
					},
				},
			},
			want:    nil,
			wantErr: controller.ErrInvalidArgument,
		},
		{
			name: "ユースケースでバリデーションエラーが発生した際、バッドリクエストエラーになる",
			fields: fields{
				share: mock.ShareUsecase{
					T:   t,
					Err: dme.NewValidationError("validation error"),
				},
				list: mock.ListUsecase{},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.ShareRequest]{
					Msg: &articlev1.ShareRequest{
						Url:         "https://example.com",
						Title:       "title",
						Description: "description",
						Thumbnail:   "https://example.com",
					},
				},
			},
			want:    nil,
			wantErr: controller.ErrInvalidArgument,
		},
		{
			name: "ユースケースでバリデーションエラー以外のエラーが発生した際、サーバーエラーになる",
			fields: fields{
				share: mock.ShareUsecase{
					T:   t,
					Err: errors.New("unknown error"),
				},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.ShareRequest]{
					Msg: &articlev1.ShareRequest{
						Url:         "https://example.com",
						Title:       "title",
						Description: "description",
						Thumbnail:   "https://example.com",
					},
				},
			},
			want:    nil,
			wantErr: controller.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := controller.NewArticle(controller.New(), tt.fields.share, tt.fields.list)
			got, err := a.Share(tt.args.ctx, tt.args.req)
			if err != nil && err != tt.wantErr {
				t.Errorf("Article.Share() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if tt.want == nil {
				return
			}
			if !reflect.DeepEqual(got.Msg.Article.Id, tt.want.Msg.Article.Id) {
				t.Errorf("Article.Share() Msg Article Id = %v, want %v", got.Msg.Article.Id, tt.want.Msg.Article.Id)
			}
			if !reflect.DeepEqual(got.Msg.Article.Url, tt.want.Msg.Article.Url) {
				t.Errorf("Article.Share() Msg Article Url = %v, want %v", got.Msg.Article.Url, tt.want.Msg.Article.Url)
			}
			if !reflect.DeepEqual(got.Msg.Article.Title, tt.want.Msg.Article.Title) {
				t.Errorf("Article.Share() Msg Article Title = %v, want %v", got.Msg.Article.Title, tt.want.Msg.Article.Title)
			}
			if !reflect.DeepEqual(got.Msg.Article.Description, tt.want.Msg.Article.Description) {
				t.Errorf("Article.Share() Msg Article Description = %v, want %v", got.Msg.Article.Description, tt.want.Msg.Article.Description)
			}
			if !reflect.DeepEqual(got.Msg.Article.Thumbnail, tt.want.Msg.Article.Thumbnail) {
				t.Errorf("Article.Share() Msg Article Thumbnail = %v, want %v", got.Msg.Article.Thumbnail, tt.want.Msg.Article.Thumbnail)
			}
		})
	}
}

func TestArticleList(t *testing.T) {
	t.Parallel()

	type fields struct {
		share usecase.Usecase[port.ShareArticleInput, port.ShareArticleOutput]
		list  usecase.Usecase[port.ListArticleInput, port.ListArticleOutput]
	}

	type args struct {
		ctx context.Context
		req *connect.Request[articlev1.ListRequest]
	}

	id := uuid.New()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *connect.Response[articlev1.ListResponse]
		wantErr bool
	}{
		{
			name: "記事の一覧が取得できる（ネクストトークンあり）",
			fields: fields{
				share: mock.ShareUsecase{},
				list: mock.ListUsecase{
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
				},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.ListRequest]{
					Msg: &articlev1.ListRequest{
						PageToken:   "",
						MaxPageSize: 1,
					},
				},
			},
			want: connect.NewResponse(&articlev1.ListResponse{
				Articles: []*articlev1.Article{
					{
						Id:          id.String(),
						Title:       "title",
						Url:         "https://example.com",
						Description: "description",
						Thumbnail:   "https://example.com",
						Tags:        []string{},
					},
				},
				NextPageToken: "MQ==",
			}),
			wantErr: false,
		},
		{
			name: "記事の一覧が取得できる（ネクストトークンなし）",
			fields: fields{
				share: mock.ShareUsecase{},
				list: mock.ListUsecase{
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
				},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.ListRequest]{
					Msg: &articlev1.ListRequest{
						PageToken:   "",
						MaxPageSize: 3,
					},
				},
			},
			want: connect.NewResponse(&articlev1.ListResponse{
				Articles: []*articlev1.Article{
					{
						Id:          id.String(),
						Title:       "title",
						Url:         "https://example.com",
						Description: "description",
						Thumbnail:   "https://example.com",
						Tags:        []string{},
					},
				},
				NextPageToken: "",
			}),
			wantErr: false,
		},
		{
			name: "不正なサイズを指定して記事の一覧が取得できない",
			fields: fields{
				share: mock.ShareUsecase{},
				list:  mock.ListUsecase{},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.ListRequest]{
					Msg: &articlev1.ListRequest{
						PageToken:   "",
						MaxPageSize: 1000,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := controller.NewArticle(controller.New(), tt.fields.share, tt.fields.list)
			got, err := a.List(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Article.List() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Article.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
