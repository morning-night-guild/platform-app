package gateway

import (
	"context"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	articlev1 "github.com/morning-night-guild/platform-app/pkg/connect/article/v1"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ repository.APIArticle = (*APIArticle)(nil)

type APIArticle struct {
	connect *Connect
}

func NewAPIArticle(connect *Connect) *APIArticle {
	return &APIArticle{
		connect: connect,
	}
}

func (ca *APIArticle) Save(
	ctx context.Context,
	url article.URL,
	title article.Title,
	description article.Description,
	thumbnail article.Thumbnail,
) (model.Article, error) {
	req := NewRequestWithTID(ctx, &articlev1.ShareRequest{
		Url:         url.String(),
		Title:       title.String(),
		Description: description.String(),
		Thumbnail:   thumbnail.String(),
	})

	res, err := ca.connect.Article.Share(ctx, req)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to share article", log.ErrorField(err))

		return model.Article{}, err
	}

	article := model.ReconstructArticle(
		uuid.MustParse(res.Msg.Article.Id),
		res.Msg.Article.Url,
		res.Msg.Article.Title,
		res.Msg.Article.Description,
		res.Msg.Article.Thumbnail,
		[]string{},
	)

	return article, nil
}

func (ca *APIArticle) FindAll(
	ctx context.Context,
	index repository.Index,
	size repository.Size,
) ([]model.Article, error) {
	req := NewRequestWithTID(ctx, &articlev1.ListRequest{
		PageToken:   string(repository.CreateNextTokenFromIndex(index)),
		MaxPageSize: uint32(size),
	})

	res, err := ca.connect.Article.List(ctx, req)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to list articles", log.ErrorField(err))

		return nil, err
	}

	articles := make([]model.Article, len(res.Msg.Articles))

	for i, article := range res.Msg.Articles {
		articles[i] = model.ReconstructArticle(
			uuid.MustParse(article.Id),
			article.Url,
			article.Title,
			article.Description,
			article.Thumbnail,
			article.Tags,
		)
	}

	return articles, nil
}
