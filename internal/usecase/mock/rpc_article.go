package mock

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

var _ rpc.Article = (*RPCArticle)(nil)

type RPCArticle struct {
	T        *testing.T
	ID       article.ID
	Articles []model.Article
	Err      error
}

func (ra *RPCArticle) Share(
	ctx context.Context,
	url article.URL,
	title article.Title,
	description article.Description,
	thumbnail article.Thumbnail,
) (model.Article, error) {
	ra.T.Helper()

	return model.Article{
		ID:          ra.ID,
		URL:         url,
		Title:       title,
		Description: description,
		Thumbnail:   thumbnail,
	}, ra.Err
}

func (ra *RPCArticle) List(
	ctx context.Context,
	index value.Index,
	size value.Size,
) ([]model.Article, error) {
	ra.T.Helper()

	return ra.Articles, ra.Err
}
