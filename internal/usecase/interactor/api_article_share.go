package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

// APIArticleShare 記事共有のインタラクター.
type APIArticleShare struct {
	articleRepository repository.APIArticle // 記事のリポジトリ
}

// NewAPIArticleShare 記事共有のインタラクターのファクトリ関数.
func NewAPIArticleShare(
	articleRepository repository.APIArticle,
) port.APIArticleShare {
	return &APIArticleShare{
		articleRepository: articleRepository,
	}
}

// Execute 記事共有のインタラクターを実行する.
func (aas *APIArticleShare) Execute(
	ctx context.Context,
	input port.APIArticleShareInput,
) (port.APIArticleShareOutput, error) {
	res, err := aas.articleRepository.Save(
		ctx,
		input.URL,
		input.Title,
		input.Description,
		input.Thumbnail,
	)
	if err != nil {
		return port.APIArticleShareOutput{}, err
	}

	return port.APIArticleShareOutput{
		Article: res,
	}, nil
}
