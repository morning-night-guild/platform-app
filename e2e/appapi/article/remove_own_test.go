package article_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/e2e/helper"
)

func TestAppAPIE2EArticleRemoveOwn(t *testing.T) {
	t.Parallel()

	url := helper.GetAppAPIEndpoint(t)

	t.Run("自身の記事を削除できる", func(t *testing.T) {
		t.Parallel()

		user := helper.NewUser(t, url)

		defer user.Delete(t)

		db := helper.NewDatabase(t, helper.GetDSN(t))

		aid := uuid.New()

		db.BulkInsertArticles([]types.UUID{aid})

		defer db.Close()

		defer db.BulkDeleteArticles([]types.UUID{aid})

		client := user.Client

		if res, err := client.Client.V1ArticleAddOwn(context.Background(), aid); err != nil {
			defer res.Body.Close()
			t.Fatalf("failed to add own article: %s", err)
		} else if res.StatusCode != http.StatusOK {
			t.Fatalf("failed to add own article: %s", res.Status)
		}

		res, err := client.Client.V1ArticleRemoveOwn(context.Background(), aid)
		if err != nil {
			t.Errorf("failed to remove own article: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("failed to remove own article: %s", res.Status)
		}
	})

	t.Run("存在しない記事は削除できない", func(t *testing.T) {
		t.Parallel()

		user := helper.NewUser(t, url)

		defer user.Delete(t)

		aid := uuid.New()

		client := user.Client

		res, err := client.Client.V1ArticleRemoveOwn(context.Background(), aid)
		if err != nil {
			t.Errorf("failed to remove own article: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("succeed to remove own article: %s", res.Status)
		}
	})

	t.Run("認証に失敗して記事が追加できない", func(t *testing.T) {
		t.Parallel()

		db := helper.NewDatabase(t, helper.GetDSN(t))

		aid := uuid.New()

		db.BulkInsertArticles([]types.UUID{aid})

		defer db.Close()

		defer db.BulkDeleteArticles([]types.UUID{aid})

		client := helper.NewOpenAPIClient(t, url)

		res, err := client.Client.V1ArticleRemoveOwn(context.Background(), aid)
		if err != nil {
			t.Errorf("failed to remove own article: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusUnauthorized {
			t.Fatalf("succeed to remove own article: %s", res.Status)
		}
	})
}
