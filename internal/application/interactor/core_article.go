package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
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
