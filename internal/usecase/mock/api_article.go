package mock

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
)

var _ repository.APIArticle = (*APIArticle)(nil)

// Article 記事リポジトリのモック.
type APIArticle struct {
	T        *testing.T
	ID       article.ID
	Articles []model.Article
	Err      error
}

// Save 記事を保存するモックメソッド.
func (ca *APIArticle) Save(
	ctx context.Context,
	url article.URL,
	title article.Title,
	description article.Description,
	thumbnail article.Thumbnail,
) (model.Article, error) {
	ca.T.Helper()

	return model.Article{
		ID:          ca.ID,
		URL:         url,
		Title:       title,
		Description: description,
		Thumbnail:   thumbnail,
	}, ca.Err
}

// FindAll 記事を一覧取得するモックメソッド.
func (ca *APIArticle) FindAll(
	ctx context.Context,
	index repository.Index,
	size repository.Size,
) ([]model.Article, error) {
	ca.T.Helper()

	return ca.Articles, ca.Err
}
