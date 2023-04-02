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

func TestE2EAuthSighOut(t *testing.T) {
	t.Parallel()

	url := helper.GetAppAPIEndpoint(t)

	t.Run("サインアウトできる", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		client := helper.NewOpenAPIClient(t, url)

		client.Client.Client = &http.Client{
			Transport: helper.NewAPIKeyTransport(t, helper.GetAPIKey(t)),
		}

		id := uuid.New().String()

		email := fmt.Sprintf("%s@example.com", id)

		password := id

		if res, err := client.Client.V1AuthSignUp(context.Background(), openapi.V1AuthSignUpJSONRequestBody{
			Email:    types.Email(email),
			Password: password,
		}); err != nil || res.StatusCode != http.StatusOK {
			defer res.Body.Close()

			t.Fatalf("failed to auth sign up: %s", err)
		} else {
			defer res.Body.Close()
		}

		prv, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}

		res, err := client.Client.V1AuthSignIn(ctx, openapi.V1AuthSignInJSONRequestBody{
			Email:     types.Email(email),
			Password:  password,
			PublicKey: helper.Public(t, prv),
		})
		if err != nil {
			t.Fatalf("failed to auth sign in: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("failed to auth sign in: %d", res.StatusCode)
		}

		client.Client.Client = &http.Client{
			Transport: helper.NewCookiesTransport(t, res.Cookies()),
		}

		res, err = client.Client.V1AuthSignOut(ctx)
		if err != nil {
			t.Fatalf("failed to auth sign out: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("failed to auth sign out: %d", res.StatusCode)
		}

		for _, cookie := range res.Cookies() {
			if cookie.Value != "" {
				t.Errorf("cookie[%s] value is not empty", cookie.Name)
			}
			if cookie.MaxAge != -1 {
				t.Errorf("cookie[%s] max age is not -1", cookie.Name)
			}
		}

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
}
