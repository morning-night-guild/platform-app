package api

import (
	"github.com/morning-night-guild/platform-app/pkg/connect/proto/article/v1/articlev1connect"
	"github.com/morning-night-guild/platform-app/pkg/connect/proto/health/v1/healthv1connect"
)

type ConnectFactory interface {
	Of(string) (*Connect, error)
}

type Connect struct {
	Article articlev1connect.ArticleServiceClient
	Health  healthv1connect.HealthServiceClient
}
