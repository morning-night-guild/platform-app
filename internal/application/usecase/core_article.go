package usecase

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

//go:generate mockgen -source core_article.go -destination core_article_mock.go -package usecase

// CoreArticle.
type CoreArticle interface {
	Share(context.Context, CoreArticleShareInput) (CoreArticleShareOutput, error)
	List(context.Context, CoreArticleListInput) (CoreArticleListOutput, error)
	ListByUser(context.Context, CoreArticleListByUserInput) (CoreArticleListByUserOutput, error)
	Delete(context.Context, CoreArticleDeleteInput) (CoreArticleDeleteOutput, error)
	AddToUser(context.Context, CoreArticleAddToUserInput) (CoreArticleAddToUserOutput, error)
	RemoveFromUser(context.Context, CoreArticleRemoveFromUserInput) (CoreArticleRemoveFromUserOutput, error)
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

// CoreArticleListByUserInput.
type CoreArticleListByUserInput struct {
	UserID user.ID
	Index  value.Index
	Size   value.Size
	Filter []value.Filter
}

// CoreArticleListByUserOutput.
type CoreArticleListByUserOutput struct {
	Articles []model.Article
}

// CoreArticleDeleteInput.
type CoreArticleDeleteInput struct {
	ArticleID article.ID
}

// CoreArticleDeleteOutput.
type CoreArticleDeleteOutput struct{}

// CoreArticleAddToUserInput.
type CoreArticleAddToUserInput struct {
	ArticleID article.ID
	UserID    user.ID
}

// CoreArticleAddToUserOutput.
type CoreArticleAddToUserOutput struct{}

// CoreArticleRemoveFromUserInput.
type CoreArticleRemoveFromUserInput struct {
	ArticleID article.ID
	UserID    user.ID
}

// CoreArticleRemoveFromUserOutput.
type CoreArticleRemoveFromUserOutput struct{}
