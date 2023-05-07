package usecase

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

//go:generate mockgen -source core_article.go -destination core_article_mock.go -package usecase

// CoreArticle.
type CoreArticle interface {
	Share(context.Context, CoreArticleShareInput) (CoreArticleShareOutput, error)
	List(context.Context, CoreArticleListInput) (CoreArticleListOutput, error)
	Delete(context.Context, CoreArticleDeleteInput) (CoreArticleDeleteOutput, error)
}

// CoreArticleShareInput.
type CoreArticleShareInput struct {
	URL         article.URL
	Title       article.Title
	Description article.Description
	Thumbnail   article.Thumbnail
}

// CoreArticleShareOutput.
type CoreArticleShareOutput struct {
	Article model.Article
}

// CoreArticleListInput.
type CoreArticleListInput struct {
	Index  value.Index
	Size   value.Size
	Filter []value.Filter
}

// CoreArticleListOutput.
type CoreArticleListOutput struct {
	Articles []model.Article
}

// CoreArticleDeleteInput.
type CoreArticleDeleteInput struct {
	ArticleID article.ID
}

// CoreArticleDeleteOutput.
type CoreArticleDeleteOutput struct{}
