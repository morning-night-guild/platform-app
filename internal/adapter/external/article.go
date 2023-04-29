package external

import (
	"context"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
	articlev1 "github.com/morning-night-guild/platform-app/pkg/connect/article/v1"
	"github.com/morning-night-guild/platform-app/pkg/connect/article/v1/articlev1connect"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

type ArticleFactory interface {
	Article(string) (*Article, error)
}

var _ rpc.Article = (*Article)(nil)

type Article struct {
	connect articlev1connect.ArticleServiceClient
}

func NewArticle(
	connect articlev1connect.ArticleServiceClient,
) *Article {
	return &Article{
		connect: connect,
	}
}

func (ext *Article) Share(
	ctx context.Context,
	url article.URL,
	title article.Title,
	description article.Description,
	thumbnail article.Thumbnail,
) (model.Article, error) {
	req := NewRequest(ctx, &articlev1.ShareRequest{
		Url:         url.String(),
		Title:       title.String(),
		Description: description.String(),
		Thumbnail:   thumbnail.String(),
	})

	res, err := ext.connect.Share(ctx, req)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to share article", log.ErrorField(err))

		return model.Article{}, err
	}

	article := model.ReconstructArticle(
		uuid.MustParse(res.Msg.Article.ArticleId),
		res.Msg.Article.Url,
		res.Msg.Article.Title,
		res.Msg.Article.Description,
		res.Msg.Article.Thumbnail,
		[]string{},
	)

	return article, nil
}

func (ext *Article) List(
	ctx context.Context,
	index value.Index,
	size value.Size,
) ([]model.Article, error) {
	req := NewRequest(ctx, &articlev1.ListRequest{
		PageToken:   string(value.CreateNextTokenFromIndex(index)),
		MaxPageSize: uint32(size),
	})

	res, err := ext.connect.List(ctx, req)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to list articles", log.ErrorField(err))

		return nil, err
	}

	articles := make([]model.Article, len(res.Msg.Articles))

	for i, article := range res.Msg.Articles {
		articles[i] = model.ReconstructArticle(
			uuid.MustParse(article.ArticleId),
			article.Url,
			article.Title,
			article.Description,
			article.Thumbnail,
			article.Tags,
		)
	}

	return articles, nil
}

func (ext *Article) Delete(
	ctx context.Context,
	articleID article.ID,
) error {
	req := NewRequest(ctx, &articlev1.DeleteRequest{
		ArticleId: articleID.String(),
	})

	if _, err := ext.connect.Delete(ctx, req); err != nil {
		log.GetLogCtx(ctx).Sugar().Warnf("failed to delete articles. articleID=%s", articleID.String(), log.ErrorField(err))

		return err
	}

	return nil
}
