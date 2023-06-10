package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
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

func TestE2EAuthJoin(t *testing.T) {
	t.Parallel()

	url := helper.GetAppAPIEndpoint(t)

	t.Run("参加できる", func(t *testing.T) {
		t.Parallel()

		// 招待する
		client := helper.NewOpenAPIClient(t, url)

		client.Client.Client = &http.Client{
			Transport: helper.NewAPIKeyTransport(t, helper.GetAPIKey(t)),
		}

		id := uuid.New().String()

		email := fmt.Sprintf("%s@example.com", id)

		tmp, err := client.Client.V1AuthInvite(context.Background(), openapi.V1AuthInviteRequestSchema{
			Email: types.Email(email),
		})
		if err != nil {
			t.Fatalf("failed to auth invite: %s", err)
		}

		defer tmp.Body.Close()

		if tmp.StatusCode != http.StatusOK {
			t.Fatalf("failed to auth invite: %d", tmp.StatusCode)
		}

		body, _ := io.ReadAll(tmp.Body)

		var invitation openapi.V1AuthInviteResponseSchema
		if err := json.Unmarshal(body, &invitation); err != nil {
			t.Fatalf("failed marshal response: %s caused by %s", body, err)

			return
		}

		// 招待コードで参加する
		client = helper.NewOpenAPIClient(t, url)

		password := id

		res, err := client.Client.V1AuthJoin(context.Background(), openapi.V1AuthJoinRequestSchema{
			Code:     invitation.Code,
			Password: password,
		})
		if err != nil {
			t.Fatalf("failed to auth join: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("failed to auth join: %d", res.StatusCode)
		}

		// サインインする
		defer func() {
			prv, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				t.Fatalf("failed to generate key: %s", err)
			}

			res, err := client.Client.V1AuthSignIn(context.Background(), openapi.V1AuthSignInJSONRequestBody{
				Email:     types.Email(email),
				Password:  password,
				PublicKey: helper.Public(t, prv),
			})
			if err != nil {
				t.Errorf("failed to auth sign in: %s", err)
			}

			defer res.Body.Close()

			uid := helper.ExtractUserID(t, res.Cookies())

			db := helper.NewDatabase(t, helper.GetDSN(t))

			defer db.Close()

			db.DeleteUser(uuid.MustParse(uid))
		}()
	})

	t.Run("招待コードが存在せず参加できない", func(t *testing.T) {
		t.Parallel()

		// 招待コードで参加する
		client := helper.NewOpenAPIClient(t, url)

		password := uuid.New().String()

		res, err := client.Client.V1AuthJoin(context.Background(), openapi.V1AuthJoinRequestSchema{
			Code:     "invalid",
			Password: password,
		})
		if err != nil {
			t.Fatalf("failed to auth join: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("failed to auth join: %d", res.StatusCode)
		}
	})
}
