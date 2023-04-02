package port

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
)

var _ APIArticleShare = (*APIArticleShareMock)(nil)

type APIArticleShareMock struct {
	T   *testing.T
	Err error
}

func (mock APIArticleShareMock) Execute(
	ctx context.Context,
	input APIArticleShareInput,
) (APIArticleShareOutput, error) {
	mock.T.Helper()

	return APIArticleShareOutput{
		Article: model.Article{
			ID:          article.GenerateID(),
			URL:         input.URL,
			Title:       input.Title,
			Description: input.Description,
			Thumbnail:   input.Thumbnail,
		},
	}, mock.Err
}

var _ APIArticleList = (*APIArticleListMock)(nil)

type APIArticleListMock struct {
	T        *testing.T
	Articles []model.Article
	Err      error
}

func (aal APIArticleListMock) Execute(
	ctx context.Context,
	input APIArticleListInput,
) (APIArticleListOutput, error) {
	aal.T.Helper()

	return APIArticleListOutput{
		Articles: aal.Articles,
	}, aal.Err
}
