package helper

import (
	"os"
	"testing"
)

func GetAppCoreEndpoint(t *testing.T) string {
	t.Helper()

	return os.Getenv("APP_CORE_ENDPOINT")
}

func GetAppAPIEndpoint(t *testing.T) string {
	t.Helper()

	return os.Getenv("APP_API_ENDPOINT")
}

func GetAPIKey(t *testing.T) string {
	t.Helper()

	return os.Getenv("API_KEY")
}

func GetDSN(t *testing.T) string {
	t.Helper()

	return os.Getenv("DATABASE_URL")
}
