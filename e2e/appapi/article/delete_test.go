package article_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/e2e/helper"
)

func TestAppAPIE2EArticleDelete(t *testing.T) {
	t.Parallel()

	url := helper.GetAppAPIEndpoint(t)

	t.Run("記事が削除できる", func(t *testing.T) {
		t.Parallel()

		db := helper.NewDatabase(t, helper.GetDSN(t))

		id := uuid.New()

		db.BulkInsertArticles([]uuid.UUID{id})

		defer db.Close()

		defer db.BulkDeleteArticles([]uuid.UUID{id})

		client := helper.NewOpenAPIClient(t, url)

		res, err := client.Client.V1ArticleDelete(context.Background(), id)
		if err != nil {
			t.Fatalf("failed to delete article: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", res.StatusCode)
		}
	})
}