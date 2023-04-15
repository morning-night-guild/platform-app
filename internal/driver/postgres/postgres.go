package postgres

import (
	"context"

	"entgo.io/ent/dialect"
	// postgres driver.
	_ "github.com/lib/pq"
	"github.com/morning-night-guild/platform-app/internal/adapter/gateway"
	"github.com/morning-night-guild/platform-app/pkg/ent"
)

var _ gateway.RDBFactory = (*Postgres)(nil)

type Postgres struct{}

func New() *Postgres {
	return &Postgres{}
}

func (c Postgres) Of(dsn string) (*gateway.RDB, error) {
	client, err := ent.Open(dialect.Postgres, dsn)
	if err != nil {
		return &gateway.RDB{}, err
	}

	if _, err := client.QueryContext(context.Background(), "SELECT 1"); err != nil {
		return &gateway.RDB{}, err
	}

	return &gateway.RDB{
		Client: client,
	}, err
}
