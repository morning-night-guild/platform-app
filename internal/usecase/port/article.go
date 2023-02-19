package port

import (
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase"
)

// ShareArticleInput.
type ShareArticleInput struct {
	usecase.Input
	URL         article.URL
	Title       article.Title
	Description article.Description
	Thumbnail   article.Thumbnail
}

// ShareArticleOutput.
type ShareArticleOutput struct {
	usecase.Output
	Article model.Article
}

// ShareArticle.
type ShareArticle interface {
	usecase.Usecase[ShareArticleInput, ShareArticleOutput]
}

// ListArticleInput.
type ListArticleInput struct {
	usecase.Input
	Index repository.Index
	Size  repository.Size
}

// ListArticleOutput.
type ListArticleOutput struct {
	usecase.Output
	Articles []model.Article
}

// ListArticle.
type ListArticle interface {
	usecase.Usecase[ListArticleInput, ListArticleOutput]
}
