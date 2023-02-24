package repository

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
)

type APIArticle interface {
	Save(context.Context, article.URL, article.Title, article.Description, article.Thumbnail) (model.Article, error)
	FindAll(context.Context, Index, Size) ([]model.Article, error)
}
