package article_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/morning-night-guild/platform-app/e2e/helper"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

func TestAppAPIE2EArticleList(t *testing.T) {
	t.Parallel()

	size := uint32(5)

	url := helper.GetAppAPIEndpoint(t)

	t.Run("記事が一覧できる", func(t *testing.T) {
		t.Parallel()

		user := helper.NewUser(t, url)

		defer user.Delete(t)

		db := helper.NewDatabase(t, helper.GetDSN(t))

		ids := helper.NewIDs(t, int(size))

		db.BulkInsertArticles(ids)

		defer db.Close()

		defer db.BulkDeleteArticles(ids)

		client := user.Client

		res, err := client.Client.V1ArticleList(context.Background(), &openapi.V1ArticleListParams{
			PageToken:   nil,
			MaxPageSize: helper.ToIntPointer(t, int(size)),
		})
		if err != nil {
			t.Fatalf("failed to list article: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("failed to list article: %s", res.Status)
		}

		body, _ := io.ReadAll(res.Body)

		var article openapi.V1ArticleListResponseSchema
		if err := json.Unmarshal(body, &article); err != nil {
			t.Fatalf("failed marshal response: %s caused by %s", body, err)

			return
		}

		if !reflect.DeepEqual(len(*article.Articles), int(size)) {
			t.Errorf("Articles length = %v, want %v", len(*article.Articles), size)
		}
	})

	t.Run("タイトル部分一致で検索して記事を一覧できる", func(t *testing.T) {
		t.Parallel()

		user := helper.NewUser(t, url)

		defer user.Delete(t)

		db := helper.NewDatabase(t, helper.GetDSN(t))

		ids := helper.NewIDs(t, int(size))

		db.BulkInsertArticles(ids)

		defer db.Close()

		defer db.BulkDeleteArticles(ids)

		client := user.Client

		res, err := client.Client.V1ArticleList(context.Background(), &openapi.V1ArticleListParams{
			PageToken:   nil,
			MaxPageSize: helper.ToIntPointer(t, int(size)),
			Title:       helper.ToStringPointer(t, string(ids[0].String())),
		})
		if err != nil {
			t.Fatalf("failed to list article: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("failed to list article: %s", res.Status)
		}

		body, _ := io.ReadAll(res.Body)

		var article openapi.V1ArticleListResponseSchema
		if err := json.Unmarshal(body, &article); err != nil {
			t.Fatalf("failed marshal response: %s caused by %s", body, err)

			return
		}

		if !reflect.DeepEqual(len(*article.Articles), 1) {
			t.Errorf("Articles length = %v, want %v", len(*article.Articles), size)
		}

		if reflect.DeepEqual((*article.Articles)[0].Title, fmt.Sprintf("title-%s", ids[0].String())) {
			t.Errorf("Articles title = %v, want %v", (*article.Articles)[0].Title, fmt.Sprintf("title-%s", ids[0].String()))
		}
	})
}
