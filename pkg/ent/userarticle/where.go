// Code generated by ent, DO NOT EDIT.

package userarticle

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/pkg/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldLTE(FieldID, id))
}

// ArticleID applies equality check predicate on the "article_id" field. It's identical to ArticleIDEQ.
func ArticleID(v uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldEQ(FieldArticleID, v))
}

// UserID applies equality check predicate on the "user_id" field. It's identical to UserIDEQ.
func UserID(v uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldEQ(FieldUserID, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldEQ(FieldUpdatedAt, v))
}

// ArticleIDEQ applies the EQ predicate on the "article_id" field.
func ArticleIDEQ(v uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldEQ(FieldArticleID, v))
}

// ArticleIDNEQ applies the NEQ predicate on the "article_id" field.
func ArticleIDNEQ(v uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldNEQ(FieldArticleID, v))
}

// ArticleIDIn applies the In predicate on the "article_id" field.
func ArticleIDIn(vs ...uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldIn(FieldArticleID, vs...))
}

// ArticleIDNotIn applies the NotIn predicate on the "article_id" field.
func ArticleIDNotIn(vs ...uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldNotIn(FieldArticleID, vs...))
}

// UserIDEQ applies the EQ predicate on the "user_id" field.
func UserIDEQ(v uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldEQ(FieldUserID, v))
}

// UserIDNEQ applies the NEQ predicate on the "user_id" field.
func UserIDNEQ(v uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldNEQ(FieldUserID, v))
}

// UserIDIn applies the In predicate on the "user_id" field.
func UserIDIn(vs ...uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldIn(FieldUserID, vs...))
}

// UserIDNotIn applies the NotIn predicate on the "user_id" field.
func UserIDNotIn(vs ...uuid.UUID) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldNotIn(FieldUserID, vs...))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.UserArticle {
	return predicate.UserArticle(sql.FieldLTE(FieldUpdatedAt, v))
}

// HasArticle applies the HasEdge predicate on the "article" edge.
func HasArticle() predicate.UserArticle {
	return predicate.UserArticle(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ArticleTable, ArticleColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasArticleWith applies the HasEdge predicate on the "article" edge with a given conditions (other predicates).
func HasArticleWith(preds ...predicate.Article) predicate.UserArticle {
	return predicate.UserArticle(func(s *sql.Selector) {
		step := newArticleStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.UserArticle {
	return predicate.UserArticle(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.UserArticle {
	return predicate.UserArticle(func(s *sql.Selector) {
		step := newUserStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.UserArticle) predicate.UserArticle {
	return predicate.UserArticle(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.UserArticle) predicate.UserArticle {
	return predicate.UserArticle(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.UserArticle) predicate.UserArticle {
	return predicate.UserArticle(func(s *sql.Selector) {
		p(s.Not())
	})
}