// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/pkg/ent/article"
	"github.com/morning-night-guild/platform-app/pkg/ent/articletag"
	"github.com/morning-night-guild/platform-app/pkg/ent/predicate"
)

// ArticleTagUpdate is the builder for updating ArticleTag entities.
type ArticleTagUpdate struct {
	config
	hooks    []Hook
	mutation *ArticleTagMutation
}

// Where appends a list predicates to the ArticleTagUpdate builder.
func (atu *ArticleTagUpdate) Where(ps ...predicate.ArticleTag) *ArticleTagUpdate {
	atu.mutation.Where(ps...)
	return atu
}

// SetTag sets the "tag" field.
func (atu *ArticleTagUpdate) SetTag(s string) *ArticleTagUpdate {
	atu.mutation.SetTag(s)
	return atu
}

// SetArticleID sets the "article_id" field.
func (atu *ArticleTagUpdate) SetArticleID(u uuid.UUID) *ArticleTagUpdate {
	atu.mutation.SetArticleID(u)
	return atu
}

// SetArticle sets the "article" edge to the Article entity.
func (atu *ArticleTagUpdate) SetArticle(a *Article) *ArticleTagUpdate {
	return atu.SetArticleID(a.ID)
}

// Mutation returns the ArticleTagMutation object of the builder.
func (atu *ArticleTagUpdate) Mutation() *ArticleTagMutation {
	return atu.mutation
}

// ClearArticle clears the "article" edge to the Article entity.
func (atu *ArticleTagUpdate) ClearArticle() *ArticleTagUpdate {
	atu.mutation.ClearArticle()
	return atu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (atu *ArticleTagUpdate) Save(ctx context.Context) (int, error) {
	return withHooks[int, ArticleTagMutation](ctx, atu.sqlSave, atu.mutation, atu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (atu *ArticleTagUpdate) SaveX(ctx context.Context) int {
	affected, err := atu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (atu *ArticleTagUpdate) Exec(ctx context.Context) error {
	_, err := atu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (atu *ArticleTagUpdate) ExecX(ctx context.Context) {
	if err := atu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (atu *ArticleTagUpdate) check() error {
	if _, ok := atu.mutation.ArticleID(); atu.mutation.ArticleCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "ArticleTag.article"`)
	}
	return nil
}

func (atu *ArticleTagUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := atu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(articletag.Table, articletag.Columns, sqlgraph.NewFieldSpec(articletag.FieldID, field.TypeUUID))
	if ps := atu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := atu.mutation.Tag(); ok {
		_spec.SetField(articletag.FieldTag, field.TypeString, value)
	}
	if atu.mutation.ArticleCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   articletag.ArticleTable,
			Columns: []string{articletag.ArticleColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: article.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := atu.mutation.ArticleIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   articletag.ArticleTable,
			Columns: []string{articletag.ArticleColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: article.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, atu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{articletag.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	atu.mutation.done = true
	return n, nil
}

// ArticleTagUpdateOne is the builder for updating a single ArticleTag entity.
type ArticleTagUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ArticleTagMutation
}

// SetTag sets the "tag" field.
func (atuo *ArticleTagUpdateOne) SetTag(s string) *ArticleTagUpdateOne {
	atuo.mutation.SetTag(s)
	return atuo
}

// SetArticleID sets the "article_id" field.
func (atuo *ArticleTagUpdateOne) SetArticleID(u uuid.UUID) *ArticleTagUpdateOne {
	atuo.mutation.SetArticleID(u)
	return atuo
}

// SetArticle sets the "article" edge to the Article entity.
func (atuo *ArticleTagUpdateOne) SetArticle(a *Article) *ArticleTagUpdateOne {
	return atuo.SetArticleID(a.ID)
}

// Mutation returns the ArticleTagMutation object of the builder.
func (atuo *ArticleTagUpdateOne) Mutation() *ArticleTagMutation {
	return atuo.mutation
}

// ClearArticle clears the "article" edge to the Article entity.
func (atuo *ArticleTagUpdateOne) ClearArticle() *ArticleTagUpdateOne {
	atuo.mutation.ClearArticle()
	return atuo
}

// Where appends a list predicates to the ArticleTagUpdate builder.
func (atuo *ArticleTagUpdateOne) Where(ps ...predicate.ArticleTag) *ArticleTagUpdateOne {
	atuo.mutation.Where(ps...)
	return atuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (atuo *ArticleTagUpdateOne) Select(field string, fields ...string) *ArticleTagUpdateOne {
	atuo.fields = append([]string{field}, fields...)
	return atuo
}

// Save executes the query and returns the updated ArticleTag entity.
func (atuo *ArticleTagUpdateOne) Save(ctx context.Context) (*ArticleTag, error) {
	return withHooks[*ArticleTag, ArticleTagMutation](ctx, atuo.sqlSave, atuo.mutation, atuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (atuo *ArticleTagUpdateOne) SaveX(ctx context.Context) *ArticleTag {
	node, err := atuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (atuo *ArticleTagUpdateOne) Exec(ctx context.Context) error {
	_, err := atuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (atuo *ArticleTagUpdateOne) ExecX(ctx context.Context) {
	if err := atuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (atuo *ArticleTagUpdateOne) check() error {
	if _, ok := atuo.mutation.ArticleID(); atuo.mutation.ArticleCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "ArticleTag.article"`)
	}
	return nil
}

func (atuo *ArticleTagUpdateOne) sqlSave(ctx context.Context) (_node *ArticleTag, err error) {
	if err := atuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(articletag.Table, articletag.Columns, sqlgraph.NewFieldSpec(articletag.FieldID, field.TypeUUID))
	id, ok := atuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "ArticleTag.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := atuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, articletag.FieldID)
		for _, f := range fields {
			if !articletag.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != articletag.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := atuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := atuo.mutation.Tag(); ok {
		_spec.SetField(articletag.FieldTag, field.TypeString, value)
	}
	if atuo.mutation.ArticleCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   articletag.ArticleTable,
			Columns: []string{articletag.ArticleColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: article.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := atuo.mutation.ArticleIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   articletag.ArticleTable,
			Columns: []string{articletag.ArticleColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: article.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &ArticleTag{config: atuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, atuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{articletag.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	atuo.mutation.done = true
	return _node, nil
}