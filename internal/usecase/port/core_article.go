package port

import (
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase"
)

// CoreArticleShareInput.
type CoreArticleShareInput struct {
	usecase.Input
	URL         article.URL
	Title       article.Title
	Description article.Description
	Thumbnail   article.Thumbnail
}

// CoreArticleShareOutput.
type CoreArticleShareOutput struct {
	usecase.Output
	Article model.Article
}

// CoreArticleShare.
type CoreArticleShare interface {
	usecase.Usecase[CoreArticleShareInput, CoreArticleShareOutput]
}

// CoreArticleListInput.
type CoreArticleListInput struct {
	usecase.Input
	Index repository.Index
	Size  repository.Size
}

// CoreArticleListOutput.
type CoreArticleListOutput struct {
	usecase.Output
	Articles []model.Article
}

// CoreArticleList.
type CoreArticleList interface {
	usecase.Usecase[CoreArticleListInput, CoreArticleListOutput]
}
