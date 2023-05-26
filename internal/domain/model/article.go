package model

import (
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

const DefaultArticleSize = value.Size(20)

// Article 記事モデル.
type Article struct {
	ArticleID   article.ID          // ID
	URL         article.URL         // 記事のURL
	Title       article.Title       // タイトル
	Description article.Description // 記事の説明
	Thumbnail   article.Thumbnail   // サムネイル
	TagList     article.TagList     // タグリスト
}

// NewArticle 記事モデルのファクトリー関数.
func NewArticle(
	id article.ID,
	url article.URL,
	title article.Title,
	description article.Description,
	thumbnail article.Thumbnail,
	tags article.TagList,
) (Article, error) {
	article := Article{
		ArticleID:   id,
		Title:       title,
		URL:         url,
		Description: description,
		Thumbnail:   thumbnail,
		TagList:     tags,
	}

	if err := article.validate(); err != nil {
		return Article{}, err
	}

	return article, nil
}

// ReconstructArticle 記事モデルの再構築関数.
func ReconstructArticle(
	id uuid.UUID,
	url string,
	title string,
	description string,
	thumbnail string,
	tags []string,
) Article {
	tagList := make([]article.Tag, 0, len(tags))

	for _, tag := range tags {
		tagList = append(tagList, article.Tag(tag))
	}

	return Article{
		ArticleID:   article.ID(id),
		URL:         article.URL(url),
		Title:       article.Title(title),
		Description: article.Description(description),
		Thumbnail:   article.Thumbnail(thumbnail),
		TagList:     article.TagList(tagList),
	}
}

// validate 記事を検証するメソッド.
func (a Article) validate() error {
	return nil
}

// CreateArticle 記事モデルを新規作成する関数.
func CreateArticle(
	url article.URL,
	title article.Title,
	description article.Description,
	thumbnail article.Thumbnail,
	tags article.TagList,
) Article {
	id := article.GenerateID()

	return Article{
		ArticleID:   id,
		URL:         url,
		Title:       title,
		Description: description,
		Thumbnail:   thumbnail,
		TagList:     tags,
	}
}
