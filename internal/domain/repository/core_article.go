package repository

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
)

type CoreArticle interface {
	Save(context.Context, model.Article) error
	FindAll(context.Context, Index, Size) ([]model.Article, error)
}
