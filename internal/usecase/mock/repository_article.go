package mock

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

var _ repository.Article = (*RepositoryArticle)(nil)

// Article 記事リポジトリのモック.
type RepositoryArticle struct {
	T        *testing.T
	Articles []model.Article
	Err      error
}

// Save 記事を保存するモックメソッド.
func (ra *RepositoryArticle) Save(ctx context.Context, article model.Article) error {
	ra.T.Helper()

	return ra.Err
}

// FindAll 記事を一覧取得するモックメソッド.
func (ra *RepositoryArticle) FindAll(
	ctx context.Context,
	index value.Index,
	size value.Size,
) ([]model.Article, error) {
	ra.T.Helper()

	return ra.Articles, ra.Err
}
