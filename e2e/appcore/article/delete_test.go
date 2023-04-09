package article_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/e2e/helper"
	articlev1 "github.com/morning-night-guild/platform-app/pkg/connect/article/v1"
)

func TestAppCoreE2EArticleDelete(t *testing.T) {
	t.Parallel()

	url := helper.GetAppCoreEndpoint(t)

	t.Run("記事が削除できる", func(t *testing.T) {
		t.Parallel()

		db := helper.NewDatabase(t, helper.GetDSN(t))

		id := uuid.New()

		defer db.Close()

		defer db.BulkDeleteArticles([]uuid.UUID{id})

		db.BulkInsertArticles([]uuid.UUID{id})

		client := helper.NewConnectClient(t, &http.Client{}, url)

		req := &articlev1.DeleteRequest{
			ArticleId: id.String(),
		}

		_, err := client.Article.Delete(context.Background(), connect.NewRequest(req))
		if err != nil {
			t.Fatalf("failed to delete articles: %s", err)
		}
	})
}
