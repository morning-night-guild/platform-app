package mock

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

type CoreArticleShare struct {
	T   *testing.T
	Err error
}

const ID = "12345678-1234-1234-1234-1234567890ab"

func (cas CoreArticleShare) Execute(ctx context.Context, input port.CoreArticleShareInput) (port.CoreArticleShareOutput, error) {
	cas.T.Helper()

	return port.CoreArticleShareOutput{
		Article: model.Article{
			ID:          article.ID(uuid.MustParse(ID)),
			URL:         input.URL,
			Title:       input.Title,
			Description: input.Description,
			Thumbnail:   input.Thumbnail,
		},
	}, cas.Err
}

type CoreArticleList struct {
	T        *testing.T
	Articles []model.Article
	Err      error
}

func (cau CoreArticleList) Execute(ctx context.Context, input port.CoreArticleListInput) (port.CoreArticleListOutput, error) {
	cau.T.Helper()

	return port.CoreArticleListOutput{
		Articles: cau.Articles,
	}, cau.Err
}
