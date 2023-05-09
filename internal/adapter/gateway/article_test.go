package gateway_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/morning-night-guild/platform-app/internal/adapter/gateway"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
	"github.com/morning-night-guild/platform-app/pkg/ent"
	entarticle "github.com/morning-night-guild/platform-app/pkg/ent/article"
	"github.com/morning-night-guild/platform-app/pkg/ent/articletag"
)

func TestCoreArticleSave(t *testing.T) {
	t.Parallel()

	t.Run("記事を保存できる", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}

		articleGateway := gateway.NewArticle(rdb)

		ctx := context.Background()

		art := model.CreateArticle(
			article.URL("https://example.com"),
			article.Title("title"),
			article.Description("description"),
			article.Thumbnail("https://example.com"),
			article.TagList{},
		)

		if err := articleGateway.Save(ctx, art); err != nil {
			t.Error(err)
		}

		found, err := rdb.Article.Get(ctx, art.ID.Value())
		if err != nil {
			t.Error(err)
		}

		got, _ := model.NewArticle(
			article.ID(found.ID),
			article.URL(found.URL),
			article.Title(found.Title),
			article.Description(found.Description),
			article.Thumbnail(found.Thumbnail),
			article.TagList{},
		)

		if !reflect.DeepEqual(got, art) {
			t.Errorf("NewArticle() = %v, want %v", got, art)
		}

		// 同じURLを保存してもerrorにならないことを確認
		if err := articleGateway.Save(ctx, model.CreateArticle(
			article.URL("https://example.com"),
			article.Title("title"),
			article.Description("description"),
			article.Thumbnail("https://example.com"),
			article.TagList{},
		)); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("タグを含む記事が保存できる", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}

		articleGateway := gateway.NewArticle(rdb)

		ctx := context.Background()

		if err := articleGateway.Save(ctx, model.CreateArticle(
			article.URL("https://example.com"),
			article.Title("title"),
			article.Description("description"),
			article.Thumbnail("https://example.com"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
				article.Tag("tag3"),
				article.Tag("tag4"),
				article.Tag("tag5"),
			}),
		)); err != nil {
			t.Error(err)
		}
	})

	t.Run("既にある記事に既にあるタグを保存しようとしてもエラーにならない", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}

		articleGateway := gateway.NewArticle(rdb)

		ctx := context.Background()

		a1 := model.CreateArticle(
			article.URL("https://example.com"),
			article.Title("title"),
			article.Description("description"),
			article.Thumbnail("https://example.com"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
			}),
		)

		a2 := model.CreateArticle(
			article.URL("https://example.com"),
			article.Title("title"),
			article.Description("description"),
			article.Thumbnail("https://example.com"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
			}),
		)

		if err := articleGateway.Save(ctx, a1); err != nil {
			t.Fatal(err)
		}

		if err := articleGateway.Save(ctx, a2); err != nil {
			t.Fatal(err)
		}

		found, err := rdb.ArticleTag.Query().
			Where(articletag.ArticleIDEQ(a1.ID.Value())).
			All(ctx)
		if err != nil {
			t.Fatal(err)
		}

		var got article.TagList
		for _, item := range found {
			got = append(got, article.Tag(item.Tag))
		}

		if !reflect.DeepEqual(a1.TagList, got) {
			t.Errorf("NewArticle() = %v, want %v", got, a1.TagList)
		}
	})
}

func TestArticleList(t *testing.T) {
	t.Parallel()

	t.Run("記事を一覧できる（単数）", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}

		articleGateway := gateway.NewArticle(rdb)

		ctx := context.Background()

		item1 := model.CreateArticle(
			article.URL("https://example.com/1"),
			article.Title("title1"),
			article.Description("description"),
			article.Thumbnail("https://example.com/1"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
			}),
		)

		if err := articleGateway.Save(ctx, item1); err != nil {
			t.Fatal(err)
		}

		item2 := model.CreateArticle(
			article.URL("https://example.com/2"),
			article.Title("title2"),
			article.Description("description"),
			article.Thumbnail("https://example.com/2"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
			}),
		)

		if err := articleGateway.Save(ctx, item2); err != nil {
			t.Fatal(err)
		}

		got, err := articleGateway.FindAll(ctx, value.Index(0), value.Size(1))
		if err != nil {
			t.Fatal(err)
		}

		articles := []model.Article{item2}

		if !reflect.DeepEqual(got, articles) {
			t.Errorf("FindAll() = %v, want %v", got, articles)
		}
	})

	t.Run("オフセットを指定して記事を一覧できる（単数）", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}

		articleGateway := gateway.NewArticle(rdb)

		ctx := context.Background()

		item1 := model.CreateArticle(
			article.URL("https://example.com/1"),
			article.Title("title1"),
			article.Description("description"),
			article.Thumbnail("https://example.com/1"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
			}),
		)

		if err := articleGateway.Save(ctx, item1); err != nil {
			t.Fatal(err)
		}

		item2 := model.CreateArticle(
			article.URL("https://example.com/2"),
			article.Title("title2"),
			article.Description("description"),
			article.Thumbnail("https://example.com/2"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
			}),
		)

		if err := articleGateway.Save(ctx, item2); err != nil {
			t.Fatal(err)
		}

		got, err := articleGateway.FindAll(ctx, value.Index(1), value.Size(1))
		if err != nil {
			t.Fatal(err)
		}

		articles := []model.Article{item1}

		if !reflect.DeepEqual(got, articles) {
			t.Errorf("FindAll() = %v, want %v", got, articles)
		}
	})

	t.Run("記事を一覧できる（複数）", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}

		articleGateway := gateway.NewArticle(rdb)

		ctx := context.Background()

		item1 := model.CreateArticle(
			article.URL("https://example.com/1"),
			article.Title("title1"),
			article.Description("description"),
			article.Thumbnail("https://example.com/1"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
			}),
		)

		if err := articleGateway.Save(ctx, item1); err != nil {
			t.Fatal(err)
		}

		item2 := model.CreateArticle(
			article.URL("https://example.com/2"),
			article.Title("title2"),
			article.Description("description"),
			article.Thumbnail("https://example.com/2"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
			}),
		)

		if err := articleGateway.Save(ctx, item2); err != nil {
			t.Fatal(err)
		}

		got, err := articleGateway.FindAll(ctx, value.Index(0), value.Size(2))
		if err != nil {
			t.Fatal(err)
		}

		articles := []model.Article{item2, item1}

		if !reflect.DeepEqual(got, articles) {
			t.Errorf("FindAll() = %v, want %v", got, articles)
		}
	})

	t.Run("保存されている記事数を超えるサイズを指定して記事を一覧できる", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}

		articleGateway := gateway.NewArticle(rdb)

		ctx := context.Background()

		item := model.CreateArticle(
			article.URL("https://example.com"),
			article.Title("title"),
			article.Description("description"),
			article.Thumbnail("https://example.com"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
			}),
		)

		if err := articleGateway.Save(ctx, item); err != nil {
			t.Fatal(err)
		}

		got, err := articleGateway.FindAll(ctx, value.Index(0), value.Size(2))
		if err != nil {
			t.Fatal(err)
		}

		articles := []model.Article{item}

		if !reflect.DeepEqual(got, articles) {
			t.Errorf("FindAll() = %v, want %v", got, articles)
		}
	})

	t.Run("保存されている記事数を超えてインデックスを指定して記事を一覧できる", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}

		articleGateway := gateway.NewArticle(rdb)

		ctx := context.Background()

		item := model.CreateArticle(
			article.URL("https://example.com"),
			article.Title("title"),
			article.Description("description"),
			article.Thumbnail("https://example.com"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
			}),
		)

		if err := articleGateway.Save(ctx, item); err != nil {
			t.Fatal(err)
		}

		got, err := articleGateway.FindAll(ctx, value.Index(2), value.Size(2))
		if err != nil {
			t.Fatal(err)
		}

		articles := []model.Article{}

		if !reflect.DeepEqual(got, articles) {
			t.Errorf("FindAll() = %v, want %v", got, articles)
		}
	})

	t.Run("記事をタイトルで検索して一覧できる", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}

		articleGateway := gateway.NewArticle(rdb)

		ctx := context.Background()

		item1 := model.CreateArticle(
			article.URL("https://example.com/1"),
			article.Title("target title"),
			article.Description("description"),
			article.Thumbnail("https://example.com/1"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
			}),
		)

		if err := articleGateway.Save(ctx, item1); err != nil {
			t.Fatal(err)
		}

		item2 := model.CreateArticle(
			article.URL("https://example.com/2"),
			article.Title("title2"),
			article.Description("description"),
			article.Thumbnail("https://example.com/2"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
			}),
		)

		if err := articleGateway.Save(ctx, item2); err != nil {
			t.Fatal(err)
		}

		got, err := articleGateway.FindAll(ctx, value.Index(0), value.Size(2), value.NewFilter("title", "target"))
		if err != nil {
			t.Fatal(err)
		}

		articles := []model.Article{item1}

		if !reflect.DeepEqual(got, articles) {
			t.Errorf("FindAll() = %v, want %v", got, articles)
		}
	})
}

func TestArticleFind(t *testing.T) {
	t.Parallel()

	t.Run("記事を取得できる", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatalf("failed to NewRDBClientMock(): %v", err)
		}

		articleGateway := gateway.NewArticle(rdb)

		ctx := context.Background()

		item := model.CreateArticle(
			article.URL("https://example.com"),
			article.Title("title"),
			article.Description("description"),
			article.Thumbnail("https://example.com"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
			}),
		)

		if err := articleGateway.Save(ctx, item); err != nil {
			t.Fatalf("failed to Save(): %v", err)
		}

		got, err := articleGateway.Find(ctx, item.ID)
		if err != nil {
			t.Fatalf("failed to Find(): %v", err)
		}

		if !reflect.DeepEqual(got, item) {
			t.Errorf("Find() = %v, want %v", got, item)
		}
	})

	t.Run("存在しない記事を取得しようとするとNotFoundエラーになる", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}

		articleGateway := gateway.NewArticle(rdb)

		ctx := context.Background()

		if _, err = articleGateway.Find(ctx, article.GenerateID()); err == nil {
			t.Fatal("error is nil")
		}

		if !errors.AsNotFoundError(err) {
			t.Errorf("error is not NotFoundError. got: %v", err)
		}
	})
}

func TestArticleDelete(t *testing.T) {
	t.Parallel()

	t.Run("記事を削除できる", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatalf("failed to create rdb client. got %v", err)
		}

		articleGateway := gateway.NewArticle(rdb)

		ctx := context.Background()

		item := model.CreateArticle(
			article.URL("https://example.com"),
			article.Title("title"),
			article.Description("description"),
			article.Thumbnail("https://example.com"),
			article.TagList([]article.Tag{
				article.Tag("tag1"),
				article.Tag("tag2"),
			}),
		)

		if err := articleGateway.Save(ctx, item); err != nil {
			t.Fatalf("failed to save. got %v", err)
		}

		if err := articleGateway.Delete(ctx, item.ID); err != nil {
			t.Errorf("unexpected error while delete. got %v", err)
		}

		article, err := rdb.Article.Query().
			Where(entarticle.IDEQ(item.ID.Value())).
			WithTags().
			First(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return
			}

			t.Errorf("unexpected error while find. got %v", err)
		}

		if article != nil {
			t.Error("Delete target still exists")
		}
	})

	t.Run("存在しない記事を削除しようとしてもエラーにならない", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatalf("failed to create rdb client. got %v", err)
		}

		articleGateway := gateway.NewArticle(rdb)

		ctx := context.Background()

		if err := articleGateway.Delete(ctx, article.GenerateID()); err != nil {
			t.Errorf("unexpected error while delete. got %v", err)
		}
	})
}
