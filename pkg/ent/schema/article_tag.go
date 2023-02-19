package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// ArticleTag holds the schema definition for the ArticleTag entity.
type ArticleTag struct {
	ent.Schema
}

// Fields of the ArticleTag.
func (ArticleTag) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("tag"),
		field.UUID("article_id", uuid.UUID{}),
	}
}

// Edges of the ArticleTag.
func (ArticleTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("article", Article.Type).
			Ref("tags").
			Unique().
			Field("article_id").
			// https://github.com/ent/ent/issues/1561
			Required(),
	}
}

// Indexes of the ArticleTag.
func (ArticleTag) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tag", "article_id").Unique(),
		index.Fields("tag"),
		index.Fields("article_id"),
	}
}
