package main

import (
	"context"
	"os"

	"github.com/morning-night-guild/platform-app/internal/driver/postgres"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")

	rdb, err := postgres.New().Of(dsn)
	if err != nil {
		panic(err)
	}
	defer rdb.Close()

	ctx := context.Background()

	if err := rdb.Debug().Schema.Create(ctx); err != nil {
		panic(err)
	}
}
