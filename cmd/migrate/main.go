package main

import (
	"context"
	"os"

	_ "github.com/lib/pq"
	"github.com/morning-night-guild/platform-app/pkg/ent"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	ctx := context.Background()

	if err := client.Debug().Schema.Create(ctx); err != nil {
		panic(err)
	}
}
