package connect

import (
	"net/http"

	"github.com/morning-night-guild/platform-app/internal/adapter/api"
	"github.com/morning-night-guild/platform-app/pkg/connect/proto/article/v1/articlev1connect"
	"github.com/morning-night-guild/platform-app/pkg/connect/proto/health/v1/healthv1connect"
)

var _ api.ConnectFactory = (*Client)(nil)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) Of(url string) (*api.Connect, error) {
	client := http.DefaultClient

	ac := articlev1connect.NewArticleServiceClient(
		client,
		url,
	)

	hc := healthv1connect.NewHealthServiceClient(
		client,
		url,
	)

	return &api.Connect{
		Article: ac,
		Health:  hc,
	}, nil
}
