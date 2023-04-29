package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ usecase.CoreArticle = (*CoreArticle)(nil)

// CoreArticle.
type CoreArticle struct {
	articleRepository repository.Article
}

func NewCoreArticle(
	articleRepository repository.Article,
) *CoreArticle {
	return &CoreArticle{
		articleRepository: articleRepository,
	}
}

// Share.
func (itr *CoreArticle) Share(
	ctx context.Context,
	input usecase.CoreArticleShareInput,
) (usecase.CoreArticleShareOutput, error) {
	item := model.CreateArticle(input.URL, input.Title, input.Description, input.Thumbnail, []article.Tag{})

	if err := itr.articleRepository.Save(ctx, item); err != nil {
		return usecase.CoreArticleShareOutput{}, err
	}

	return usecase.CoreArticleShareOutput{
		Article: item,
	}, nil
}

// List.
func (itr *CoreArticle) List(
	ctx context.Context,
	input usecase.CoreArticleListInput,
) (usecase.CoreArticleListOutput, error) {
	articles, err := itr.articleRepository.FindAll(ctx, input.Index, input.Size)
	if err != nil {
		return usecase.CoreArticleListOutput{}, err
	}

	return usecase.CoreArticleListOutput{
		Articles: articles,
	}, nil
}

// Delete.
func (itr *CoreArticle) Delete(
	ctx context.Context,
	input usecase.CoreArticleDeleteInput,
) (usecase.CoreArticleDeleteOutput, error) {
	if _, err := itr.articleRepository.Find(ctx, input.ArticleID); err != nil {
		if errors.AsNotFoundError(err) {
			log.GetLogCtx(ctx).Sugar().Warnf("article not found. id=%s", input.ArticleID)

			return usecase.CoreArticleDeleteOutput{}, nil
		}

		return usecase.CoreArticleDeleteOutput{}, err
	}

	if err := itr.articleRepository.Delete(ctx, input.ArticleID); err != nil {
		return usecase.CoreArticleDeleteOutput{}, err
	}

	return usecase.CoreArticleDeleteOutput{}, nil
}
