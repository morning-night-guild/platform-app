package mock

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

var _ repository.Article = (*RepositoryArticle)(nil)

type RepositoryArticle struct {
	T        *testing.T
	Articles []model.Article
	Err      error
}

func (ra *RepositoryArticle) Save(ctx context.Context, article model.Article) error {
	ra.T.Helper()

	return ra.Err
}

func (ra *RepositoryArticle) FindAll(
	ctx context.Context,
	index value.Index,
	size value.Size,
) ([]model.Article, error) {
	ra.T.Helper()

	return ra.Articles, ra.Err
}
