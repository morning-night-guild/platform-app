package rpc

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

//go:generate mockgen -source article.go -destination article_mock.go -package rpc

type Article interface {
	Share(context.Context, article.URL, article.Title, article.Description, article.Thumbnail) (model.Article, error)
	List(context.Context, value.Index, value.Size, ...value.Filter) ([]model.Article, error)
	ListByUser(context.Context, user.ID, value.Index, value.Size, ...value.Filter) ([]model.Article, error)
	Delete(context.Context, article.ID) error
	AddToUser(context.Context, article.ID, user.ID) error
}
