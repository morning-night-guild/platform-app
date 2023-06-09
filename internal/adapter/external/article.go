package external

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
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
	connect  articlev1connect.ArticleServiceClient
	external *External
}

func NewArticle(
	connect articlev1connect.ArticleServiceClient,
) *Article {
	return &Article{
		connect:  connect,
		external: New(),
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

		return model.Article{}, ext.external.HandleError(ctx, err)
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
	filter ...value.Filter,
) ([]model.Article, error) {
	req := NewRequest(ctx, &articlev1.ListRequest{
		PageToken:   string(value.CreateNextTokenFromIndex(index)),
		MaxPageSize: uint32(size),
	})

	if len(filter) > 0 {
		for _, f := range filter {
			val := f.Value
			if f.Name == "title" {
				req.Msg.Title = &val
			}
		}
	}

	res, err := ext.connect.List(ctx, req)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to list articles", log.ErrorField(err))

		return nil, ext.external.HandleError(ctx, err)
	}

	return ext.toModels(res.Msg.Articles), nil
}

func (ext *Article) ListByUser(
	ctx context.Context,
	userID user.ID,
	index value.Index,
	size value.Size,
	filter ...value.Filter,
) ([]model.Article, error) {
	req := NewRequest(ctx, &articlev1.ListByUserRequest{
		UserId:      userID.String(),
		PageToken:   string(value.CreateNextTokenFromIndex(index)),
		MaxPageSize: uint32(size),
	})

	if len(filter) > 0 {
		for _, f := range filter {
			val := f.Value
			if f.Name == "title" {
				req.Msg.Title = &val
			}
		}
	}

	res, err := ext.connect.ListByUser(ctx, req)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to list articles", log.ErrorField(err))

		return nil, ext.external.HandleError(ctx, err)
	}

	return ext.toModels(res.Msg.Articles), nil
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

		return ext.external.HandleError(ctx, err)
	}

	return nil
}

func (ext *Article) AddToUser(
	ctx context.Context,
	articleID article.ID,
	userID user.ID,
) error {
	req := NewRequest(ctx, &articlev1.AddToUserRequest{
		ArticleId: articleID.String(),
		UserId:    userID.String(),
	})

	if _, err := ext.connect.AddToUser(ctx, req); err != nil {
		msg := fmt.Sprintf("failed to add article to user. articleID=%s, userID=%s", articleID.String(), userID.String())

		log.GetLogCtx(ctx).Warn(msg, log.ErrorField(err))

		return ext.external.HandleError(ctx, err)
	}

	return nil
}

func (ext *Article) RemoveFromUser(
	ctx context.Context,
	articleID article.ID,
	userID user.ID,
) error {
	req := NewRequest(ctx, &articlev1.RemoveFromUserRequest{
		ArticleId: articleID.String(),
		UserId:    userID.String(),
	})

	if _, err := ext.connect.RemoveFromUser(ctx, req); err != nil {
		msg := fmt.Sprintf("failed to remove article from user. articleID=%s, userID=%s", articleID.String(), userID.String())

		log.GetLogCtx(ctx).Warn(msg, log.ErrorField(err))

		return ext.external.HandleError(ctx, err)
	}

	return nil
}

func (ext *Article) toModels(
	articles []*articlev1.Article,
) []model.Article {
	res := make([]model.Article, len(articles))

	for i, article := range articles {
		res[i] = model.ReconstructArticle(
			uuid.MustParse(article.ArticleId),
			article.Url,
			article.Title,
			article.Description,
			article.Thumbnail,
			article.Tags,
		)
	}

	return res
}
