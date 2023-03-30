package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/e2e/helper"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

func TestE2EAuthSighUp(t *testing.T) {
	t.Parallel()

	url := helper.GetAppAPIEndpoint(t)

	t.Run("サインアップできる", func(t *testing.T) {
		t.Parallel()

		client := helper.NewOpenAPIClient(t, url)

		client.Client.Client = &http.Client{
			Transport: helper.NewAPIKeyTransport(t, helper.GetAPIKey(t)),
		}

		id := uuid.New().String()

		email := fmt.Sprintf("%s@example.com", id)

		password := id

		res, err := client.Client.V1AuthSignUp(context.Background(), openapi.V1AuthSignUpJSONRequestBody{
			Email:    types.Email(email),
			Password: password,
		})
		if err != nil {
			t.Fatalf("failed to auth sign up: %s", err)
		}

		if res.StatusCode != http.StatusOK {
			t.Fatalf("failed to auth sign up: %d", res.StatusCode)
		}

		defer res.Body.Close()

		defer func() {
			prv, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Client.V1AuthSignIn(context.Background(), openapi.V1AuthSignInJSONRequestBody{
				Email:     types.Email(email),
				Password:  password,
				PublicKey: helper.Public(t, prv),
			})
			if err != nil {
				t.Fatalf("failed to auth sign in: %s", err)
			}

			defer res.Body.Close()

			uid := helper.ExtractUserID(t, res.Cookies())

			db := helper.NewDatabase(t, helper.GetDSN(t))

			defer db.Close()

			db.DeleteUser(uuid.MustParse(uid))
		}()
	})

	t.Run("Api-Keyがなくサインアップできない", func(t *testing.T) {
		t.Parallel()

		client := helper.NewOpenAPIClient(t, url)

		id := uuid.New().String()

		email := fmt.Sprintf("%s@example.com", id)

		password := id

		res, err := client.Client.V1AuthSignUp(context.Background(), openapi.V1AuthSignUpJSONRequestBody{
			Email:    types.Email(email),
			Password: password,
		})
		if err != nil {
			t.Fatalf("failed to auth sign up: %s", err)
		}

		if res.StatusCode != http.StatusUnauthorized {
			t.Fatalf("failed to auth sign up: %d", res.StatusCode)
		}
	})
}
