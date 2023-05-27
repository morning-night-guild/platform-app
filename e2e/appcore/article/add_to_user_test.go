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

func TestAppCoreE2EArticleAddToUser(t *testing.T) {
	t.Parallel()

	url := helper.GetAppCoreEndpoint(t)

	t.Run("ユーザーに記事が追加できる", func(t *testing.T) {
		t.Parallel()

		client := helper.NewConnectClient(t, &http.Client{}, url)

		db := helper.NewDatabase(t, helper.GetDSN(t))

		defer db.Close()

		aid := uuid.New()

		uid := uuid.New()

		db.InsertUser(uid)

		db.BulkInsertArticles([]uuid.UUID{aid})

		defer db.BulkDeleteArticles([]uuid.UUID{})

		defer db.DeleteUser(uid)

		req := &articlev1.AddToUserRequest{
			UserId:    uid.String(),
			ArticleId: aid.String(),
		}

		if _, err := client.Article.AddToUser(context.Background(), connect.NewRequest(req)); err != nil {
			t.Errorf("failed to add article: %s", err)
		}
	})

	t.Run("ユーザーが存在せず記事が追加できない", func(t *testing.T) {
		t.Parallel()

		client := helper.NewConnectClient(t, &http.Client{}, url)

		db := helper.NewDatabase(t, helper.GetDSN(t))

		defer db.Close()

		aid := uuid.New()

		db.BulkInsertArticles([]uuid.UUID{aid})

		// NOTE: 記事またはユーザーを削除すればユーザーに紐づく記事も削除される

		defer db.BulkDeleteArticles([]uuid.UUID{})

		req := &articlev1.AddToUserRequest{
			UserId:    uuid.New().String(),
			ArticleId: aid.String(),
		}

		if _, err := client.Article.AddToUser(context.Background(), connect.NewRequest(req)); err == nil {
			t.Error("err is nil")
		}
	})

	t.Run("記事が存在せず記事が追加できない", func(t *testing.T) {
		t.Parallel()

		client := helper.NewConnectClient(t, &http.Client{}, url)

		db := helper.NewDatabase(t, helper.GetDSN(t))

		defer db.Close()

		uid := uuid.New()

		db.InsertUser(uid)

		defer db.BulkDeleteArticles([]uuid.UUID{})

		defer db.DeleteUser(uid)

		req := &articlev1.AddToUserRequest{
			UserId:    uid.String(),
			ArticleId: uuid.New().String(),
		}

		if _, err := client.Article.AddToUser(context.Background(), connect.NewRequest(req)); err == nil {
			t.Error("err is nil")
		}
	})
}
