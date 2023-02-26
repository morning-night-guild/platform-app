package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

// CoreArticleList 記事一覧のインタラクター.
type CoreArticleList struct {
	articleRepository repository.CoreArticle // 記事のリポジトリ
}

// NewCoreArticleList 記事一覧のインタラクターのファクトリ関数.
func NewCoreArticleList(
	articleRepository repository.CoreArticle,
) port.CoreArticleList {
	return &CoreArticleList{
		articleRepository: articleRepository,
	}
}

// Execute 記事一覧のインタラクターを実行する.
func (l *CoreArticleList) Execute(
	ctx context.Context,
	input port.CoreArticleListInput,
) (port.CoreArticleListOutput, error) {
	articles, err := l.articleRepository.FindAll(ctx, input.Index, input.Size)
	if err != nil {
		return port.CoreArticleListOutput{}, err
	}

	return port.CoreArticleListOutput{
		Articles: articles,
	}, nil
}
