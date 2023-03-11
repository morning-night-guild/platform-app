package external

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/pkg/connect/article/v1/articlev1connect"
	"github.com/morning-night-guild/platform-app/pkg/connect/health/v1/healthv1connect"
	"github.com/morning-night-guild/platform-app/pkg/trace"
)

type ConnectFactory interface {
	Of(string) (*Connect, error)
}

type Connect struct {
	Article articlev1connect.ArticleServiceClient
	Health  healthv1connect.HealthServiceClient
}

func NewRequestWithTID[T any](ctx context.Context, msg *T) *connect.Request[T] {
	req := connect.NewRequest(msg)

	req.Header().Set("tid", trace.GetTIDCtx(ctx))

	return req
}
