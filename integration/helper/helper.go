package helper

import (
	"os"
	"testing"
)

func GetRedisURL(t *testing.T) string {
	t.Helper()

	url := os.Getenv("TEST_REDIS_URL")

	if url != "" {
		return url
	}

	url = os.Getenv("REDIS_URL")

	if url != "" {
		return url
	}

	if url == "" {
		t.Fatal("TEST_REDIS_URL or REDIS_URL is not set")
	}

	return url
}

func GetPostgresURL(t *testing.T) string {
	t.Helper()

	url := os.Getenv("TEST_POSTGRES_URL")

	if url != "" {
		return url
	}

	url = os.Getenv("POSTGRES_URL")

	if url != "" {
		return url
	}

	if url == "" {
		t.Fatal("TEST_POSTGRES_URL or POSTGRES_URL is not set")
	}

	return url
}

func GetDSN(t *testing.T) string {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")

	if dsn != "" {
		return dsn
	}

	dsn = os.Getenv("DATABASE_URL")

	if dsn != "" {
		return dsn
	}

	if dsn == "" {
		t.Fatal("TEST_DATABASE_URL or DATABASE_URL is not set")
	}

	return ""
}
