package gateway

import (
	"context"
	"time"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	da "github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
	"github.com/morning-night-guild/platform-app/pkg/ent"
	"github.com/morning-night-guild/platform-app/pkg/ent/article"
	"github.com/pkg/errors"
)

var _ repository.Article = (*Article)(nil)

// Article.
type Article struct {
	rdb *RDB
}

// NewAArticle AArticleGatewayを生成するファクトリー関数.
func NewArticle(rdb *RDB) *Article {
	return &Article{
		rdb: rdb,
	}
}

// Save 記事を保存するメソッド.
func (art *Article) Save(ctx context.Context, item model.Article) error {
	id := item.ID.Value()

	now := time.Now().UTC()

	err := art.rdb.Article.Create().
		SetID(id).
		SetTitle(item.Title.String()).
		SetURL(item.URL.String()).
		SetDescription(item.Description.String()).
		SetThumbnail(item.Thumbnail.String()).
		SetCreatedAt(now).
		SetUpdatedAt(now).
		OnConflict().
		DoNothing().
		Exec(ctx)

	if err != nil && art.rdb.IsDuplicatedError(ctx, err) {
		if ea, err := art.rdb.Article.Query().Where(article.URLEQ(item.URL.String())).First(ctx); err == nil {
			id = ea.ID
		} else {
			return errors.Wrap(err, "failed to save")
		}
	} else if err != nil {
		return errors.Wrap(err, "failed to save")
	}

	if item.TagList.Len() == 0 {
		return nil
	}

	bulk := make([]*ent.ArticleTagCreate, item.TagList.Len())
	for i, tag := range item.TagList {
		bulk[i] = art.rdb.ArticleTag.Create().
			SetTag(tag.String()).
			SetArticleID(id)
	}

	if err = art.rdb.ArticleTag.CreateBulk(bulk...).
		OnConflict().
		DoNothing().
		Exec(ctx); err == nil {
		return nil
	}

	if art.rdb.IsDuplicatedError(ctx, err) {
		return nil
	}

	return errors.Wrap(err, "failed to save")
}

// FindAll 記事を取得するメソッド.
func (art *Article) FindAll(
	ctx context.Context,
	index value.Index,
	size value.Size,
) ([]model.Article, error) {
	// ent articles
	eas, err := art.rdb.Article.Query().
		WithTags().
		Order(ent.Desc(article.FieldCreatedAt)).
		Offset(index.Int()).
		Limit(size.Int()).
		All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to article query")
	}

	articles := make([]model.Article, len(eas))

	for i, ea := range eas {
		tags := make([]string, len(ea.Edges.Tags))
		for i, tag := range ea.Edges.Tags {
			tags[i] = tag.Tag
		}

		articles[i] = model.ReconstructArticle(
			ea.ID,
			ea.URL,
			ea.Title,
			ea.Description,
			ea.Thumbnail,
			tags,
		)
	}

	return articles, nil
}

// Find ID指定で記事を取得するメソッド.
func (art *Article) Find(ctx context.Context, id da.ID) (*model.Article, error) {
	ea, err := art.rdb.Article.Query().
		Where(article.IDEQ(id.Value())).
		WithTags().
		First(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find")
	}

	if ea == nil {
		return nil, nil
	}

	tags := make([]string, len(ea.Edges.Tags))
	for i, tag := range ea.Edges.Tags {
		tags[i] = tag.Tag
	}

	article := model.ReconstructArticle(
		ea.ID,
		ea.URL,
		ea.Title,
		ea.Description,
		ea.Thumbnail,
		tags,
	)

	return &article, nil
}

func (art *Article) Delete(ctx context.Context, id da.ID) error {
	err := art.rdb.Article.DeleteOneID(id.Value()).Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to delete")
	}

	return nil
}
