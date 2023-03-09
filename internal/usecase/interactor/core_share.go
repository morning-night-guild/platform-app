package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

var _ port.CoreArticleShare = (*CoreArticleShare)(nil)

// CoreArticleShare 記事共有のインタラクター.
type CoreArticleShare struct {
	articleRepository repository.Article // 記事のリポジトリ
}

// NewCoreArticleShare 記事共有のインタラクターのファクトリ関数.
func NewCoreArticleShare(
	articleRepository repository.Article,
) *CoreArticleShare {
	return &CoreArticleShare{
		articleRepository: articleRepository,
	}
}

// Execute 記事共有のインタラクターを実行する.
func (cas *CoreArticleShare) Execute(
	ctx context.Context,
	input port.CoreArticleShareInput,
) (port.CoreArticleShareOutput, error) {
	art := model.CreateArticle(input.URL, input.Title, input.Description, input.Thumbnail, []article.Tag{})

	if err := cas.articleRepository.Save(ctx, art); err != nil {
		return port.CoreArticleShareOutput{}, err
	}

	return port.CoreArticleShareOutput{
		Article: art,
	}, nil
}
