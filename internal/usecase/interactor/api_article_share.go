package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

var _ port.APIArticleShare = (*APIArticleShare)(nil)

// APIArticleShare 記事共有のインタラクター.
type APIArticleShare struct {
	articleRPC rpc.Article // 記事のリポジトリ
}

// NewAPIArticleShare 記事共有のインタラクターのファクトリ関数.
func NewAPIArticleShare(
	articleRPC rpc.Article,
) *APIArticleShare {
	return &APIArticleShare{
		articleRPC: articleRPC,
	}
}

// Execute 記事共有のインタラクターを実行する.
func (aas *APIArticleShare) Execute(
	ctx context.Context,
	input port.APIArticleShareInput,
) (port.APIArticleShareOutput, error) {
	res, err := aas.articleRPC.Share(
		ctx,
		input.URL,
		input.Title,
		input.Description,
		input.Thumbnail,
	)
	if err != nil {
		return port.APIArticleShareOutput{}, err
	}

	return port.APIArticleShareOutput{
		Article: res,
	}, nil
}
