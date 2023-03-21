package connect

import (
	"net/http"

	"github.com/morning-night-guild/platform-app/internal/adapter/external"
	"github.com/morning-night-guild/platform-app/pkg/connect/article/v1/articlev1connect"
	"github.com/morning-night-guild/platform-app/pkg/connect/health/v1/healthv1connect"
)

var (
	_ external.ArticleFactory = (*Connect)(nil)
	_ external.HealthFactory  = (*Connect)(nil)
)

type Connect struct {
	client *http.Client
}

func New() *Connect {
	return &Connect{
		client: http.DefaultClient,
	}
}

func (cn *Connect) Article(url string) (*external.Article, error) {
	return external.NewArticle(articlev1connect.NewArticleServiceClient(
		cn.client,
		url,
	)), nil
}

func (cn *Connect) Health(url string) (*external.Health, error) {
	return external.NewHealth(healthv1connect.NewHealthServiceClient(
		cn.client,
		url,
	)), nil
}
