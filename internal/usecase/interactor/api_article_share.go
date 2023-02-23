package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

var _ port.APIArticleShare = (*APIArticleShare)(nil)

// APIArticleShare 記事共有のインタラクター.
type APIArticleShare struct {
	articleRepository repository.Article // 記事のリポジトリ
}

// NewAPIArticleShare 記事共有のインタラクターのファクトリ関数.
func NewAPIArticleShare(
	articleRepository repository.Article,
) *APIArticleShare {
	return &APIArticleShare{
		articleRepository: articleRepository,
	}
}

// Execute 記事共有のインタラクターを実行する.
func (aas *APIArticleShare) Execute(
	ctx context.Context,
	input port.APIArticleShareInput,
) (port.APIArticleShareOutput, error) {
	art := model.CreateArticle(input.URL, input.Title, input.Description, input.Thumbnail, []article.Tag{})

	if err := aas.articleRepository.Save(ctx, art); err != nil {
		return port.APIArticleShareOutput{}, err
	}

	return port.APIArticleShareOutput{
		Article: art,
	}, nil
}
