package port

import (
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
	"github.com/morning-night-guild/platform-app/internal/usecase"
)

// APIArticleShareInput.
type APIArticleShareInput struct {
	usecase.Input
	URL         article.URL
	Title       article.Title
	Description article.Description
	Thumbnail   article.Thumbnail
}

// APIArticleShareOutput.
type APIArticleShareOutput struct {
	usecase.Output
	Article model.Article
}

// APIArticleShare.
type APIArticleShare interface {
	usecase.Usecase[APIArticleShareInput, APIArticleShareOutput]
}

// APIArticleListInput.
type APIArticleListInput struct {
	usecase.Input
	Index value.Index
	Size  value.Size
}

// APIArticleListOutput.
type APIArticleListOutput struct {
	usecase.Output
	Articles []model.Article
}

// APIArticleList.
type APIArticleList interface {
	usecase.Usecase[APIArticleListInput, APIArticleListOutput]
}
