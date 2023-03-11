package connect

import (
	"net/http"

	"github.com/morning-night-guild/platform-app/internal/adapter/external"
	"github.com/morning-night-guild/platform-app/pkg/connect/article/v1/articlev1connect"
	"github.com/morning-night-guild/platform-app/pkg/connect/health/v1/healthv1connect"
)

var _ external.ConnectFactory = (*Client)(nil)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (cl *Client) Of(url string) (*external.Connect, error) {
	client := http.DefaultClient

	ac := articlev1connect.NewArticleServiceClient(
		client,
		url,
	)

	hc := healthv1connect.NewHealthServiceClient(
		client,
		url,
	)

	return &external.Connect{
		Article: ac,
		Health:  hc,
	}, nil
}
