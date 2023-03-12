package main

import (
	"context"
	"flag"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/morning-night-guild/platform-app/pkg/ent"
)

const LocalDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

type Mode string

func main() {
	flag.Parse()
	fmt.Println(flag.Args())

	mode := flag.Arg(0)

	dsn := flag.Arg(1)
	if dsn == "" {
		panic("dsn must be set")
	}

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
	srcCli, err := ent.Open("postgres", src)
	if err != nil {
		return fmt.Errorf("failed to open src client: %w", err)
	}

	defer srcCli.Close()

	articles, err := srcCli.Article.Query().All(ctx)
	if err != nil {
		return fmt.Errorf("failed to query articles: %w", err)
	}

	dstCli, err := ent.Open("postgres", dst)
	if err != nil {
		return fmt.Errorf("failed to open dst client: %w", err)
	}

	defer dstCli.Close()

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

	if _, err := dstCli.Article.CreateBulk(bulk...).Save(ctx); err != nil {
		return fmt.Errorf("failed to bulk create articles: %w", err)
	}

	return nil
}
