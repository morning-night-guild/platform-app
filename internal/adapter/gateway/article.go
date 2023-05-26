package gateway

import (
	"context"
	"time"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	domainerrors "github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
	"github.com/morning-night-guild/platform-app/pkg/ent"
	entarticle "github.com/morning-night-guild/platform-app/pkg/ent/article"
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
func (gtw *Article) Save(ctx context.Context, item model.Article) error {
	id := item.ArticleID.Value()

	now := time.Now().UTC()

	err := gtw.rdb.Article.Create().
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

	if err != nil && gtw.rdb.IsDuplicatedError(ctx, err) {
		if ea, err := gtw.rdb.Article.Query().Where(entarticle.URLEQ(item.URL.String())).First(ctx); err == nil {
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
		bulk[i] = gtw.rdb.ArticleTag.Create().
			SetTag(tag.String()).
			SetArticleID(id)
	}

	if err = gtw.rdb.ArticleTag.CreateBulk(bulk...).
		OnConflict().
		DoNothing().
		Exec(ctx); err == nil {
		return nil
	}

	if gtw.rdb.IsDuplicatedError(ctx, err) {
		return nil
	}

	return errors.Wrap(err, "failed to save")
}

// FindAll 記事を取得するメソッド.
func (gtw *Article) FindAll(
	ctx context.Context,
	index value.Index,
	size value.Size,
	filter ...value.Filter,
) ([]model.Article, error) {
	query := gtw.rdb.Article.Query().
		WithTags().
		Order(ent.Desc(entarticle.FieldCreatedAt)).
		Offset(index.Int()).
		Limit(size.Int())

	if len(filter) > 0 {
		for _, f := range filter {
			if f.Name == "title" {
				query = query.Where(entarticle.TitleContains(f.Value))
			}
		}
	}

	eas, err := query.All(ctx)
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
func (gtw *Article) Find(ctx context.Context, id article.ID) (model.Article, error) {
	ea, err := gtw.rdb.Article.Query().
		Where(entarticle.IDEQ(id.Value())).
		WithTags().
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return model.Article{}, domainerrors.NewNotFoundError("article not found")
		}

		return model.Article{}, errors.Wrap(err, "failed to find")
	}

	tags := make([]string, len(ea.Edges.Tags))
	for i, tag := range ea.Edges.Tags {
		tags[i] = tag.Tag
	}

	return model.ReconstructArticle(
		ea.ID,
		ea.URL,
		ea.Title,
		ea.Description,
		ea.Thumbnail,
		tags,
	), nil
}

func (gtw *Article) Delete(ctx context.Context, id article.ID) error {
	if err := gtw.rdb.Article.DeleteOneID(id.Value()).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return nil
		}

		return errors.Wrap(err, "failed to delete article")
	}

	return nil
}

func (gtw *Article) AddToUser(ctx context.Context, articleID article.ID, userID user.ID) error {
	if err := gtw.rdb.UserArticle.Create().
		SetArticleID(articleID.Value()).
		SetUserID(userID.Value()).
		OnConflict().
		DoNothing().
		Exec(ctx); err != nil {
		return errors.Wrap(err, "failed to add to user")
	}

	return nil
}
