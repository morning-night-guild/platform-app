package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
)

var _ usecase.APIArticle = (*APIArticle)(nil)

// APIArticle.
type APIArticle struct {
	articleRPC rpc.Article
}

func NewAPIArticle(
	articleRPC rpc.Article,
) *APIArticle {
	return &APIArticle{
		articleRPC: articleRPC,
	}
}

func (aa *APIArticle) Share(
	ctx context.Context,
	input usecase.APIArticleShareInput,
) (usecase.APIArticleShareOutput, error) {
	article, err := aa.articleRPC.Share(ctx, input.URL, input.Title, input.Description, input.Thumbnail)
	if err != nil {
		return usecase.APIArticleShareOutput{}, err
	}

	return usecase.APIArticleShareOutput{
		Article: article,
	}, nil
}

func (aa *APIArticle) List(
	ctx context.Context,
	input usecase.APIArticleListInput,
) (usecase.APIArticleListOutput, error) {
	articles, err := aa.articleRPC.List(ctx, input.Index, input.Size)
	if err != nil {
		return usecase.APIArticleListOutput{}, err
	}

	return usecase.APIArticleListOutput{
		Articles: articles,
	}, nil
}
