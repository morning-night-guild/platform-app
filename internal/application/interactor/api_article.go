package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ usecase.APIArticle = (*APIArticle)(nil)

// APIArticle.
type APIArticle struct {
	authCache  cache.Cache[model.Auth]
	articleRPC rpc.Article
}

func NewAPIArticle(
	authCache cache.Cache[model.Auth],
	articleRPC rpc.Article,
) *APIArticle {
	return &APIArticle{
		authCache:  authCache,
		articleRPC: articleRPC,
	}
}

func (itr *APIArticle) Share(
	ctx context.Context,
	input usecase.APIArticleShareInput,
) (usecase.APIArticleShareOutput, error) {
	article, err := itr.articleRPC.Share(ctx, input.URL, input.Title, input.Description, input.Thumbnail)
	if err != nil {
		return usecase.APIArticleShareOutput{}, err
	}

	return usecase.APIArticleShareOutput{
		Article: article,
	}, nil
}

func (itr *APIArticle) List(
	ctx context.Context,
	input usecase.APIArticleListInput,
) (usecase.APIArticleListOutput, error) {
	auth, err := itr.authCache.Get(ctx, input.UserID.String())
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get auth cache", log.ErrorField(err))

		return usecase.APIArticleListOutput{}, errors.NewUnauthorizedError("failed to get auth cache", err)
	}

	if auth.IsExpired() {
		return usecase.APIArticleListOutput{}, errors.NewUnauthorizedError("auth token is expired")
	}

	switch input.Scope {
	case article.All:
		articles, err := itr.articleRPC.List(ctx, input.Index, input.Size, input.Filter...)
		if err != nil {
			return usecase.APIArticleListOutput{}, err
		}

		return usecase.APIArticleListOutput{
			Articles: articles,
		}, nil
	case article.Own:
		articles, err := itr.articleRPC.ListByUser(ctx, input.UserID, input.Index, input.Size, input.Filter...)
		if err != nil {
			return usecase.APIArticleListOutput{}, err
		}

		return usecase.APIArticleListOutput{
			Articles: articles,
		}, nil
	}

	return usecase.APIArticleListOutput{}, errors.NewValidationError("unknown scope")
}

func (itr *APIArticle) Delete(
	ctx context.Context,
	input usecase.APIArticleDeleteInput,
) (usecase.APIArticleDeleteOutput, error) {
	if err := itr.articleRPC.Delete(ctx, input.ArticleID); err != nil {
		return usecase.APIArticleDeleteOutput{}, err
	}

	return usecase.APIArticleDeleteOutput{}, nil
}

func (itr *APIArticle) AddToUser(
	ctx context.Context,
	input usecase.APIArticleAddToUserInput,
) (usecase.APIArticleAddToUserOutput, error) {
	auth, err := itr.authCache.Get(ctx, input.UserID.String())
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get auth cache", log.ErrorField(err))

		return usecase.APIArticleAddToUserOutput{}, errors.NewUnauthorizedError("failed to get auth cache", err)
	}

	if auth.IsExpired() {
		return usecase.APIArticleAddToUserOutput{}, errors.NewUnauthorizedError("auth token is expired")
	}

	if err := itr.articleRPC.AddToUser(ctx, input.ArticleID, input.UserID); err != nil {
		return usecase.APIArticleAddToUserOutput{}, err
	}

	return usecase.APIArticleAddToUserOutput{}, nil
}

func (itr *APIArticle) RemoveFromUser(
	ctx context.Context,
	input usecase.APIArticleRemoveFromUserInput,
) (usecase.APIArticleRemoveFromUserOutput, error) {
	auth, err := itr.authCache.Get(ctx, input.UserID.String())
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get auth cache", log.ErrorField(err))

		return usecase.APIArticleRemoveFromUserOutput{}, errors.NewUnauthorizedError("failed to get auth cache", err)
	}

	if auth.IsExpired() {
		return usecase.APIArticleRemoveFromUserOutput{}, errors.NewUnauthorizedError("auth token is expired")
	}

	if err := itr.articleRPC.RemoveFromUser(ctx, input.ArticleID, input.UserID); err != nil {
		return usecase.APIArticleRemoveFromUserOutput{}, err
	}

	return usecase.APIArticleRemoveFromUserOutput{}, nil
}
