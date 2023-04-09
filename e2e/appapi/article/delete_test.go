package article_test

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/e2e/helper"
)

func TestAppAPIE2EArticleDelete(t *testing.T) {
	t.Parallel()

	size := uint32(5)

	url := helper.GetAppAPIEndpoint(t)

	t.Run("記事が削除できる", func(t *testing.T) {
		t.Parallel()

		db := helper.NewDatabase(t, helper.GetDSN(t))

		ids := helper.NewIDs(t, int(size))

		db.BulkInsertArticles(ids)

		defer db.Close()

		defer db.BulkDeleteArticles(ids)

		client := helper.NewOpenAPIClient(t, url)

		res, err := client.Client.V1ArticleDelete(context.Background(), ids[0])
		if err != nil {
			t.Fatalf("failed to delete article: %s", err)
		}

		defer res.Body.Close()

		article := db.SelectArticleByID(ids[0])

		if article != nil {
			t.Fatalf("article is not deleted")
		}
	})
}
