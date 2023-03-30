package helper

import (
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/pkg/connect/article/v1/articlev1connect"
	"github.com/morning-night-guild/platform-app/pkg/connect/health/v1/healthv1connect"
	"github.com/morning-night-guild/platform-app/pkg/connect/user/v1/userv1connect"
)

type ConnectClient struct {
	Article articlev1connect.ArticleServiceClient
	User    userv1connect.UserServiceClient
	Health  healthv1connect.HealthServiceClient
}

func NewConnectClient(t *testing.T, client connect.HTTPClient, url string) *ConnectClient {
	t.Helper()

	ac := articlev1connect.NewArticleServiceClient(
		client,
		url,
	)

	uc := userv1connect.NewUserServiceClient(
		client,
		url,
	)

	hc := healthv1connect.NewHealthServiceClient(
		client,
		url,
	)

	return &ConnectClient{
		Article: ac,
		User:    uc,
		Health:  hc,
	}
}
