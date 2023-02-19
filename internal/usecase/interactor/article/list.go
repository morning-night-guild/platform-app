package article

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

var _ port.ListArticle = (*ListInteractor)(nil)

// ListInteractor 記事一覧のインタラクター.
type ListInteractor struct {
	articleRepository repository.Article // 記事のリポジトリ
}

// NewListInteractor 記事一覧のインタラクターのファクトリ関数.
func NewListInteractor(
	articleRepository repository.Article,
) *ListInteractor {
	return &ListInteractor{
		articleRepository: articleRepository,
	}
}

// Execute 記事一覧のインタラクターを実行する.
func (l *ListInteractor) Execute(ctx context.Context, input port.ListArticleInput) (port.ListArticleOutput, error) {
	articles, err := l.articleRepository.FindAll(ctx, input.Index, input.Size)
	if err != nil {
		return port.ListArticleOutput{}, err
	}

	return port.ListArticleOutput{
		Articles: articles,
	}, nil
}
