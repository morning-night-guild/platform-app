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

func TestE2EAuthVerify(t *testing.T) {
	t.Parallel()

	url := helper.GetAppAPIEndpoint(t)

	t.Run("検証できる", func(t *testing.T) {
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
			t.Fatalf("failed to auth sign up: %s", err)
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

		if res.StatusCode != http.StatusOK {
			t.Fatalf("failed to auth sign in: %d", res.StatusCode)
		}

		client.Client.Client = &http.Client{
			Transport: helper.NewCookiesTransport(t, res.Cookies()),
		}

		res, err = client.Client.V1AuthVerify(ctx)
		if err != nil {
			t.Fatalf("failed to verify in: %s", err)
		}

		if res.StatusCode != http.StatusOK {
			t.Fatalf("failed to verify in: %d", res.StatusCode)
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

	t.Run("cookie[session-token]がなくて認証できない", func(t *testing.T) {
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
			t.Fatalf("failed to auth sign up: %s", err)
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

		if res.StatusCode != http.StatusOK {
			t.Fatalf("failed to auth sign in: %d", res.StatusCode)
		}

		client.Client.Client = &http.Client{
			Transport: helper.NewOnlyIDTokenCookieTransport(t, res.Cookies()),
		}

		res, err = client.Client.V1AuthVerify(ctx)
		if err != nil {
			t.Fatalf("failed to verify in: %s", err)
		}

		if res.StatusCode == http.StatusOK {
			t.Fatalf("success to verify in: %d", res.StatusCode)
		}

		if res.StatusCode != http.StatusUnauthorized {
			t.Fatalf("failed to verify in: %d", res.StatusCode)
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

	t.Run("cookie[auth-token]がなくて認証できない", func(t *testing.T) {
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
			t.Fatalf("failed to auth sign up: %s", err)
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

		if res.StatusCode != http.StatusOK {
			t.Fatalf("failed to auth sign in: %d", res.StatusCode)
		}

		client.Client.Client = &http.Client{
			Transport: helper.NewOnlySessionTokenCookieTransport(t, res.Cookies()),
		}

		res, err = client.Client.V1AuthVerify(ctx)
		if err != nil {
			t.Fatalf("failed to verify in: %s", err)
		}

		if res.StatusCode == http.StatusOK {
			t.Fatalf("success to verify in: %d", res.StatusCode)
		}

		if res.StatusCode != http.StatusUnauthorized {
			t.Fatalf("failed to verify in: %d", res.StatusCode)
		}

		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)

		var unauthorized openapi.V1AuthVerifyUnauthorizedResponseSchema
		if err := json.Unmarshal(body, &unauthorized); err != nil {
			t.Fatalf("failed marshal response: %s caused by %s", body, err)
			return
		}

		// Codeの有無をチェックできればいいが、せっかくだから uuid かどうかもチェックしておく
		if _, err := uuid.Parse(unauthorized.Code.String()); err != nil {
			t.Errorf("failed to parse code: %s caused by %s", unauthorized.Code.String(), err)
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
