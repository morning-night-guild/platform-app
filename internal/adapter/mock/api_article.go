package mock

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

var _ port.APIArticleShare = (*APIArticleShare)(nil)

type APIArticleShare struct {
	T   *testing.T
	Err error
}

func (cas APIArticleShare) Execute(
	ctx context.Context,
	input port.APIArticleShareInput,
) (port.APIArticleShareOutput, error) {
	cas.T.Helper()

	return port.APIArticleShareOutput{
		Article: model.Article{
			ID:          article.ID(uuid.MustParse(ID)),
			URL:         input.URL,
			Title:       input.Title,
			Description: input.Description,
			Thumbnail:   input.Thumbnail,
		},
	}, cas.Err
}

var _ port.APIArticleList = (*APIArticleList)(nil)

type APIArticleList struct {
	T        *testing.T
	Articles []model.Article
	Err      error
}

func (cau APIArticleList) Execute(
	ctx context.Context,
	input port.APIArticleListInput,
) (port.APIArticleListOutput, error) {
	cau.T.Helper()

	return port.APIArticleListOutput{
		Articles: cau.Articles,
	}, cau.Err
}
