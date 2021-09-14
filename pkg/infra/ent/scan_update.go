// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/octovy/pkg/infra/ent/packagerecord"
	"github.com/m-mizutani/octovy/pkg/infra/ent/predicate"
	"github.com/m-mizutani/octovy/pkg/infra/ent/repository"
	"github.com/m-mizutani/octovy/pkg/infra/ent/scan"
)

// ScanUpdate is the builder for updating Scan entities.
type ScanUpdate struct {
	config
	hooks    []Hook
	mutation *ScanMutation
}

// Where appends a list predicates to the ScanUpdate builder.
func (su *ScanUpdate) Where(ps ...predicate.Scan) *ScanUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetRequestedAt sets the "requested_at" field.
func (su *ScanUpdate) SetRequestedAt(i int64) *ScanUpdate {
	su.mutation.ResetRequestedAt()
	su.mutation.SetRequestedAt(i)
	return su
}

// AddRequestedAt adds i to the "requested_at" field.
func (su *ScanUpdate) AddRequestedAt(i int64) *ScanUpdate {
	su.mutation.AddRequestedAt(i)
	return su
}

// SetScannedAt sets the "scanned_at" field.
func (su *ScanUpdate) SetScannedAt(i int64) *ScanUpdate {
	su.mutation.ResetScannedAt()
	su.mutation.SetScannedAt(i)
	return su
}

// SetNillableScannedAt sets the "scanned_at" field if the given value is not nil.
func (su *ScanUpdate) SetNillableScannedAt(i *int64) *ScanUpdate {
	if i != nil {
		su.SetScannedAt(*i)
	}
	return su
}

// AddScannedAt adds i to the "scanned_at" field.
func (su *ScanUpdate) AddScannedAt(i int64) *ScanUpdate {
	su.mutation.AddScannedAt(i)
	return su
}

// ClearScannedAt clears the value of the "scanned_at" field.
func (su *ScanUpdate) ClearScannedAt() *ScanUpdate {
	su.mutation.ClearScannedAt()
	return su
}

// SetCheckID sets the "check_id" field.
func (su *ScanUpdate) SetCheckID(i int64) *ScanUpdate {
	su.mutation.ResetCheckID()
	su.mutation.SetCheckID(i)
	return su
}

// AddCheckID adds i to the "check_id" field.
func (su *ScanUpdate) AddCheckID(i int64) *ScanUpdate {
	su.mutation.AddCheckID(i)
	return su
}

// SetPullRequestTarget sets the "pull_request_target" field.
func (su *ScanUpdate) SetPullRequestTarget(s string) *ScanUpdate {
	su.mutation.SetPullRequestTarget(s)
	return su
}

// SetNillablePullRequestTarget sets the "pull_request_target" field if the given value is not nil.
func (su *ScanUpdate) SetNillablePullRequestTarget(s *string) *ScanUpdate {
	if s != nil {
		su.SetPullRequestTarget(*s)
	}
	return su
}

// ClearPullRequestTarget clears the value of the "pull_request_target" field.
func (su *ScanUpdate) ClearPullRequestTarget() *ScanUpdate {
	su.mutation.ClearPullRequestTarget()
	return su
}

// AddRepositoryIDs adds the "repository" edge to the Repository entity by IDs.
func (su *ScanUpdate) AddRepositoryIDs(ids ...int) *ScanUpdate {
	su.mutation.AddRepositoryIDs(ids...)
	return su
}

// AddRepository adds the "repository" edges to the Repository entity.
func (su *ScanUpdate) AddRepository(r ...*Repository) *ScanUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return su.AddRepositoryIDs(ids...)
}

// AddPackageIDs adds the "packages" edge to the PackageRecord entity by IDs.
func (su *ScanUpdate) AddPackageIDs(ids ...int) *ScanUpdate {
	su.mutation.AddPackageIDs(ids...)
	return su
}

// AddPackages adds the "packages" edges to the PackageRecord entity.
func (su *ScanUpdate) AddPackages(p ...*PackageRecord) *ScanUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return su.AddPackageIDs(ids...)
}

// Mutation returns the ScanMutation object of the builder.
func (su *ScanUpdate) Mutation() *ScanMutation {
	return su.mutation
}

// ClearRepository clears all "repository" edges to the Repository entity.
func (su *ScanUpdate) ClearRepository() *ScanUpdate {
	su.mutation.ClearRepository()
	return su
}

// RemoveRepositoryIDs removes the "repository" edge to Repository entities by IDs.
func (su *ScanUpdate) RemoveRepositoryIDs(ids ...int) *ScanUpdate {
	su.mutation.RemoveRepositoryIDs(ids...)
	return su
}

// RemoveRepository removes "repository" edges to Repository entities.
func (su *ScanUpdate) RemoveRepository(r ...*Repository) *ScanUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return su.RemoveRepositoryIDs(ids...)
}

// ClearPackages clears all "packages" edges to the PackageRecord entity.
func (su *ScanUpdate) ClearPackages() *ScanUpdate {
	su.mutation.ClearPackages()
	return su
}

// RemovePackageIDs removes the "packages" edge to PackageRecord entities by IDs.
func (su *ScanUpdate) RemovePackageIDs(ids ...int) *ScanUpdate {
	su.mutation.RemovePackageIDs(ids...)
	return su
}

// RemovePackages removes "packages" edges to PackageRecord entities.
func (su *ScanUpdate) RemovePackages(p ...*PackageRecord) *ScanUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return su.RemovePackageIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *ScanUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(su.hooks) == 0 {
		affected, err = su.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			su.mutation = mutation
			affected, err = su.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(su.hooks) - 1; i >= 0; i-- {
			if su.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = su.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, su.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (su *ScanUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *ScanUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *ScanUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

func (su *ScanUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   scan.Table,
			Columns: scan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: scan.FieldID,
			},
		},
	}
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.RequestedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: scan.FieldRequestedAt,
		})
	}
	if value, ok := su.mutation.AddedRequestedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: scan.FieldRequestedAt,
		})
	}
	if value, ok := su.mutation.ScannedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: scan.FieldScannedAt,
		})
	}
	if value, ok := su.mutation.AddedScannedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: scan.FieldScannedAt,
		})
	}
	if su.mutation.ScannedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Column: scan.FieldScannedAt,
		})
	}
	if value, ok := su.mutation.CheckID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: scan.FieldCheckID,
		})
	}
	if value, ok := su.mutation.AddedCheckID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: scan.FieldCheckID,
		})
	}
	if value, ok := su.mutation.PullRequestTarget(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: scan.FieldPullRequestTarget,
		})
	}
	if su.mutation.PullRequestTargetCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: scan.FieldPullRequestTarget,
		})
	}
	if su.mutation.RepositoryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   scan.RepositoryTable,
			Columns: scan.RepositoryPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: repository.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedRepositoryIDs(); len(nodes) > 0 && !su.mutation.RepositoryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   scan.RepositoryTable,
			Columns: scan.RepositoryPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: repository.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RepositoryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   scan.RepositoryTable,
			Columns: scan.RepositoryPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: repository.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.mutation.PackagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   scan.PackagesTable,
			Columns: scan.PackagesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: packagerecord.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedPackagesIDs(); len(nodes) > 0 && !su.mutation.PackagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   scan.PackagesTable,
			Columns: scan.PackagesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: packagerecord.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.PackagesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   scan.PackagesTable,
			Columns: scan.PackagesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: packagerecord.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{scan.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// ScanUpdateOne is the builder for updating a single Scan entity.
type ScanUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ScanMutation
}

// SetRequestedAt sets the "requested_at" field.
func (suo *ScanUpdateOne) SetRequestedAt(i int64) *ScanUpdateOne {
	suo.mutation.ResetRequestedAt()
	suo.mutation.SetRequestedAt(i)
	return suo
}

// AddRequestedAt adds i to the "requested_at" field.
func (suo *ScanUpdateOne) AddRequestedAt(i int64) *ScanUpdateOne {
	suo.mutation.AddRequestedAt(i)
	return suo
}

// SetScannedAt sets the "scanned_at" field.
func (suo *ScanUpdateOne) SetScannedAt(i int64) *ScanUpdateOne {
	suo.mutation.ResetScannedAt()
	suo.mutation.SetScannedAt(i)
	return suo
}

// SetNillableScannedAt sets the "scanned_at" field if the given value is not nil.
func (suo *ScanUpdateOne) SetNillableScannedAt(i *int64) *ScanUpdateOne {
	if i != nil {
		suo.SetScannedAt(*i)
	}
	return suo
}

// AddScannedAt adds i to the "scanned_at" field.
func (suo *ScanUpdateOne) AddScannedAt(i int64) *ScanUpdateOne {
	suo.mutation.AddScannedAt(i)
	return suo
}

// ClearScannedAt clears the value of the "scanned_at" field.
func (suo *ScanUpdateOne) ClearScannedAt() *ScanUpdateOne {
	suo.mutation.ClearScannedAt()
	return suo
}

// SetCheckID sets the "check_id" field.
func (suo *ScanUpdateOne) SetCheckID(i int64) *ScanUpdateOne {
	suo.mutation.ResetCheckID()
	suo.mutation.SetCheckID(i)
	return suo
}

// AddCheckID adds i to the "check_id" field.
func (suo *ScanUpdateOne) AddCheckID(i int64) *ScanUpdateOne {
	suo.mutation.AddCheckID(i)
	return suo
}

// SetPullRequestTarget sets the "pull_request_target" field.
func (suo *ScanUpdateOne) SetPullRequestTarget(s string) *ScanUpdateOne {
	suo.mutation.SetPullRequestTarget(s)
	return suo
}

// SetNillablePullRequestTarget sets the "pull_request_target" field if the given value is not nil.
func (suo *ScanUpdateOne) SetNillablePullRequestTarget(s *string) *ScanUpdateOne {
	if s != nil {
		suo.SetPullRequestTarget(*s)
	}
	return suo
}

// ClearPullRequestTarget clears the value of the "pull_request_target" field.
func (suo *ScanUpdateOne) ClearPullRequestTarget() *ScanUpdateOne {
	suo.mutation.ClearPullRequestTarget()
	return suo
}

// AddRepositoryIDs adds the "repository" edge to the Repository entity by IDs.
func (suo *ScanUpdateOne) AddRepositoryIDs(ids ...int) *ScanUpdateOne {
	suo.mutation.AddRepositoryIDs(ids...)
	return suo
}

// AddRepository adds the "repository" edges to the Repository entity.
func (suo *ScanUpdateOne) AddRepository(r ...*Repository) *ScanUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return suo.AddRepositoryIDs(ids...)
}

// AddPackageIDs adds the "packages" edge to the PackageRecord entity by IDs.
func (suo *ScanUpdateOne) AddPackageIDs(ids ...int) *ScanUpdateOne {
	suo.mutation.AddPackageIDs(ids...)
	return suo
}

// AddPackages adds the "packages" edges to the PackageRecord entity.
func (suo *ScanUpdateOne) AddPackages(p ...*PackageRecord) *ScanUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return suo.AddPackageIDs(ids...)
}

// Mutation returns the ScanMutation object of the builder.
func (suo *ScanUpdateOne) Mutation() *ScanMutation {
	return suo.mutation
}

// ClearRepository clears all "repository" edges to the Repository entity.
func (suo *ScanUpdateOne) ClearRepository() *ScanUpdateOne {
	suo.mutation.ClearRepository()
	return suo
}

// RemoveRepositoryIDs removes the "repository" edge to Repository entities by IDs.
func (suo *ScanUpdateOne) RemoveRepositoryIDs(ids ...int) *ScanUpdateOne {
	suo.mutation.RemoveRepositoryIDs(ids...)
	return suo
}

// RemoveRepository removes "repository" edges to Repository entities.
func (suo *ScanUpdateOne) RemoveRepository(r ...*Repository) *ScanUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return suo.RemoveRepositoryIDs(ids...)
}

// ClearPackages clears all "packages" edges to the PackageRecord entity.
func (suo *ScanUpdateOne) ClearPackages() *ScanUpdateOne {
	suo.mutation.ClearPackages()
	return suo
}

// RemovePackageIDs removes the "packages" edge to PackageRecord entities by IDs.
func (suo *ScanUpdateOne) RemovePackageIDs(ids ...int) *ScanUpdateOne {
	suo.mutation.RemovePackageIDs(ids...)
	return suo
}

// RemovePackages removes "packages" edges to PackageRecord entities.
func (suo *ScanUpdateOne) RemovePackages(p ...*PackageRecord) *ScanUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return suo.RemovePackageIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *ScanUpdateOne) Select(field string, fields ...string) *ScanUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Scan entity.
func (suo *ScanUpdateOne) Save(ctx context.Context) (*Scan, error) {
	var (
		err  error
		node *Scan
	)
	if len(suo.hooks) == 0 {
		node, err = suo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ScanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			suo.mutation = mutation
			node, err = suo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(suo.hooks) - 1; i >= 0; i-- {
			if suo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = suo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, suo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (suo *ScanUpdateOne) SaveX(ctx context.Context) *Scan {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *ScanUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *ScanUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (suo *ScanUpdateOne) sqlSave(ctx context.Context) (_node *Scan, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   scan.Table,
			Columns: scan.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: scan.FieldID,
			},
		},
	}
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Scan.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, scan.FieldID)
		for _, f := range fields {
			if !scan.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != scan.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := suo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := suo.mutation.RequestedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: scan.FieldRequestedAt,
		})
	}
	if value, ok := suo.mutation.AddedRequestedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: scan.FieldRequestedAt,
		})
	}
	if value, ok := suo.mutation.ScannedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: scan.FieldScannedAt,
		})
	}
	if value, ok := suo.mutation.AddedScannedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: scan.FieldScannedAt,
		})
	}
	if suo.mutation.ScannedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Column: scan.FieldScannedAt,
		})
	}
	if value, ok := suo.mutation.CheckID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: scan.FieldCheckID,
		})
	}
	if value, ok := suo.mutation.AddedCheckID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: scan.FieldCheckID,
		})
	}
	if value, ok := suo.mutation.PullRequestTarget(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: scan.FieldPullRequestTarget,
		})
	}
	if suo.mutation.PullRequestTargetCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: scan.FieldPullRequestTarget,
		})
	}
	if suo.mutation.RepositoryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   scan.RepositoryTable,
			Columns: scan.RepositoryPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: repository.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedRepositoryIDs(); len(nodes) > 0 && !suo.mutation.RepositoryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   scan.RepositoryTable,
			Columns: scan.RepositoryPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: repository.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RepositoryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   scan.RepositoryTable,
			Columns: scan.RepositoryPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: repository.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.mutation.PackagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   scan.PackagesTable,
			Columns: scan.PackagesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: packagerecord.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedPackagesIDs(); len(nodes) > 0 && !suo.mutation.PackagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   scan.PackagesTable,
			Columns: scan.PackagesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: packagerecord.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.PackagesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   scan.PackagesTable,
			Columns: scan.PackagesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: packagerecord.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Scan{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{scan.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
