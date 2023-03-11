package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

var _ port.APIArticleList = (*APIArticleList)(nil)

// APIArticleList 記事一覧のインタラクター.
type APIArticleList struct {
	articleRPC rpc.Article // 記事のリポジトリ
}

// NewAPIArticleList 記事一覧のインタラクターのファクトリ関数.
func NewAPIArticleList(
	articleRPC rpc.Article,
) *APIArticleList {
	return &APIArticleList{
		articleRPC: articleRPC,
	}
}

// Execute 記事一覧のインタラクターを実行する.
func (aal *APIArticleList) Execute(
	ctx context.Context,
	input port.APIArticleListInput,
) (port.APIArticleListOutput, error) {
	articles, err := aal.articleRPC.List(ctx, input.Index, input.Size)
	if err != nil {
		return port.APIArticleListOutput{}, err
	}

	return port.APIArticleListOutput{
		Articles: articles,
	}, nil
}
