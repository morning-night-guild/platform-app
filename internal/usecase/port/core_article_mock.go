package port

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
)

const id = "12345678-1234-1234-1234-1234567890ab"

func MockID(t *testing.T) string {
	t.Helper()

	return id
}

var _ CoreArticleShare = (*CoreArticleShareMock)(nil)

type CoreArticleShareMock struct {
	T   *testing.T
	Err error
}

func (mock CoreArticleShareMock) Execute(
	ctx context.Context,
	input CoreArticleShareInput,
) (CoreArticleShareOutput, error) {
	mock.T.Helper()

	return CoreArticleShareOutput{
		Article: model.Article{
			ID:          article.ID(uuid.MustParse(MockID(mock.T))),
			URL:         input.URL,
			Title:       input.Title,
			Description: input.Description,
			Thumbnail:   input.Thumbnail,
		},
	}, mock.Err
}

var _ CoreArticleList = (*CoreArticleListMock)(nil)

type CoreArticleListMock struct {
	T        *testing.T
	Articles []model.Article
	Err      error
}

func (mock CoreArticleListMock) Execute(
	ctx context.Context,
	input CoreArticleListInput,
) (CoreArticleListOutput, error) {
	mock.T.Helper()

	return CoreArticleListOutput{
		Articles: mock.Articles,
	}, mock.Err
}
