package article_test

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/e2e/helper"
	articlev1 "github.com/morning-night-guild/platform-app/pkg/connect/article/v1"
)

func TestAppCoreE2EArticleListByUser(t *testing.T) {
	t.Parallel()

	url := helper.GetAppCoreEndpoint(t)

	t.Run("ユーザーに紐づく記事が一覧できる", func(t *testing.T) {
		t.Parallel()

		db := helper.NewDatabase(t, helper.GetDSN(t))

		ids := helper.NewIDs(t, int(articleCount))

		defer db.Close()

		uid := uuid.New()

		db.InsertUser(uid)

		defer db.DeleteUser(uid)

		db.BulkInsertArticles(ids)

		defer db.BulkDeleteArticles(ids)

		client := helper.NewConnectClient(t, &http.Client{}, url)

		if _, err := client.Article.AddToUser(context.Background(), connect.NewRequest(&articlev1.AddToUserRequest{
			UserId:    uid.String(),
			ArticleId: ids[0].String(),
		})); err != nil {
			t.Fatalf("failed to add article: %s", err)
		}

		if _, err := client.Article.AddToUser(context.Background(), connect.NewRequest(&articlev1.AddToUserRequest{
			UserId:    uid.String(),
			ArticleId: ids[1].String(),
		})); err != nil {
			t.Fatalf("failed to add article: %s", err)
		}

		req := &articlev1.ListByUserRequest{
			UserId:      uid.String(),
			MaxPageSize: articleCount,
		}

		res, err := client.Article.ListByUser(context.Background(), connect.NewRequest(req))
		if err != nil {
			t.Fatalf("failed to list by user articles: %s", err)
		}

		if !reflect.DeepEqual(len(res.Msg.Articles), 2) {
			t.Errorf("Articles length = %v, want %v", len(res.Msg.Articles), 2)
		}

		if res.Msg.NextPageToken != "" {
			t.Error("Next Page Token is not empty")
		}
	})

	t.Run("ユーザーに紐づくタイトル部分一致で検索して記事を一覧できる", func(t *testing.T) {
		t.Parallel()

		toPointer := func(s string) *string {
			return &s
		}

		db := helper.NewDatabase(t, helper.GetDSN(t))

		defer db.Close()

		uid := uuid.New()

		db.InsertUser(uid)

		defer db.DeleteUser(uid)

		ids := helper.NewIDs(t, int(articleCount))

		db.BulkInsertArticles(ids)

		defer db.BulkDeleteArticles(ids)

		client := helper.NewConnectClient(t, &http.Client{}, url)

		if _, err := client.Article.AddToUser(context.Background(), connect.NewRequest(&articlev1.AddToUserRequest{
			UserId:    uid.String(),
			ArticleId: ids[0].String(),
		})); err != nil {
			t.Fatalf("failed to add article: %s", err)
		}

		if _, err := client.Article.AddToUser(context.Background(), connect.NewRequest(&articlev1.AddToUserRequest{
			UserId:    uid.String(),
			ArticleId: ids[1].String(),
		})); err != nil {
			t.Fatalf("failed to add article: %s", err)
		}

		req := &articlev1.ListByUserRequest{
			UserId:      uid.String(),
			MaxPageSize: articleCount,
			Title:       toPointer(ids[0].String()),
		}

		res, err := client.Article.ListByUser(context.Background(), connect.NewRequest(req))
		if err != nil {
			t.Fatalf("failed to list articles: %s", err)
		}

		if !reflect.DeepEqual(len(res.Msg.Articles), 1) {
			t.Errorf("Articles length = %v, want %v", len(res.Msg.Articles), 1)
		}

		if !reflect.DeepEqual(res.Msg.Articles[0].Title, fmt.Sprintf("title-%s", ids[0].String())) {
			t.Errorf("Articles Title = %v, want %v", res.Msg.Articles[0].Title, fmt.Sprintf("title-%s", ids[0].String()))
		}
	})

	t.Run("ユーザーが存在しないので紐づく記事が一覧できない", func(t *testing.T) {
		t.Parallel()

		db := helper.NewDatabase(t, helper.GetDSN(t))

		ids := helper.NewIDs(t, int(articleCount))

		defer db.Close()

		db.BulkInsertArticles(ids)

		defer db.BulkDeleteArticles(ids)

		client := helper.NewConnectClient(t, &http.Client{}, url)

		req := &articlev1.ListByUserRequest{
			UserId:      uuid.New().String(),
			MaxPageSize: articleCount,
		}

		if _, err := client.Article.ListByUser(context.Background(), connect.NewRequest(req)); err == nil {
			t.Error("error is nil")
		} else if connect.CodeOf(err) != connect.CodeNotFound {
			t.Errorf("err = %v, want %v", connect.CodeOf(err), connect.CodeNotFound)
		}
	})
}
