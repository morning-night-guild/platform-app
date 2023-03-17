package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/morning-night-guild/platform-app/internal/driver/postgres"
	"github.com/morning-night-guild/platform-app/pkg/ent"
)

const LocalDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

type Mode string

func main() {
	flag.Parse()

	mode := flag.Arg(0)

	dsn := flag.Arg(1)
	if dsn == "" {
		panic("dsn must be set")
	}

	log.Printf("mode: %s, dsn: %s\n", mode, dsn)

	switch mode {
	case "export":
		if err := Execute(context.Background(), dsn, LocalDSN); err != nil {
			panic(err)
		}
	case "import":
		if err := Execute(context.Background(), LocalDSN, dsn); err != nil {
			panic(err)
		}
	default:
		panic("invalid mode")
	}
}

func Execute(ctx context.Context, src string, dst string) error {
	srcCli, err := postgres.New().Of(src)
	if err != nil {
		return fmt.Errorf("failed to open src client: %w", err)
	}

	defer srcCli.Close()

	dstCli, err := postgres.New().Of(dst)
	if err != nil {
		return fmt.Errorf("failed to open dst client: %w", err)
	}

	defer dstCli.Close()

	articles, err := srcCli.Debug().Article.Query().All(ctx)
	if err != nil {
		return fmt.Errorf("failed to query articles: %w", err)
	}

	log.Printf("completed select %d articles from %s\n", len(articles), src)

	bulk := make([]*ent.ArticleCreate, len(articles))
	for i, article := range articles {
		bulk[i] = dstCli.Article.Create().
			SetID(article.ID).
			SetTitle(article.Title).
			SetURL(article.URL).
			SetDescription(article.Description).
			SetThumbnail(article.Thumbnail).
			SetCreatedAt(article.CreatedAt).
			SetUpdatedAt(article.UpdatedAt)
	}

	if _, err := dstCli.Debug().Article.CreateBulk(bulk...).Save(ctx); err != nil {
		return fmt.Errorf("failed to bulk create articles: %w", err)
	}

	log.Printf("completed insert %d articles to %s\n", len(articles), dst)

	return nil
}
