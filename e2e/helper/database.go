package helper

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	// postgres driver.
	_ "github.com/lib/pq"
	"github.com/morning-night-guild/platform-app/pkg/ent"
	"github.com/morning-night-guild/platform-app/pkg/ent/article"
	"github.com/morning-night-guild/platform-app/pkg/ent/user"
	"github.com/morning-night-guild/platform-app/pkg/ent/userarticle"
)

type Database struct {
	T      *testing.T
	client *ent.Client
}

func NewDatabase(
	t *testing.T,
	dsn string,
) *Database {
	t.Helper()

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	return &Database{
		T:      t,
		client: client,
	}
}

func (db *Database) Close() {
	if err := db.client.Close(); err != nil {
		db.T.Error(err)
	}
}

func (db *Database) BulkInsertArticles(ids []uuid.UUID) {
	db.T.Helper()

	count := len(ids)

	bulk := make([]*ent.ArticleCreate, count)

	for i, id := range ids {
		bulk[i] = db.client.Article.Create().
			SetID(id).
			SetTitle("title-" + id.String()).
			SetURL("https://example.com/" + id.String()).
			SetDescription("description").
			SetThumbnail("https://example.com/" + id.String()).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now())
	}

	if err := db.client.Article.CreateBulk(bulk...).
		OnConflict().
		UpdateNewValues().
		DoNothing().
		Exec(context.Background()); err != nil {
		db.T.Fatal(err)
	}
}

func (db *Database) SelectArticleByTitle(title string) *ent.Article {
	db.T.Helper()

	article, err := db.client.Article.Query().Where(article.TitleEQ(title)).Only(context.Background())
	if err != nil {
		db.T.Error(err)
	}

	return article
}

func (db *Database) BulkDeleteArticles(ids []uuid.UUID) {
	db.T.Helper()

	if _, err := db.client.Article.Delete().Where(article.IDIn(ids...)).Exec(context.Background()); err != nil {
		db.T.Error(err)
	}
}

func (db *Database) DeleteArticleByTitle(title string) {
	db.T.Helper()

	if _, err := db.client.Article.Delete().Where(article.TitleEQ(title)).Exec(context.Background()); err != nil {
		db.T.Error(err)
	}
}

func (db *Database) DeleteUserArticle(
	uid uuid.UUID,
	aid uuid.UUID,
) {
	db.T.Helper()

	if _, err := db.client.UserArticle.Delete().
		Where(
			userarticle.ArticleIDEQ(aid),
			userarticle.UserIDEQ(uid),
		).Exec(context.Background()); err != nil {
		db.T.Error(err)
	}
}

func (db *Database) InsertUser(uid uuid.UUID) {
	db.T.Helper()

	if err := db.client.User.Create().SetID(uid).Exec(context.Background()); err != nil {
		db.T.Fatal(err)
	}
}

func (db *Database) DeleteUser(uid uuid.UUID) {
	db.T.Helper()

	if _, err := db.client.User.Delete().Where(user.IDEQ(uid)).Exec(context.Background()); err != nil {
		db.T.Error(err)
	}
}
