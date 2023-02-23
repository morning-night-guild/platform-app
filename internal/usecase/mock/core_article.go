package mock

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
)

var _ repository.CoreArticle = (*CoreArticle)(nil)

// Article 記事リポジトリのモック.
type CoreArticle struct {
	T        *testing.T
	Articles []model.Article
	Err      error
}

// Save 記事を保存するモックメソッド.
func (ca *CoreArticle) Save(ctx context.Context, article model.Article) error {
	ca.T.Helper()

	return ca.Err
}

// FindAll 記事を一覧取得するモックメソッド.
func (ca *CoreArticle) FindAll(
	ctx context.Context,
	index repository.Index,
	size repository.Size,
) ([]model.Article, error) {
	ca.T.Helper()

	return ca.Articles, ca.Err
}
