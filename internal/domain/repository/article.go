package repository

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

type Article interface {
	Save(context.Context, model.Article) error
	FindAll(context.Context, value.Index, value.Size) ([]model.Article, error)
}
