package repository

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

//go:generate mockgen -source article.go -destination article_mock.go -package repository

type Article interface {
	Save(context.Context, model.Article) error
	FindAll(context.Context, value.Index, value.Size) ([]model.Article, error)
}
