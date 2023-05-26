package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// UserArticle holds the schema definition for the UserArticle entity.
type UserArticle struct {
	ent.Schema
}

// Fields of the UserArticle.
func (UserArticle) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("article_id", uuid.UUID{}),
		field.UUID("user_id", uuid.UUID{}),
		field.Time("created_at").Default(time.Now().UTC),
		field.Time("updated_at").Default(time.Now().UTC).UpdateDefault(time.Now().UTC),
	}
}

// Edges of the UserArticle.
func (UserArticle) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("article", Article.Type).
			Ref("user_articles").
			Field("article_id").
			Required().
			Unique(),
		edge.From("user", User.Type).
			Ref("user_articles").
			Field("user_id").
			Required().
			Unique(),
	}
}

// Indexes of the UserArticle.
func (UserArticle) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "article_id").Unique(),
		index.Fields("user_id"),
	}
}
