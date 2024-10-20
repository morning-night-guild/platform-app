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

		db.BulkInsertArticles([]uuid.UUID{id})

		client := helper.NewConnectClient(t, &http.Client{}, url)

		req := &articlev1.DeleteRequest{
			ArticleId: id.String(),
		}

		if _, err := client.Article.Delete(context.Background(), connect.NewRequest(req)); err != nil {
			t.Errorf("failed to delete articles: %s", err)
		}
	})

	t.Run("存在しない記事は削除できない", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()

		client := helper.NewConnectClient(t, &http.Client{}, url)

		req := &articlev1.DeleteRequest{
			ArticleId: id.String(),
		}

		if _, err := client.Article.Delete(context.Background(), connect.NewRequest(req)); err == nil {
			t.Errorf("succeeded to delete articles: %s", err)
		} else if connect.CodeOf(err) != connect.CodeNotFound {
			t.Errorf("err = %v, want %v", connect.CodeOf(err), connect.CodeNotFound)
		}
	})
}
