package mock

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

type ShareUsecase struct {
	T   *testing.T
	Err error
}

const ID = "12345678-1234-1234-1234-1234567890ab"

func (su ShareUsecase) Execute(ctx context.Context, input port.ShareArticleInput) (port.ShareArticleOutput, error) {
	su.T.Helper()

	return port.ShareArticleOutput{
		Article: model.Article{
			ID:          article.ID(uuid.MustParse(ID)),
			URL:         input.URL,
			Title:       input.Title,
			Description: input.Description,
			Thumbnail:   input.Thumbnail,
		},
	}, su.Err
}

type ListUsecase struct {
	T        *testing.T
	Articles []model.Article
	Err      error
}

func (lu ListUsecase) Execute(ctx context.Context, input port.ListArticleInput) (port.ListArticleOutput, error) {
	lu.T.Helper()

	return port.ListArticleOutput{
		Articles: lu.Articles,
	}, lu.Err
}
