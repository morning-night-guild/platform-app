package article_test

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/e2e/helper"
	articlev1 "github.com/morning-night-guild/platform-app/pkg/connect/article/v1"
)

func TestAppCoreE2EArticleShare(t *testing.T) {
	t.Parallel()

	url := helper.GetAppCoreEndpoint(t)

	t.Run("記事が共有できる", func(t *testing.T) {
		t.Parallel()

		client := helper.NewConnectClient(t, &http.Client{}, url)

		title := uuid.NewString()

		db := helper.NewDatabase(t, helper.GetDSN(t))

		defer db.Close()

		defer db.DeleteArticleByTitle(title)

		req := &articlev1.ShareRequest{
			Url:         "https://www.example.com",
			Title:       title,
			Description: "description",
			Thumbnail:   "https://www.example.com/thumbnail.jpg",
		}

		res, err := client.Article.Share(context.Background(), connect.NewRequest(req))
		if err != nil {
			t.Fatalf("failed to share article: %s", err)
		}

		if !reflect.DeepEqual(res.Msg.Article.Url, req.Url) {
			t.Errorf("Url = %v, want %v", res.Msg.Article.Url, req.Url)
		}
		if !reflect.DeepEqual(res.Msg.Article.Title, req.Title) {
			t.Errorf("Title = %v, want %v", res.Msg.Article.Title, req.Title)
		}
		if !reflect.DeepEqual(res.Msg.Article.Description, req.Description) {
			t.Errorf("Description = %v, want %v", res.Msg.Article.Description, req.Description)
		}
		if !reflect.DeepEqual(res.Msg.Article.Thumbnail, req.Thumbnail) {
			t.Errorf("Thumbnail = %v, want %v", res.Msg.Article.Thumbnail, req.Thumbnail)
		}

		got := db.SelectArticleByTitle(title)
		if !reflect.DeepEqual(got.URL, req.Url) {
			t.Errorf("Url = %v, want %v", got.URL, req.Url)
		}
		if !reflect.DeepEqual(got.Title, req.Title) {
			t.Errorf("Title = %v, want %v", got.Title, req.Title)
		}
		if !reflect.DeepEqual(got.Description, req.Description) {
			t.Errorf("Description = %v, want %v", got.Description, req.Description)
		}
		if !reflect.DeepEqual(got.Thumbnail, req.Thumbnail) {
			t.Errorf("Thumbnail = %v, want %v", got.Thumbnail, req.Thumbnail)
		}
	})

	t.Run("不正なURLが指定されて記事が共有できない", func(t *testing.T) {
		t.Parallel()

		client := helper.NewConnectClient(t, &http.Client{}, url)

		req := &articlev1.ShareRequest{
			Url:         "http://www.example.com",
			Title:       "title",
			Description: "description",
			Thumbnail:   "https://www.example.com/thumbnail.jpg",
		}

		_, err := client.Article.Share(context.Background(), connect.NewRequest(req))
		if !strings.Contains(err.Error(), "invalid_argument") {
			t.Errorf("err = %v", err)
		}
		if !strings.Contains(err.Error(), "bad request") {
			t.Errorf("err = %v", err)
		}
	})

	t.Run("不正なThumbnailが指定されて記事が共有できない", func(t *testing.T) {
		t.Parallel()

		client := helper.NewConnectClient(t, &http.Client{}, url)

		req := &articlev1.ShareRequest{
			Url:         "https://www.example.com",
			Title:       "title",
			Description: "description",
			Thumbnail:   "http://www.example.com/thumbnail.jpg",
		}

		_, err := client.Article.Share(context.Background(), connect.NewRequest(req))
		if !strings.Contains(err.Error(), "invalid_argument") {
			t.Errorf("err = %v", err)
		}
		if !strings.Contains(err.Error(), "bad request") {
			t.Errorf("err = %v", err)
		}
	})
}
