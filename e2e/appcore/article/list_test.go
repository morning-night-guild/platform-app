package article_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/e2e/helper"
	articlev1 "github.com/morning-night-guild/platform-app/pkg/connect/proto/article/v1"
)

const articleCount = uint32(5)

func TestAppCoreE2EArticleList(t *testing.T) {
	t.Parallel()

	url := helper.GetAppCoreEndpoint(t)

	t.Run("記事が一覧できる", func(t *testing.T) {
		t.Parallel()

		db := helper.NewDatabase(t, helper.GetDSN(t))

		ids := helper.GenerateIDs(t, int(articleCount))

		defer db.Close()

		defer db.BulkDeleteArticles(ids)

		db.BulkInsertArticles(ids)

		client := helper.NewConnectClient(t, &http.Client{}, url)

		req := &articlev1.ListRequest{
			MaxPageSize: articleCount,
		}

		res, err := client.Article.List(context.Background(), connect.NewRequest(req))
		if err != nil {
			t.Fatalf("failed to list articles: %s", err)
		}

		if !reflect.DeepEqual(len(res.Msg.Articles), int(articleCount)) {
			t.Errorf("Articles length = %v, want %v", len(res.Msg.Articles), articleCount)
		}
		if res.Msg.NextPageToken == "" {
			t.Error("Next Page Token is empty")
		}
	})
}
