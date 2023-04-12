package gateway_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/integration/helper"
	"github.com/morning-night-guild/platform-app/internal/driver/postgres"
)

func TestGateway(t *testing.T) {
	t.Parallel()

	dsn := helper.GetDSN(t)

	t.Run("PostgreSQLに接続できる", func(t *testing.T) {
		t.Parallel()

		if _, err := postgres.New().Of(dsn); err != nil {
			t.Fatalf("failed to connect to postgres: %v", err)
		}
	})
}
