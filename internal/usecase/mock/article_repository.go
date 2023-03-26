package mock

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

var _ repository.Article = (*ArticleRepository)(nil)

type ArticleRepository struct {
	T        *testing.T
	Articles []model.Article
	Err      error
}

func (ar *ArticleRepository) Save(ctx context.Context, article model.Article) error {
	ar.T.Helper()

	return ar.Err
}

func (ar *ArticleRepository) FindAll(
	ctx context.Context,
	index value.Index,
	size value.Size,
) ([]model.Article, error) {
	ar.T.Helper()

	return ar.Articles, ar.Err
}
