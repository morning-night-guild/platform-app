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

func TestAppCoreE2EArticleRemoveFromUser(t *testing.T) {
	t.Parallel()

	url := helper.GetAppCoreEndpoint(t)

	t.Run("ユーザーの記事が削除できる", func(t *testing.T) {
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

		if _, err := client.Article.AddToUser(context.Background(), connect.NewRequest(&articlev1.AddToUserRequest{
			UserId:    uid.String(),
			ArticleId: aid.String(),
		})); err != nil {
			t.Fatalf("failed to add article to user: %s", err)
		}

		req := &articlev1.RemoveFromUserRequest{
			UserId:    uid.String(),
			ArticleId: aid.String(),
		}

		if _, err := client.Article.RemoveFromUser(context.Background(), connect.NewRequest(req)); err != nil {
			t.Errorf("failed to remove article from user: %s", err)
		}
	})

	t.Run("ユーザーが存在せず記事が削除できない", func(t *testing.T) {
		t.Parallel()

		client := helper.NewConnectClient(t, &http.Client{}, url)

		db := helper.NewDatabase(t, helper.GetDSN(t))

		defer db.Close()

		aid := uuid.New()

		db.BulkInsertArticles([]uuid.UUID{aid})

		defer db.BulkDeleteArticles([]uuid.UUID{})

		req := &articlev1.RemoveFromUserRequest{
			UserId:    uuid.New().String(),
			ArticleId: aid.String(),
		}

		if _, err := client.Article.RemoveFromUser(context.Background(), connect.NewRequest(req)); err == nil {
			t.Error("err is nil")
		} else if connect.CodeOf(err) != connect.CodeNotFound {
			t.Errorf("err = %v, want %v", connect.CodeOf(err), connect.CodeNotFound)
		}
	})

	t.Run("記事が存在せず記事が削除できない", func(t *testing.T) {
		t.Parallel()

		client := helper.NewConnectClient(t, &http.Client{}, url)

		db := helper.NewDatabase(t, helper.GetDSN(t))

		defer db.Close()

		uid := uuid.New()

		db.InsertUser(uid)

		defer db.BulkDeleteArticles([]uuid.UUID{})

		defer db.DeleteUser(uid)

		req := &articlev1.RemoveFromUserRequest{
			UserId:    uid.String(),
			ArticleId: uuid.New().String(),
		}

		if _, err := client.Article.RemoveFromUser(context.Background(), connect.NewRequest(req)); err == nil {
			t.Error("err is nil")
		} else if connect.CodeOf(err) != connect.CodeNotFound {
			t.Errorf("err = %v, want %v", connect.CodeOf(err), connect.CodeNotFound)
		}
	})
}
