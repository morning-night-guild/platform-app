package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
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
func (ca *CoreArticle) Share(
	ctx context.Context,
	input usecase.CoreArticleShareInput,
) (usecase.CoreArticleShareOutput, error) {
	item := model.CreateArticle(input.URL, input.Title, input.Description, input.Thumbnail, []article.Tag{})

	if err := ca.articleRepository.Save(ctx, item); err != nil {
		return usecase.CoreArticleShareOutput{}, err
	}

	return usecase.CoreArticleShareOutput{
		Article: item,
	}, nil
}

// List.
func (ca *CoreArticle) List(
	ctx context.Context,
	input usecase.CoreArticleListInput,
) (usecase.CoreArticleListOutput, error) {
	articles, err := ca.articleRepository.FindAll(ctx, input.Index, input.Size)
	if err != nil {
		return usecase.CoreArticleListOutput{}, err
	}

	return usecase.CoreArticleListOutput{
		Articles: articles,
	}, nil
}

// Delete.
func (ca *CoreArticle) Delete(
	ctx context.Context,
	input usecase.CoreArticleDeleteInput,
) (usecase.CoreArticleDeleteOutput, error) {
	article, err := ca.articleRepository.Find(ctx, input.ID)
	if err != nil {
		return usecase.CoreArticleDeleteOutput{}, err
	}

	if article.ID != input.ID {
		log.GetLogCtx(ctx).Sugar().Warnf("article not found. id=%s", input.ID)
		return usecase.CoreArticleDeleteOutput{}, nil
	}

	if err := ca.articleRepository.Delete(ctx, input.ID); err != nil {
		return usecase.CoreArticleDeleteOutput{}, err
	}

	return usecase.CoreArticleDeleteOutput{}, nil
}
