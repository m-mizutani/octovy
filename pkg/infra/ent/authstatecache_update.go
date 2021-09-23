// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/octovy/pkg/infra/ent/authstatecache"
	"github.com/m-mizutani/octovy/pkg/infra/ent/predicate"
)

// AuthStateCacheUpdate is the builder for updating AuthStateCache entities.
type AuthStateCacheUpdate struct {
	config
	hooks    []Hook
	mutation *AuthStateCacheMutation
}

// Where appends a list predicates to the AuthStateCacheUpdate builder.
func (ascu *AuthStateCacheUpdate) Where(ps ...predicate.AuthStateCache) *AuthStateCacheUpdate {
	ascu.mutation.Where(ps...)
	return ascu
}

// SetExpiresAt sets the "expires_at" field.
func (ascu *AuthStateCacheUpdate) SetExpiresAt(i int64) *AuthStateCacheUpdate {
	ascu.mutation.ResetExpiresAt()
	ascu.mutation.SetExpiresAt(i)
	return ascu
}

// AddExpiresAt adds i to the "expires_at" field.
func (ascu *AuthStateCacheUpdate) AddExpiresAt(i int64) *AuthStateCacheUpdate {
	ascu.mutation.AddExpiresAt(i)
	return ascu
}

// Mutation returns the AuthStateCacheMutation object of the builder.
func (ascu *AuthStateCacheUpdate) Mutation() *AuthStateCacheMutation {
	return ascu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ascu *AuthStateCacheUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(ascu.hooks) == 0 {
		affected, err = ascu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AuthStateCacheMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ascu.mutation = mutation
			affected, err = ascu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(ascu.hooks) - 1; i >= 0; i-- {
			if ascu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ascu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ascu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (ascu *AuthStateCacheUpdate) SaveX(ctx context.Context) int {
	affected, err := ascu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ascu *AuthStateCacheUpdate) Exec(ctx context.Context) error {
	_, err := ascu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ascu *AuthStateCacheUpdate) ExecX(ctx context.Context) {
	if err := ascu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ascu *AuthStateCacheUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   authstatecache.Table,
			Columns: authstatecache.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: authstatecache.FieldID,
			},
		},
	}
	if ps := ascu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ascu.mutation.ExpiresAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: authstatecache.FieldExpiresAt,
		})
	}
	if value, ok := ascu.mutation.AddedExpiresAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: authstatecache.FieldExpiresAt,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ascu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{authstatecache.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// AuthStateCacheUpdateOne is the builder for updating a single AuthStateCache entity.
type AuthStateCacheUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *AuthStateCacheMutation
}

// SetExpiresAt sets the "expires_at" field.
func (ascuo *AuthStateCacheUpdateOne) SetExpiresAt(i int64) *AuthStateCacheUpdateOne {
	ascuo.mutation.ResetExpiresAt()
	ascuo.mutation.SetExpiresAt(i)
	return ascuo
}

// AddExpiresAt adds i to the "expires_at" field.
func (ascuo *AuthStateCacheUpdateOne) AddExpiresAt(i int64) *AuthStateCacheUpdateOne {
	ascuo.mutation.AddExpiresAt(i)
	return ascuo
}

// Mutation returns the AuthStateCacheMutation object of the builder.
func (ascuo *AuthStateCacheUpdateOne) Mutation() *AuthStateCacheMutation {
	return ascuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ascuo *AuthStateCacheUpdateOne) Select(field string, fields ...string) *AuthStateCacheUpdateOne {
	ascuo.fields = append([]string{field}, fields...)
	return ascuo
}

// Save executes the query and returns the updated AuthStateCache entity.
func (ascuo *AuthStateCacheUpdateOne) Save(ctx context.Context) (*AuthStateCache, error) {
	var (
		err  error
		node *AuthStateCache
	)
	if len(ascuo.hooks) == 0 {
		node, err = ascuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AuthStateCacheMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ascuo.mutation = mutation
			node, err = ascuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(ascuo.hooks) - 1; i >= 0; i-- {
			if ascuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ascuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ascuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (ascuo *AuthStateCacheUpdateOne) SaveX(ctx context.Context) *AuthStateCache {
	node, err := ascuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ascuo *AuthStateCacheUpdateOne) Exec(ctx context.Context) error {
	_, err := ascuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ascuo *AuthStateCacheUpdateOne) ExecX(ctx context.Context) {
	if err := ascuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ascuo *AuthStateCacheUpdateOne) sqlSave(ctx context.Context) (_node *AuthStateCache, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   authstatecache.Table,
			Columns: authstatecache.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: authstatecache.FieldID,
			},
		},
	}
	id, ok := ascuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing AuthStateCache.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := ascuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, authstatecache.FieldID)
		for _, f := range fields {
			if !authstatecache.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != authstatecache.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ascuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ascuo.mutation.ExpiresAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: authstatecache.FieldExpiresAt,
		})
	}
	if value, ok := ascuo.mutation.AddedExpiresAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: authstatecache.FieldExpiresAt,
		})
	}
	_node = &AuthStateCache{config: ascuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ascuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{authstatecache.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
