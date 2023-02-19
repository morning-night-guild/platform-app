package helper

import (
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/pkg/connect/proto/article/v1/articlev1connect"
	"github.com/morning-night-guild/platform-app/pkg/connect/proto/health/v1/healthv1connect"
)

type ConnectClient struct {
	Article articlev1connect.ArticleServiceClient
	Health  healthv1connect.HealthServiceClient
}

func NewConnectClient(t *testing.T, client connect.HTTPClient, url string) *ConnectClient {
	t.Helper()

	ac := articlev1connect.NewArticleServiceClient(
		client,
		url,
	)

	hc := healthv1connect.NewHealthServiceClient(
		client,
		url,
	)

	return &ConnectClient{
		Article: ac,
		Health:  hc,
	}
}
