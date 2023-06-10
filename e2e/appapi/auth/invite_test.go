package auth_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/e2e/helper"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

func TestE2EAuthInvite(t *testing.T) {
	t.Parallel()

	url := helper.GetAppAPIEndpoint(t)

	t.Run("招待できる", func(t *testing.T) {
		t.Parallel()

		client := helper.NewOpenAPIClient(t, url)

		client.Client.Client = &http.Client{
			Transport: helper.NewAPIKeyTransport(t, helper.GetAPIKey(t)),
		}

		id := uuid.New().String()

		email := fmt.Sprintf("%s@example.com", id)

		res, err := client.Client.V1AuthInvite(context.Background(), openapi.V1AuthInviteRequestSchema{
			Email: types.Email(email),
		})
		if err != nil {
			t.Fatalf("failed to auth invite: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("failed to auth invite: %d", res.StatusCode)
		}

		body, _ := io.ReadAll(res.Body)

		var invitation openapi.V1AuthInviteResponseSchema
		if err := json.Unmarshal(body, &invitation); err != nil {
			t.Fatalf("failed marshal response: %s caused by %s", body, err)

			return
		}

		if len(invitation.Code) != 8 {
			t.Errorf("failed to auth invite: %s", invitation.Code)
		}

		rds := helper.NewRedis(t, helper.GetRedisURL(t))

		defer rds.Close()

		defer rds.DeleteInvitationCode(invitation.Code)
	})

	t.Run("Api-Keyがなく招待できない", func(t *testing.T) {
		t.Parallel()

		client := helper.NewOpenAPIClient(t, url)

		id := uuid.New().String()

		email := fmt.Sprintf("%s@example.com", id)

		res, err := client.Client.V1AuthInvite(context.Background(), openapi.V1AuthInviteRequestSchema{
			Email: types.Email(email),
		})
		if err != nil {
			t.Fatalf("failed to auth invite: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusUnauthorized {
			t.Fatalf("failed to auth invite: %d", res.StatusCode)
		}
	})
}
