package main

import (
	"context"
	"fmt"
	"os"

	"github.com/morning-night-guild/platform-app/internal/adapter/gateway"
	"github.com/morning-night-guild/platform-app/internal/driver/postgres"
	"github.com/morning-night-guild/platform-app/pkg/ent"
	"github.com/morning-night-guild/platform-app/pkg/ent/article"
	"github.com/morning-night-guild/platform-app/pkg/ent/articletag"
	"github.com/morning-night-guild/platform-app/pkg/ent/user"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

func main() {
	ctx := context.Background()

	primaryDSN := os.Getenv("PRIMARY_DATABASE_URL")

	secondaryDSN := os.Getenv("SECONDARY_DATABASE_URL")

	primary, err := postgres.New().Of(primaryDSN)
	if err != nil {
		log.GetLogCtx(ctx).Panic("failed to connect to primary database", log.ErrorField(err))
	}

	defer primary.Close()

	entity, err := Export(ctx, primary)
	if err != nil {
		log.GetLogCtx(ctx).Panic("failed to export", log.ErrorField(err))
	}

	secondary, err := postgres.New().Of(secondaryDSN)
	if err != nil {
		log.GetLogCtx(ctx).Panic("failed to connect to secondary database", log.ErrorField(err))
	}

	defer secondary.Close()

	log.GetLogCtx(ctx).Info("start transaction")

	tx, err := secondary.Tx(ctx)
	if err != nil {
		log.GetLogCtx(ctx).Panic("failed to begin transaction", log.ErrorField(err))
	}

	if err := ResetTable(ctx, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			log.GetLogCtx(ctx).Panic("failed to rollback transaction", log.ErrorField(err))
		}

		log.GetLogCtx(ctx).Panic("failed to drop table", log.ErrorField(err))
	}

	if err := Import(ctx, tx, entity); err != nil {
		if err := tx.Rollback(); err != nil {
			log.GetLogCtx(ctx).Panic("failed to rollback transaction", log.ErrorField(err))
		}

		log.GetLogCtx(ctx).Panic("failed to import", log.ErrorField(err))
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			log.GetLogCtx(ctx).Panic("failed to rollback transaction", log.ErrorField(err))
		}

		log.GetLogCtx(ctx).Panic("failed to commit transaction", log.ErrorField(err))
	}

	log.GetLogCtx(ctx).Info("end transaction")

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

const dropTableQuery = "DROP TABLE IF EXISTS %s CASCADE"

//nolint:cyclop
func ResetTable(ctx context.Context, tx *ent.Tx) error {
	log.GetLogCtx(ctx).Info("start drop table")

	if _, err := tx.ExecContext(ctx, fmt.Sprintf(dropTableQuery, user.Table)); err != nil {
		return fmt.Errorf("failed to drop user table: %w", err)
	}

	if _, err := tx.ExecContext(ctx, fmt.Sprintf(dropTableQuery, article.Table)); err != nil {
		return fmt.Errorf("failed to drop article table: %w", err)
	}

	if _, err := tx.ExecContext(ctx, fmt.Sprintf(dropTableQuery, articletag.Table)); err != nil {
		return fmt.Errorf("failed to drop article tag table: %w", err)
	}

	if err := tx.Client().Debug().Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed to create primary schema: %w", err)
	}

	log.GetLogCtx(ctx).Info("end drop table")

	return nil
}

//nolint:funlen
func Import(ctx context.Context, tx *ent.Tx, entity Entity) error { //nolint:cyclop
	log.GetLogCtx(ctx).Info("start import data")

	userBulk := make([]*ent.UserCreate, len(entity.Users))
	for i, user := range entity.Users {
		userBulk[i] = tx.User.Create().
			SetID(user.ID).
			SetCreatedAt(user.CreatedAt).
			SetUpdatedAt(user.UpdatedAt)
	}

	if _, err := tx.User.CreateBulk(userBulk...).Save(ctx); err != nil {
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
		return fmt.Errorf("failed to bulk create article tags: %w", err)
	}

	log.GetLogCtx(ctx).Info("end import data")

	return nil
}
