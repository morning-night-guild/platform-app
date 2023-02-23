package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

var _ port.APIArticleList = (*APIArticleList)(nil)

// APIArticleList 記事一覧のインタラクター.
type APIArticleList struct {
	articleRepository repository.Article // 記事のリポジトリ
}

// NewAPIArticleList 記事一覧のインタラクターのファクトリ関数.
func NewAPIArticleList(
	articleRepository repository.Article,
) *APIArticleList {
	return &APIArticleList{
		articleRepository: articleRepository,
	}
}

// Execute 記事一覧のインタラクターを実行する.
func (aal *APIArticleList) Execute(
	ctx context.Context,
	input port.APIArticleListInput,
) (port.APIArticleListOutput, error) {
	articles, err := aal.articleRepository.FindAll(ctx, input.Index, input.Size)
	if err != nil {
		return port.APIArticleListOutput{}, err
	}

	return port.APIArticleListOutput{
		Articles: articles,
	}, nil
}
