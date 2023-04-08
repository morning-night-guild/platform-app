package usecase

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

//go:generate mockgen -source api_article.go -destination api_article_mock.go -package usecase

// APIArticle.
type APIArticle interface {
	Share(context.Context, APIArticleShareInput) (APIArticleShareOutput, error)
	List(context.Context, APIArticleListInput) (APIArticleListOutput, error)
}

// APIArticleShareInput.
type APIArticleShareInput struct {
	URL         article.URL
	Title       article.Title
	Description article.Description
	Thumbnail   article.Thumbnail
}

// APIArticleShareOutput.
type APIArticleShareOutput struct {
	Article model.Article
}

// APIArticleListInput.
type APIArticleListInput struct {
	Index value.Index
	Size  value.Size
}

// APIArticleListOutput.
type APIArticleListOutput struct {
	Articles []model.Article
}
