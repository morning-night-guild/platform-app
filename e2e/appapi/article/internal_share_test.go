package article_test

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/e2e/helper"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

func TestAppAPIE2EArticleInternalShare(t *testing.T) {
	t.Parallel()

	url := helper.GetAppAPIEndpoint(t)

	t.Run("記事が共有できる", func(t *testing.T) {
		t.Parallel()

		title := uuid.NewString()

		db := helper.NewDatabase(t, helper.GetDSN(t))

		defer db.Close()

		defer db.DeleteArticleByTitle(title)

		client := helper.NewOpenAPIClientWithAPIKey(t, url, helper.GetAPIKey(t))

		res, err := client.Client.V1InternalArticleShare(context.Background(), openapi.V1ArticleShareRequestSchema{
			Url:         fmt.Sprintf("https://example.com/%s", title),
			Title:       helper.ToStringPointer(t, title),
			Description: helper.ToStringPointer(t, "description"),
			Thumbnail:   helper.ToStringPointer(t, "https://example.com/thumbnail.jpg"),
		})
		if err != nil {
			t.Fatalf("failed to share article: %s", err)
		}

		defer res.Body.Close()

		if !reflect.DeepEqual(res.StatusCode, http.StatusOK) {
			t.Errorf("StatusCode = %v, want %v", res.StatusCode, http.StatusOK)
		}
	})

	t.Run("タイトルがない記事が共有できる", func(t *testing.T) {
		t.Parallel()

		title := ""

		db := helper.NewDatabase(t, helper.GetDSN(t))

		defer db.Close()

		defer db.DeleteArticleByTitle(title)

		client := helper.NewOpenAPIClientWithAPIKey(t, url, helper.GetAPIKey(t))

		res, err := client.Client.V1InternalArticleShare(context.Background(), openapi.V1ArticleShareRequestSchema{
			Url:         fmt.Sprintf("https://example.com/%s", title),
			Title:       nil,
			Description: helper.ToStringPointer(t, "description"),
			Thumbnail:   helper.ToStringPointer(t, "https://example.com/thumbnail.jpg"),
		})
		if err != nil {
			t.Fatalf("failed to share article: %s", err)
		}

		defer res.Body.Close()

		if !reflect.DeepEqual(res.StatusCode, http.StatusOK) {
			t.Errorf("StatusCode = %v, want %v", res.StatusCode, http.StatusOK)
		}
	})

	t.Run("詳細がない記事が共有できる", func(t *testing.T) {
		t.Parallel()

		title := uuid.NewString()

		db := helper.NewDatabase(t, helper.GetDSN(t))

		defer db.Close()

		defer db.DeleteArticleByTitle(title)

		client := helper.NewOpenAPIClientWithAPIKey(t, url, helper.GetAPIKey(t))

		res, err := client.Client.V1InternalArticleShare(context.Background(), openapi.V1ArticleShareRequestSchema{
			Url:         fmt.Sprintf("https://example.com/%s", title),
			Title:       helper.ToStringPointer(t, title),
			Description: nil,
			Thumbnail:   helper.ToStringPointer(t, "https://example.com/thumbnail.jpg"),
		})
		if err != nil {
			t.Fatalf("failed to share article: %s", err)
		}

		defer res.Body.Close()

		if !reflect.DeepEqual(res.StatusCode, http.StatusOK) {
			t.Errorf("StatusCode = %v, want %v", res.StatusCode, http.StatusOK)
		}
	})

	t.Run("サムネイルがない記事が共有できる", func(t *testing.T) {
		t.Parallel()

		title := uuid.NewString()

		db := helper.NewDatabase(t, helper.GetDSN(t))

		defer db.Close()

		defer db.DeleteArticleByTitle(title)

		client := helper.NewOpenAPIClientWithAPIKey(t, url, helper.GetAPIKey(t))

		res, err := client.Client.V1InternalArticleShare(context.Background(), openapi.V1ArticleShareRequestSchema{
			Url:         fmt.Sprintf("https://example.com/%s", title),
			Title:       helper.ToStringPointer(t, title),
			Description: helper.ToStringPointer(t, "description"),
			Thumbnail:   nil,
		})
		if err != nil {
			t.Fatalf("failed to share article: %s", err)
		}

		defer res.Body.Close()

		if !reflect.DeepEqual(res.StatusCode, http.StatusOK) {
			t.Errorf("StatusCode = %v, want %v", res.StatusCode, http.StatusOK)
		}
	})

	t.Run("urlが空文字で記事が共有できない", func(t *testing.T) {
		t.Parallel()

		client := helper.NewOpenAPIClientWithAPIKey(t, url, helper.GetAPIKey(t))

		res, err := client.Client.V1InternalArticleShare(context.Background(), openapi.V1ArticleShareRequestSchema{
			Url:         "",
			Title:       helper.ToStringPointer(t, "title"),
			Description: helper.ToStringPointer(t, "description"),
			Thumbnail:   helper.ToStringPointer(t, "https://example.com/thumbnail.jpg"),
		})
		if err != nil {
			t.Fatalf("failed to share article: %s", err)
		}

		defer res.Body.Close()

		if !reflect.DeepEqual(res.StatusCode, http.StatusBadRequest) {
			t.Errorf("StatusCode = %v, want %v", res.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("Api-Keyがなくて記事が共有できない", func(t *testing.T) {
		t.Parallel()

		title := uuid.NewString()

		client := helper.NewOpenAPIClient(t, url)

		res, err := client.Client.V1InternalArticleShare(context.Background(), openapi.V1ArticleShareRequestSchema{
			Url:         fmt.Sprintf("https://example.com/%s", title),
			Title:       helper.ToStringPointer(t, "title"),
			Description: helper.ToStringPointer(t, "description"),
			Thumbnail:   helper.ToStringPointer(t, "https://example.com/thumbnail.jpg"),
		})
		if err != nil {
			t.Fatalf("failed to share article: %s", err)
		}

		defer res.Body.Close()

		if !reflect.DeepEqual(res.StatusCode, http.StatusUnauthorized) {
			t.Errorf("StatusCode = %v, want %v", res.StatusCode, http.StatusUnauthorized)
		}
	})
}
