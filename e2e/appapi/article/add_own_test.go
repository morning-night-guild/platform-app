package article_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/e2e/helper"
)

func TestAppAPIE2EArticleAddOwn(t *testing.T) {
	t.Parallel()

	url := helper.GetAppAPIEndpoint(t)

	t.Run("自身の記事として追加できる", func(t *testing.T) {
		t.Parallel()

		user := helper.NewUser(t, url)

		defer user.Delete(t)

		db := helper.NewDatabase(t, helper.GetDSN(t))

		aid := uuid.New()

		db.BulkInsertArticles([]types.UUID{aid})

		defer db.Close()

		defer db.BulkDeleteArticles([]types.UUID{aid})

		client := user.Client

		res, err := client.Client.V1ArticleAddOwn(context.Background(), aid)
		if err != nil {
			t.Fatalf("failed to add own article: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("failed to add own article: %s", res.Status)
		}
	})
}
