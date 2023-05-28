package repository

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

//go:generate mockgen -source article.go -destination article_mock.go -package repository

type Article interface {
	Save(context.Context, model.Article) error
	List(context.Context, value.Index, value.Size, ...value.Filter) ([]model.Article, error)
	ListByUser(context.Context, user.ID, value.Index, value.Size, ...value.Filter) ([]model.Article, error)
	Find(context.Context, article.ID) (model.Article, error)
	Delete(context.Context, article.ID) error
	ExistsByUser(context.Context, article.ID, user.ID) (bool, error)
	AddToUser(context.Context, article.ID, user.ID) error
	RemoveFromUser(context.Context, article.ID, user.ID) error
}
