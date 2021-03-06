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

// AuthStateCacheDelete is the builder for deleting a AuthStateCache entity.
type AuthStateCacheDelete struct {
	config
	hooks    []Hook
	mutation *AuthStateCacheMutation
}

// Where appends a list predicates to the AuthStateCacheDelete builder.
func (ascd *AuthStateCacheDelete) Where(ps ...predicate.AuthStateCache) *AuthStateCacheDelete {
	ascd.mutation.Where(ps...)
	return ascd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ascd *AuthStateCacheDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(ascd.hooks) == 0 {
		affected, err = ascd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AuthStateCacheMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ascd.mutation = mutation
			affected, err = ascd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(ascd.hooks) - 1; i >= 0; i-- {
			if ascd.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ascd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ascd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (ascd *AuthStateCacheDelete) ExecX(ctx context.Context) int {
	n, err := ascd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ascd *AuthStateCacheDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: authstatecache.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: authstatecache.FieldID,
			},
		},
	}
	if ps := ascd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, ascd.driver, _spec)
}

// AuthStateCacheDeleteOne is the builder for deleting a single AuthStateCache entity.
type AuthStateCacheDeleteOne struct {
	ascd *AuthStateCacheDelete
}

// Exec executes the deletion query.
func (ascdo *AuthStateCacheDeleteOne) Exec(ctx context.Context) error {
	n, err := ascdo.ascd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{authstatecache.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ascdo *AuthStateCacheDeleteOne) ExecX(ctx context.Context) {
	ascdo.ascd.ExecX(ctx)
}
