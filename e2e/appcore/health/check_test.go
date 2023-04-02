package health_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/e2e/helper"
	healthv1 "github.com/morning-night-guild/platform-app/pkg/connect/health/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

func TestAppCoreE2EHealthCheck(t *testing.T) {
	t.Parallel()

	url := helper.GetAppCoreEndpoint(t)

	t.Run("ヘルスチェックが成功する", func(t *testing.T) {
		t.Parallel()

		client := helper.NewConnectClient(t, http.DefaultClient, url)

		req := &healthv1.CheckRequest{}

		_, err := client.Health.Check(context.Background(), connect.NewRequest(req))
		if err != nil {
			t.Errorf("failed to health check: %s", err)
		}
	})

	t.Run("healthチェックが成功する", func(t *testing.T) {
		t.Parallel()

		// http://localhost:8888
		target := strings.Split(url, "//")[1]

		conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatal("failed to connect to grpc server: ", err)
		}

		defer conn.Close()

		client := health.NewHealthClient(conn)

		req := &health.HealthCheckRequest{}

		res, err := client.Check(context.Background(), req)
		if err != nil {
			t.Fatal("failed to health check: ", err)
		}

		if res.Status != health.HealthCheckResponse_SERVING {
			t.Errorf("health check status = %v, want %v", res.Status, health.HealthCheckResponse_SERVING)
		}
	})
}
