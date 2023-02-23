package connect

import (
	"net/http"

	"github.com/morning-night-guild/platform-app/internal/adapter/gateway"
	"github.com/morning-night-guild/platform-app/pkg/connect/article/v1/articlev1connect"
	"github.com/morning-night-guild/platform-app/pkg/connect/health/v1/healthv1connect"
)

var _ gateway.ConnectFactory = (*Client)(nil)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) Of(url string) (*gateway.Connect, error) {
	client := http.DefaultClient

	ac := articlev1connect.NewArticleServiceClient(
		client,
		url,
	)

	hc := healthv1connect.NewHealthServiceClient(
		client,
		url,
	)

	return &gateway.Connect{
		Article: ac,
		Health:  hc,
	}, nil
}
