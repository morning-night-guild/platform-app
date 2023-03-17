package main

import (
	"context"
	"os"

	"github.com/morning-night-guild/platform-app/internal/driver/postgres"
	"github.com/morning-night-guild/platform-app/pkg/ent"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

func main() {
	primaryDSN := os.Getenv("PRIMARY_DATABASE_URL")

	secondaryDSN := os.Getenv("SECONDARY_DATABASE_URL")

	primary, err := postgres.New().Of(primaryDSN)
	if err != nil {
		log.Log().Panic("failed to connect to primary database", log.ErrorField(err))
	}

	defer primary.Close()

	secondary, err := postgres.New().Of(secondaryDSN)
	if err != nil {
		log.Log().Panic("failed to connect to secondary database", log.ErrorField(err))
	}

	defer secondary.Close()

	ctx := context.Background()

	if err := secondary.Debug().Schema.Create(ctx); err != nil {
		log.Log().Panic("failed to create secondary schema", log.ErrorField(err))
	}

	articleTags, err := primary.ArticleTag.Query().All(ctx)
	if err != nil {
		log.Log().Panic("failed to query article tags", log.ErrorField(err))
	}

	articles, err := primary.Article.Query().All(ctx)
	if err != nil {
		log.Log().Panic("failed to query articles", log.ErrorField(err))
	}

	log.GetLogCtx(ctx).Info("start transaction")

	tx, err := secondary.BeginTx(ctx, nil)
	if err != nil {
		log.Log().Panic("failed to begin transaction", log.ErrorField(err))
	}

	if _, err := tx.ArticleTag.Delete().Exec(ctx); err != nil {
		tx.Rollback()

		log.Log().Panic("failed to delete article tags", log.ErrorField(err))
	}

	if _, err := tx.Article.Delete().Exec(ctx); err != nil {
		tx.Rollback()

		log.Log().Panic("failed to delete articles", log.ErrorField(err))
	}

	articleBulk := make([]*ent.ArticleCreate, 0, len(articles))
	for i, article := range articles {
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
		tx.Rollback()

		log.Log().Panic("failed to bulk create articles", log.ErrorField(err))
	}

	articleTagBulk := make([]*ent.ArticleTagCreate, 0, len(articleTags))
	for i, articleTag := range articleTags {
		articleTagBulk[i] = tx.ArticleTag.Create().
			SetID(articleTag.ID).
			SetTag(articleTag.Tag).
			SetArticleID(articleTag.ArticleID)
	}

	if _, err := tx.ArticleTag.CreateBulk(articleTagBulk...).Save(ctx); err != nil {
		tx.Rollback()

		log.Log().Panic("failed to bulk create article tags", log.ErrorField(err))
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()

		log.Log().Panic("failed to commit transaction", log.ErrorField(err))
	}

	log.GetLogCtx(ctx).Info("commit transaction")
}
