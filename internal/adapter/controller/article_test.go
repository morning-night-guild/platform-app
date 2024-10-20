package controller_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/adapter/controller"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
	articlev1 "github.com/morning-night-guild/platform-app/pkg/connect/article/v1"
)

const aid = "01234567-0123-0123-0123-0123456789ab"

func TestArticleShare(t *testing.T) {
	t.Parallel()

	type fields struct {
		usecase func(*testing.T) usecase.CoreArticle
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
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					mock.EXPECT().Share(gomock.Any(), usecase.CoreArticleShareInput{
						URL:         article.URL("https://example.com"),
						Title:       article.Title("title"),
						Description: article.Description("description"),
						Thumbnail:   article.Thumbnail("https://example.com"),
					}).Return(usecase.CoreArticleShareOutput{
						Article: model.Article{
							ArticleID:   article.ID(uuid.MustParse(aid)),
							URL:         article.URL("https://example.com"),
							Title:       article.Title("title"),
							Description: article.Description("description"),
							Thumbnail:   article.Thumbnail("https://example.com"),
							TagList:     []article.Tag{},
						},
					}, nil)
					return mock
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
			want: connect.NewResponse(&articlev1.ShareResponse{
				Article: &articlev1.Article{
					ArticleId:   aid,
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
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					return mock
				},
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
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					return mock
				},
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
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					mock.EXPECT().Share(gomock.Any(), usecase.CoreArticleShareInput{
						URL:         article.URL("https://example.com"),
						Title:       article.Title("title"),
						Description: article.Description("description"),
						Thumbnail:   article.Thumbnail("https://example.com"),
					}).Return(usecase.CoreArticleShareOutput{
						Article: model.Article{},
					}, errors.NewValidationError("error"))
					return mock
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
			wantErr: controller.ErrInvalidArgument,
		},
		{
			name: "ユースケースでバリデーションエラー以外のエラーが発生した際、サーバーエラーになる",
			fields: fields{
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					mock.EXPECT().Share(gomock.Any(), usecase.CoreArticleShareInput{
						URL:         article.URL("https://example.com"),
						Title:       article.Title("title"),
						Description: article.Description("description"),
						Thumbnail:   article.Thumbnail("https://example.com"),
					}).Return(usecase.CoreArticleShareOutput{
						Article: model.Article{},
					}, fmt.Errorf("error"))
					return mock
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
			ctrl := controller.NewArticle(controller.New(), tt.fields.usecase(t))
			got, err := ctrl.Share(tt.args.ctx, tt.args.req)
			if err != nil && err != tt.wantErr {
				t.Errorf("Article.Share() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if tt.want == nil {
				return
			}
			if !reflect.DeepEqual(got.Msg.Article.ArticleId, tt.want.Msg.Article.ArticleId) {
				t.Errorf("Article.Share() Msg Article Id = %v, want %v", got.Msg.Article.ArticleId, tt.want.Msg.Article.ArticleId)
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

	toPointer := func(s string) *string {
		return &s
	}

	type fields struct {
		usecase func(*testing.T) usecase.CoreArticle
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
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					mock.EXPECT().List(gomock.Any(), usecase.CoreArticleListInput{
						Index: value.Index(0),
						Size:  value.Size(1),
					}).Return(usecase.CoreArticleListOutput{
						Articles: []model.Article{
							{
								ArticleID:   article.ID(id),
								URL:         article.URL("https://example.com"),
								Title:       article.Title("title"),
								Description: article.Description("description"),
								Thumbnail:   article.Thumbnail("https://example.com"),
								TagList:     []article.Tag{},
							},
						},
					}, nil)
					return mock
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
						ArticleId:   id.String(),
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
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					mock.EXPECT().List(gomock.Any(), usecase.CoreArticleListInput{
						Index: value.Index(0),
						Size:  value.Size(3),
					}).Return(usecase.CoreArticleListOutput{
						Articles: []model.Article{
							{
								ArticleID:   article.ID(id),
								URL:         article.URL("https://example.com"),
								Title:       article.Title("title"),
								Description: article.Description("description"),
								Thumbnail:   article.Thumbnail("https://example.com"),
								TagList:     []article.Tag{},
							},
						},
					}, nil)
					return mock
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
						ArticleId:   id.String(),
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
			name: "タイトルの部分一致検索で記事の一覧が取得できる",
			fields: fields{
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					mock.EXPECT().List(gomock.Any(), usecase.CoreArticleListInput{
						Index:  value.Index(0),
						Size:   value.Size(3),
						Filter: []value.Filter{value.NewFilter("title", "title")},
					}).Return(usecase.CoreArticleListOutput{
						Articles: []model.Article{
							{
								ArticleID:   article.ID(id),
								URL:         article.URL("https://example.com"),
								Title:       article.Title("title"),
								Description: article.Description("description"),
								Thumbnail:   article.Thumbnail("https://example.com"),
								TagList:     []article.Tag{},
							},
						},
					}, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.ListRequest]{
					Msg: &articlev1.ListRequest{
						PageToken:   "",
						MaxPageSize: 3,
						Title:       toPointer("title"),
					},
				},
			},
			want: connect.NewResponse(&articlev1.ListResponse{
				Articles: []*articlev1.Article{
					{
						ArticleId:   id.String(),
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
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					return mock
				},
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
			ctrl := controller.NewArticle(controller.New(), tt.fields.usecase(t))
			got, err := ctrl.List(tt.args.ctx, tt.args.req)
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

func TestArticleListByUser(t *testing.T) {
	t.Parallel()

	toPointer := func(s string) *string {
		return &s
	}

	type fields struct {
		usecase func(*testing.T) usecase.CoreArticle
	}

	type args struct {
		ctx context.Context
		req *connect.Request[articlev1.ListByUserRequest]
	}

	id := uuid.New()

	uid := uuid.New()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *connect.Response[articlev1.ListByUserResponse]
		wantErr bool
	}{
		{
			name: "記事の一覧が取得できる（ネクストトークンあり）",
			fields: fields{
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					mock.EXPECT().ListByUser(gomock.Any(), usecase.CoreArticleListByUserInput{
						UserID: user.ID(uid),
						Index:  value.Index(0),
						Size:   value.Size(1),
					}).Return(usecase.CoreArticleListByUserOutput{
						Articles: []model.Article{
							{
								ArticleID:   article.ID(id),
								URL:         article.URL("https://example.com"),
								Title:       article.Title("title"),
								Description: article.Description("description"),
								Thumbnail:   article.Thumbnail("https://example.com"),
								TagList:     []article.Tag{},
							},
						},
					}, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.ListByUserRequest]{
					Msg: &articlev1.ListByUserRequest{
						UserId:      uid.String(),
						PageToken:   "",
						MaxPageSize: 1,
					},
				},
			},
			want: connect.NewResponse(&articlev1.ListByUserResponse{
				Articles: []*articlev1.Article{
					{
						ArticleId:   id.String(),
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
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					mock.EXPECT().ListByUser(gomock.Any(), usecase.CoreArticleListByUserInput{
						UserID: user.ID(uid),
						Index:  value.Index(0),
						Size:   value.Size(3),
					}).Return(usecase.CoreArticleListByUserOutput{
						Articles: []model.Article{
							{
								ArticleID:   article.ID(id),
								URL:         article.URL("https://example.com"),
								Title:       article.Title("title"),
								Description: article.Description("description"),
								Thumbnail:   article.Thumbnail("https://example.com"),
								TagList:     []article.Tag{},
							},
						},
					}, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.ListByUserRequest]{
					Msg: &articlev1.ListByUserRequest{
						UserId:      uid.String(),
						PageToken:   "",
						MaxPageSize: 3,
					},
				},
			},
			want: connect.NewResponse(&articlev1.ListByUserResponse{
				Articles: []*articlev1.Article{
					{
						ArticleId:   id.String(),
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
			name: "タイトルの部分一致検索で記事の一覧が取得できる",
			fields: fields{
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					mock.EXPECT().ListByUser(gomock.Any(), usecase.CoreArticleListByUserInput{
						UserID: user.ID(uid),
						Index:  value.Index(0),
						Size:   value.Size(3),
						Filter: []value.Filter{value.NewFilter("title", "title")},
					}).Return(usecase.CoreArticleListByUserOutput{
						Articles: []model.Article{
							{
								ArticleID:   article.ID(id),
								URL:         article.URL("https://example.com"),
								Title:       article.Title("title"),
								Description: article.Description("description"),
								Thumbnail:   article.Thumbnail("https://example.com"),
								TagList:     []article.Tag{},
							},
						},
					}, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.ListByUserRequest]{
					Msg: &articlev1.ListByUserRequest{
						UserId:      uid.String(),
						PageToken:   "",
						MaxPageSize: 3,
						Title:       toPointer("title"),
					},
				},
			},
			want: connect.NewResponse(&articlev1.ListByUserResponse{
				Articles: []*articlev1.Article{
					{
						ArticleId:   id.String(),
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
			name: "usecaseでerrorが発生して記事の一覧が取得できない",
			fields: fields{
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					mock.EXPECT().ListByUser(gomock.Any(), usecase.CoreArticleListByUserInput{
						UserID: user.ID(uid),
						Index:  value.Index(0),
						Size:   value.Size(3),
					}).Return(usecase.CoreArticleListByUserOutput{}, fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.ListByUserRequest]{
					Msg: &articlev1.ListByUserRequest{
						UserId:      uid.String(),
						PageToken:   "",
						MaxPageSize: 3,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "不正なサイズを指定して記事の一覧が取得できない",
			fields: fields{
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.ListByUserRequest]{
					Msg: &articlev1.ListByUserRequest{
						UserId:      uid.String(),
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
			ctrl := controller.NewArticle(controller.New(), tt.fields.usecase(t))
			got, err := ctrl.ListByUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Article.ListByUser() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Article.ListByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleDelete(t *testing.T) {
	t.Parallel()

	type fields struct {
		usecase func(*testing.T) usecase.CoreArticle
	}

	type args struct {
		ctx context.Context
		req *connect.Request[articlev1.DeleteRequest]
	}

	id := uuid.New()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *connect.Response[articlev1.DeleteResponse]
		wantErr bool
	}{
		{
			name: "記事が削除できる",
			fields: fields{
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					mock.EXPECT().Delete(gomock.Any(), usecase.CoreArticleDeleteInput{
						ArticleID: article.ID(id),
					}).Return(usecase.CoreArticleDeleteOutput{}, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.DeleteRequest]{
					Msg: &articlev1.DeleteRequest{
						ArticleId: id.String(),
					},
				},
			},
			want:    connect.NewResponse(&articlev1.DeleteResponse{}),
			wantErr: false,
		},
		{
			name: "不正なIDを指定して記事が削除できない",
			fields: fields{
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.DeleteRequest]{
					Msg: &articlev1.DeleteRequest{
						ArticleId: "invalid",
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
			ctrl := controller.NewArticle(controller.New(), tt.fields.usecase(t))
			got, err := ctrl.Delete(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Article.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Article.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleAddToUser(t *testing.T) {
	t.Parallel()

	type fields struct {
		usecase func(*testing.T) usecase.CoreArticle
	}

	type args struct {
		ctx context.Context
		req *connect.Request[articlev1.AddToUserRequest]
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *connect.Response[articlev1.AddToUserResponse]
		wantErr bool
	}{
		{
			name: "記事をユーザーに追加できる",
			fields: fields{
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					mock.EXPECT().AddToUser(
						gomock.Any(),
						usecase.CoreArticleAddToUserInput{
							ArticleID: article.ID(uuid.MustParse(aid)),
							UserID:    user.ID(uuid.MustParse(uid)),
						},
					).Return(usecase.CoreArticleAddToUserOutput{}, nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.AddToUserRequest]{
					Msg: &articlev1.AddToUserRequest{
						ArticleId: aid,
						UserId:    uid,
					},
				},
			},
			want:    connect.NewResponse(&articlev1.AddToUserResponse{}),
			wantErr: false,
		},
		{
			name: "記事IDが不正な値でユーザーに記事を追加できない",
			fields: fields{
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.AddToUserRequest]{
					Msg: &articlev1.AddToUserRequest{
						ArticleId: "invalid",
						UserId:    uid,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ユーザーIDが不正な値でユーザーに記事を追加できない",
			fields: fields{
				usecase: func(t *testing.T) usecase.CoreArticle {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := usecase.NewMockCoreArticle(ctrl)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: &connect.Request[articlev1.AddToUserRequest]{
					Msg: &articlev1.AddToUserRequest{
						ArticleId: aid,
						UserId:    "invalid",
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
			ctrl := controller.NewArticle(controller.New(), tt.fields.usecase(t))
			got, err := ctrl.AddToUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Article.AddToUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Article.AddToUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
