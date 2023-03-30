package main

import (
	"context"
	"fmt"
	"os"

	"github.com/morning-night-guild/platform-app/internal/adapter/gateway"
	"github.com/morning-night-guild/platform-app/internal/driver/postgres"
	"github.com/morning-night-guild/platform-app/pkg/ent"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

func main() {
	ctx := context.Background()

	primaryDSN := os.Getenv("PRIMARY_DATABASE_URL")

	secondaryDSN := os.Getenv("SECONDARY_DATABASE_URL")

	primary, err := postgres.New().Of(primaryDSN)
	if err != nil {
		log.Log().Panic("failed to connect to primary database", log.ErrorField(err))
	}

	defer primary.Close()

	entity, err := Export(ctx, primary)
	if err != nil {
		log.Log().Panic("failed to export", log.ErrorField(err))
	}

	secondary, err := postgres.New().Of(secondaryDSN)
	if err != nil {
		log.Log().Panic("failed to connect to secondary database", log.ErrorField(err))
	}

	defer secondary.Close()

	if err := secondary.Debug().Schema.Create(ctx); err != nil {
		log.Log().Panic("failed to create secondary schema", log.ErrorField(err))
	}

	if err := Import(ctx, secondary, entity); err != nil {
		log.Log().Panic("failed to import", log.ErrorField(err))
	}

	log.GetLogCtx(ctx).Info("success backup")
}

type Entity struct {
	Users       []*ent.User
	Articles    []*ent.Article
	ArticleTags []*ent.ArticleTag
}

func Export(ctx context.Context, client *gateway.RDB) (Entity, error) {
	log.GetLogCtx(ctx).Info("start export")

	users, err := client.User.Query().All(ctx)
	if err != nil {
		return Entity{}, fmt.Errorf("failed to query users: %w", err)
	}

	articleTags, err := client.ArticleTag.Query().All(ctx)
	if err != nil {
		return Entity{}, fmt.Errorf("failed to query article tags: %w", err)
	}

	articles, err := client.Article.Query().All(ctx)
	if err != nil {
		return Entity{}, fmt.Errorf("failed to query articles: %w", err)
	}

	return Entity{
		Users:       users,
		Articles:    articles,
		ArticleTags: articleTags,
	}, nil
}

//nolint:funlen
func Import(ctx context.Context, client *gateway.RDB, entity Entity) error { //nolint:cyclop
	log.GetLogCtx(ctx).Info("start import")

	tx, err := client.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if _, err := tx.User.Delete().Exec(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}

		return fmt.Errorf("failed to delete users: %w", err)
	}

	if _, err := tx.ArticleTag.Delete().Exec(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}

		return fmt.Errorf("failed to delete article tags: %w", err)
	}

	if _, err := tx.Article.Delete().Exec(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}

		return fmt.Errorf("failed to delete articles: %w", err)
	}

	userBulk := make([]*ent.UserCreate, len(entity.Users))
	for i, user := range entity.Users {
		userBulk[i] = tx.User.Create().
			SetID(user.ID).
			SetCreatedAt(user.CreatedAt).
			SetUpdatedAt(user.UpdatedAt)
	}

	if _, err := tx.User.CreateBulk(userBulk...).Save(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}

		return fmt.Errorf("failed to bulk create users: %w", err)
	}

	articleBulk := make([]*ent.ArticleCreate, len(entity.Articles))
	for i, article := range entity.Articles {
		articleBulk[i] = tx.Article.Create().
			SetID(article.ID).
			SetTitle(article.Title).
			SetURL(article.URL).
			SetDescription(article.Description).
			SetThumbnail(article.Thumbnail).
			SetCreatedAt(article.CreatedAt).
			SetUpdatedAt(article.UpdatedAt)
	}

	if _, err := tx.Article.CreateBulk(articleBulk...).Save(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}

		return fmt.Errorf("failed to bulk create articles: %w", err)
	}

	articleTagBulk := make([]*ent.ArticleTagCreate, len(entity.ArticleTags))
	for i, articleTag := range entity.ArticleTags {
		articleTagBulk[i] = tx.ArticleTag.Create().
			SetID(articleTag.ID).
			SetTag(articleTag.Tag).
			SetArticleID(articleTag.ArticleID)
	}

	if _, err := tx.ArticleTag.CreateBulk(articleTagBulk...).Save(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}

		return fmt.Errorf("failed to bulk create article tags: %w", err)
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}

		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.GetLogCtx(ctx).Info("commit transaction")

	return nil
}
