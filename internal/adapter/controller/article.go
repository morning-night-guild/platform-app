package controller

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
	articlev1 "github.com/morning-night-guild/platform-app/pkg/connect/article/v1"
	"github.com/morning-night-guild/platform-app/pkg/connect/article/v1/articlev1connect"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ articlev1connect.ArticleServiceHandler = (*Article)(nil)

// Article.
type Article struct {
	controller *Controller
	usecase    usecase.CoreArticle
}

// NewArticle 記事のコントローラを新規作成する関数.
func NewArticle(
	controller *Controller,
	usecase usecase.CoreArticle,
) *Article {
	return &Article{
		controller: controller,
		usecase:    usecase,
	}
}

// Share 記事を共有するコントローラメソッド.
func (ctrl *Article) Share(
	ctx context.Context,
	req *connect.Request[articlev1.ShareRequest],
) (*connect.Response[articlev1.ShareResponse], error) {
	url, err := article.NewURL(req.Msg.Url)
	if err != nil {
		return nil, ctrl.controller.HandleConnectError(ctx, err)
	}

	title, err := article.NewTitle(req.Msg.Title)
	if err != nil {
		return nil, ctrl.controller.HandleConnectError(ctx, err)
	}

	description, err := article.NewDescription(req.Msg.Description)
	if err != nil {
		return nil, ctrl.controller.HandleConnectError(ctx, err)
	}

	thumbnail, err := article.NewThumbnail(req.Msg.Thumbnail)
	if err != nil {
		return nil, ctrl.controller.HandleConnectError(ctx, err)
	}

	input := usecase.CoreArticleShareInput{
		URL:         url,
		Title:       title,
		Description: description,
		Thumbnail:   thumbnail,
	}

	output, err := ctrl.usecase.Share(ctx, input)
	if err != nil {
		return nil, ctrl.controller.HandleConnectError(ctx, err)
	}

	return connect.NewResponse(&articlev1.ShareResponse{
		Article: &articlev1.Article{
			ArticleId:   output.Article.ArticleID.String(),
			Title:       output.Article.Title.String(),
			Url:         output.Article.URL.String(),
			Description: output.Article.Description.String(),
			Thumbnail:   output.Article.Thumbnail.String(),
			Tags:        output.Article.TagList.StringSlice(),
		},
	}), nil
}

// List 記事を取得するコントローラメソッド.
func (ctrl *Article) List(
	ctx context.Context,
	req *connect.Request[articlev1.ListRequest],
) (*connect.Response[articlev1.ListResponse], error) {
	token := value.NewNextToken(req.Msg.PageToken)

	index := token.ToIndex()

	size, err := value.NewSize(int(req.Msg.MaxPageSize))
	if err != nil {
		return nil, ctrl.controller.HandleConnectError(ctx, err)
	}

	input := usecase.CoreArticleListInput{
		Index: index,
		Size:  size,
	}

	if req.Msg.Title != nil {
		input.Filter = []value.Filter{
			value.NewFilter("title", *req.Msg.Title),
		}
	}

	output, err := ctrl.usecase.List(ctx, input)
	if err != nil {
		return nil, ctrl.controller.HandleConnectError(ctx, err)
	}

	result := make([]*articlev1.Article, len(output.Articles))

	for i, article := range output.Articles {
		result[i] = &articlev1.Article{
			ArticleId:   article.ArticleID.String(),
			Title:       article.Title.String(),
			Url:         article.URL.String(),
			Description: article.Description.String(),
			Thumbnail:   article.Thumbnail.String(),
			Tags:        article.TagList.StringSlice(),
		}
	}

	next := token.CreateNextToken(size).String()
	if len(output.Articles) < size.Int() {
		next = ""
	}

	return connect.NewResponse(&articlev1.ListResponse{
		Articles:      result,
		NextPageToken: next,
	}), nil
}

// Delete 記事を削除するコントローラメソッド.
func (ctrl *Article) Delete(
	ctx context.Context,
	req *connect.Request[articlev1.DeleteRequest],
) (*connect.Response[articlev1.DeleteResponse], error) {
	articleID, err := article.NewID(req.Msg.ArticleId)
	if err != nil {
		return nil, ctrl.controller.HandleConnectError(ctx, err)
	}

	input := usecase.CoreArticleDeleteInput{
		ArticleID: articleID,
	}

	if _, err = ctrl.usecase.Delete(ctx, input); err != nil {
		return nil, ctrl.controller.HandleConnectError(ctx, err)
	}

	return connect.NewResponse(&articlev1.DeleteResponse{}), nil
}

func (ctrl *Article) AddToUser(
	ctx context.Context,
	req *connect.Request[articlev1.AddToUserRequest],
) (*connect.Response[articlev1.AddToUserResponse], error) {
	articleID, err := article.NewID(req.Msg.ArticleId)
	if err != nil {
		return nil, ctrl.controller.HandleConnectError(ctx, err)
	}

	userID, err := user.NewID(req.Msg.UserId)
	if err != nil {
		return nil, ctrl.controller.HandleConnectError(ctx, err)
	}

	input := usecase.CoreArticleAddToUserInput{
		ArticleID: articleID,
		UserID:    userID,
	}

	if _, err := ctrl.usecase.AddToUser(ctx, input); err != nil {
		return nil, ctrl.controller.HandleConnectError(ctx, err)
	}

	return connect.NewResponse(&articlev1.AddToUserResponse{}), nil
}

func (ctrl *Article) RemoveFromUser(
	ctx context.Context,
	req *connect.Request[articlev1.RemoveFromUserRequest],
) (*connect.Response[articlev1.RemoveFromUserResponse], error) {
	log.GetLogCtx(ctx).Debug(fmt.Sprintf("%+v", req))

	return connect.NewResponse(&articlev1.RemoveFromUserResponse{}), nil
}
