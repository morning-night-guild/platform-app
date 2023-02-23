package gateway

import (
	"context"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	articlev1 "github.com/morning-night-guild/platform-app/pkg/connect/article/v1"
	healthv1 "github.com/morning-night-guild/platform-app/pkg/connect/health/v1"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ repository.Article = (*CoreArticle)(nil)

type CoreArticle struct {
	connect *Connect
}

func NewCoreArticle(connect *Connect) *CoreArticle {
	return &CoreArticle{
		connect: connect,
	}
}

func (ca *CoreArticle) Save(ctx context.Context, item model.Article) error {
	req := NewRequestWithTID(ctx, &articlev1.ShareRequest{
		Url:         item.URL.String(),
		Title:       item.Title.String(),
		Description: item.Description.String(),
		Thumbnail:   item.Thumbnail.String(),
	})

	if _, err := ca.connect.Article.Share(ctx, req); err != nil {
		log.GetLogCtx(ctx).Warn("failed to share article", log.ErrorField(err))

		return err
	}

	return nil
}

func (ca *CoreArticle) FindAll(ctx context.Context, index repository.Index, size repository.Size) ([]model.Article, error) {
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

var _ repository.Health = (*CoreHealth)(nil)

type CoreHealth struct {
	connect *Connect
}

func NewCoreHealth(connect *Connect) *CoreHealth {
	return &CoreHealth{
		connect: connect,
	}
}

func (ch *CoreHealth) Check(ctx context.Context) error {
	req := NewRequestWithTID(ctx, &healthv1.CheckRequest{})

	if _, err := ch.connect.Health.Check(ctx, req); err != nil {
		log.GetLogCtx(ctx).Warn("failed to check health core", log.ErrorField(err))

		return err
	}

	return nil
}
