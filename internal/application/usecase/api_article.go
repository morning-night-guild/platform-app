package usecase

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

//go:generate mockgen -source api_article.go -destination api_article_mock.go -package usecase

// APIArticle.
type APIArticle interface {
	Share(context.Context, APIArticleShareInput) (APIArticleShareOutput, error)
	List(context.Context, APIArticleListInput) (APIArticleListOutput, error)
	Delete(context.Context, APIArticleDeleteInput) (APIArticleDeleteOutput, error)
	AddToUser(context.Context, APIArticleAddToUserInput) (APIArticleAddToUserOutput, error)
	RemoveFromUser(context.Context, APIArticleRemoveFromUserInput) (APIArticleRemoveFromUserOutput, error)
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
	UserID user.ID
	Scope  article.Scope
	Index  value.Index
	Size   value.Size
	Filter []value.Filter
}

// APIArticleListOutput.
type APIArticleListOutput struct {
	Articles []model.Article
}

// APIArticleDeleteInput.
type APIArticleDeleteInput struct {
	ArticleID article.ID
}

// APIArticleDeleteOutput.
type APIArticleDeleteOutput struct{}

// APIArticleAddToUserInput.
type APIArticleAddToUserInput struct {
	ArticleID article.ID
	UserID    user.ID
}

// APIArticleAddToUserOutput.
type APIArticleAddToUserOutput struct{}

// APIArticleRemoveFromUserInput.
type APIArticleRemoveFromUserInput struct {
	ArticleID article.ID
	UserID    user.ID
}

// APIArticleRemoveFromUserOutput.
type APIArticleRemoveFromUserOutput struct{}
