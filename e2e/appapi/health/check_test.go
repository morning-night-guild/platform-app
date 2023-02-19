package health_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/morning-night-guild/platform-app/e2e/helper"
)

func TestAppAPIE2EHealthCheck(t *testing.T) {
	t.Parallel()

	url := helper.GetAppAPIEndpoint(t)

	t.Run("apiヘルスチェックが成功する", func(t *testing.T) {
		t.Parallel()

		client := helper.NewOpenAPIClient(t, url)

		res, err := client.Client.V1HealthAPI(context.Background())
		if err != nil {
			t.Fatalf("failed to health check: %s", err)
		}

		defer res.Body.Close()

		if !reflect.DeepEqual(res.StatusCode, http.StatusOK) {
			t.Errorf("Articles actual = %v, want %v", res.StatusCode, http.StatusOK)
		}
	})

	t.Run("coreヘルスチェックが成功する", func(t *testing.T) {
		t.Parallel()

		client := helper.NewOpenAPIClient(t, url)

		res, err := client.Client.V1HealthCore(context.Background())
		if err != nil {
			t.Fatalf("failed to health check: %s", err)
		}

		defer res.Body.Close()

		if !reflect.DeepEqual(res.StatusCode, http.StatusOK) {
			t.Errorf("Articles actual = %v, want %v", res.StatusCode, http.StatusOK)
		}
	})
}
