package mock

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

var _ rpc.Article = (*ArticleRPC)(nil)

type ArticleRPC struct {
	T        *testing.T
	ID       article.ID
	Articles []model.Article
	Err      error
}

func (ar *ArticleRPC) Share(
	ctx context.Context,
	url article.URL,
	title article.Title,
	description article.Description,
	thumbnail article.Thumbnail,
) (model.Article, error) {
	ar.T.Helper()

	return model.Article{
		ID:          ar.ID,
		URL:         url,
		Title:       title,
		Description: description,
		Thumbnail:   thumbnail,
	}, ar.Err
}

func (ar *ArticleRPC) List(
	ctx context.Context,
	index value.Index,
	size value.Size,
) ([]model.Article, error) {
	ar.T.Helper()

	return ar.Articles, ar.Err
}
