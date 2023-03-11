package controller

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	articlev1 "github.com/morning-night-guild/platform-app/pkg/connect/article/v1"
)

// Article.
type Article struct {
	ctl   *Controller
	share port.CoreArticleShare
	list  port.CoreArticleList
}

// NewArticle 記事のコントローラを新規作成する関数.
func NewArticle(
	ctl *Controller,
	share port.CoreArticleShare,
	list port.CoreArticleList,
) *Article {
	return &Article{
		ctl:   ctl,
		share: share,
		list:  list,
	}
}

// Share 記事を共有するコントローラメソッド.
func (a *Article) Share(
	ctx context.Context,
	req *connect.Request[articlev1.ShareRequest],
) (*connect.Response[articlev1.ShareResponse], error) {
	url, err := article.NewURL(req.Msg.Url)
	if err != nil {
		return nil, a.ctl.HandleConnectError(ctx, err)
	}

	title, err := article.NewTitle(req.Msg.Title)
	if err != nil {
		return nil, a.ctl.HandleConnectError(ctx, err)
	}

	description, err := article.NewDescription(req.Msg.Description)
	if err != nil {
		return nil, a.ctl.HandleConnectError(ctx, err)
	}

	thumbnail, err := article.NewThumbnail(req.Msg.Thumbnail)
	if err != nil {
		return nil, a.ctl.HandleConnectError(ctx, err)
	}

	input := port.CoreArticleShareInput{
		URL:         url,
		Title:       title,
		Description: description,
		Thumbnail:   thumbnail,
	}

	output, err := a.share.Execute(ctx, input)
	if err != nil {
		return nil, a.ctl.HandleConnectError(ctx, err)
	}

	return connect.NewResponse(&articlev1.ShareResponse{
		Article: &articlev1.Article{
			Id:          output.Article.ID.String(),
			Title:       output.Article.Title.String(),
			Url:         output.Article.URL.String(),
			Description: output.Article.Description.String(),
			Thumbnail:   output.Article.Thumbnail.String(),
			Tags:        output.Article.TagList.StringSlice(),
		},
	}), nil
}

// List 記事を取得するコントローラメソッド.
func (a *Article) List(
	ctx context.Context,
	req *connect.Request[articlev1.ListRequest],
) (*connect.Response[articlev1.ListResponse], error) {
	token := value.NewNextToken(req.Msg.PageToken)

	index := token.ToIndex()

	size, err := value.NewSize(int(req.Msg.MaxPageSize))
	if err != nil {
		return nil, a.ctl.HandleConnectError(ctx, err)
	}

	input := port.CoreArticleListInput{
		Index: index,
		Size:  size,
	}

	output, err := a.list.Execute(ctx, input)
	if err != nil {
		return nil, a.ctl.HandleConnectError(ctx, err)
	}

	result := make([]*articlev1.Article, len(output.Articles))

	for i, article := range output.Articles {
		result[i] = &articlev1.Article{
			Id:          article.ID.String(),
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
