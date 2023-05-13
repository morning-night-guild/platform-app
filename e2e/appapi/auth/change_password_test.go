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
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

func TestE2EChangePassword(t *testing.T) {
	t.Parallel()

	url := helper.GetAppAPIEndpoint(t)

	t.Run("パスワードの更新ができる", func(t *testing.T) {
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

			t.Fatalf("failed to auth sign up: %v", err)
		} else {
			defer res.Body.Close()
		}

		prv1, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("failed to generate key: %v", err)
		}

		pre, err := client.Client.V1AuthSignIn(ctx, openapi.V1AuthSignInJSONRequestBody{
			Email:     types.Email(email),
			Password:  password,
			PublicKey: helper.Public(t, prv1),
		})
		if err != nil {
			t.Fatalf("failed to auth sign in: %v", err)
		}

		defer pre.Body.Close()

		if pre.StatusCode != http.StatusOK {
			t.Fatalf("failed to auth sign in: %d", pre.StatusCode)
		}

		client.Client.Client = &http.Client{
			Transport: helper.NewCookiesTransport(t, pre.Cookies()),
		}

		oldPass := password

		newPass := uuid.New().String()

		prv2, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("failed to generate key: %v", err)
		}

		res, err := client.Client.V1AuthChangePassword(ctx, openapi.V1AuthChangePasswordRequestSchema{
			Email:       types.Email(email),
			NewPassword: newPass,
			OldPassword: oldPass,
			PublicKey:   helper.Public(t, prv2),
		})
		if err != nil {
			t.Errorf("failed to auth change password: %v", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("failed to auth change password: %d", res.StatusCode)
		}

		for _, cookie := range res.Cookies() {
			if cookie.Name != auth.AuthTokenKey && cookie.Name != auth.SessionTokenKey {
				t.Errorf("failed to auth change password: %s", cookie.Name)
			}
			if cookie.Value == "" {
				t.Errorf("failed to auth change password: %s", cookie.Value)
			}
		}

		// パスワード変更後のトークンで認証できることを確認する
		client.Client.Client = &http.Client{
			Transport: helper.NewCookiesTransport(t, res.Cookies()),
		}

		if res, err := client.Client.V1AuthVerify(ctx); err != nil {
			defer res.Body.Close()

			t.Errorf("failed to auth verify: %v", err)
		} else {
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("failed to auth verify: %d", res.StatusCode)
			}
		}

		// パスワード変更前のセッショントークンが無効であることを確認する
		client.Client.Client = &http.Client{
			Transport: helper.NewCookiesTransport(t, pre.Cookies()),
		}

		if res, err := client.Client.V1AuthVerify(ctx); err != nil {
			defer res.Body.Close()

			t.Errorf("failed to auth verify: %v", err)
		} else {
			defer res.Body.Close()

			if res.StatusCode != http.StatusUnauthorized {
				t.Errorf("succeeded to auth verify: %d", res.StatusCode)
			}
		}

		// 古いパスワードでログインできないことを確認する
		client.Client.Client = http.DefaultClient

		if res, err := client.Client.V1AuthSignIn(ctx, openapi.V1AuthSignInRequestSchema{
			Email:     types.Email(email),
			Password:  password,
			PublicKey: helper.Public(t, prv1),
		}); err != nil {
			defer res.Body.Close()

			t.Errorf("failed to auth sign in: %v", err)
		} else {
			defer res.Body.Close()

			if res.StatusCode != http.StatusBadRequest {
				t.Errorf("succeeded to auth sign in: %d", res.StatusCode)
			}
		}

		defer func() {
			prv, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				t.Fatal(err)
			}

			// 新しいパスワードでログインできることを確認する
			res, err := client.Client.V1AuthSignIn(context.Background(), openapi.V1AuthSignInRequestSchema{
				Email:     types.Email(email),
				Password:  newPass,
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

	t.Run("同じパスワードで更新することはできない", func(t *testing.T) {
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

			t.Fatalf("failed to auth sign up: %v", err)
		} else {
			defer res.Body.Close()
		}

		prv1, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("failed to generate key: %v", err)
		}

		pre, err := client.Client.V1AuthSignIn(ctx, openapi.V1AuthSignInJSONRequestBody{
			Email:     types.Email(email),
			Password:  password,
			PublicKey: helper.Public(t, prv1),
		})
		if err != nil {
			t.Fatalf("failed to auth sign in: %v", err)
		}

		defer pre.Body.Close()

		if pre.StatusCode != http.StatusOK {
			t.Fatalf("failed to auth sign in: %d", pre.StatusCode)
		}

		client.Client.Client = &http.Client{
			Transport: helper.NewCookiesTransport(t, pre.Cookies()),
		}

		prv2, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("failed to generate key: %v", err)
		}

		res, err := client.Client.V1AuthChangePassword(ctx, openapi.V1AuthChangePasswordRequestSchema{
			Email:       types.Email(email),
			NewPassword: password,
			OldPassword: password,
			PublicKey:   helper.Public(t, prv2),
		})
		if err != nil {
			t.Errorf("failed to auth change password: %v", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("failed to auth change password: %d", res.StatusCode)
		}

		defer func() {
			prv, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Client.V1AuthSignIn(context.Background(), openapi.V1AuthSignInRequestSchema{
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
