// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnstatus"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnstatusindex"
)

// VulnStatusIndexCreate is the builder for creating a VulnStatusIndex entity.
type VulnStatusIndexCreate struct {
	config
	mutation *VulnStatusIndexMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetID sets the "id" field.
func (vsic *VulnStatusIndexCreate) SetID(s string) *VulnStatusIndexCreate {
	vsic.mutation.SetID(s)
	return vsic
}

// SetLatestID sets the "latest" edge to the VulnStatus entity by ID.
func (vsic *VulnStatusIndexCreate) SetLatestID(id int) *VulnStatusIndexCreate {
	vsic.mutation.SetLatestID(id)
	return vsic
}

// SetNillableLatestID sets the "latest" edge to the VulnStatus entity by ID if the given value is not nil.
func (vsic *VulnStatusIndexCreate) SetNillableLatestID(id *int) *VulnStatusIndexCreate {
	if id != nil {
		vsic = vsic.SetLatestID(*id)
	}
	return vsic
}

// SetLatest sets the "latest" edge to the VulnStatus entity.
func (vsic *VulnStatusIndexCreate) SetLatest(v *VulnStatus) *VulnStatusIndexCreate {
	return vsic.SetLatestID(v.ID)
}

// AddStatuIDs adds the "status" edge to the VulnStatus entity by IDs.
func (vsic *VulnStatusIndexCreate) AddStatuIDs(ids ...int) *VulnStatusIndexCreate {
	vsic.mutation.AddStatuIDs(ids...)
	return vsic
}

// AddStatus adds the "status" edges to the VulnStatus entity.
func (vsic *VulnStatusIndexCreate) AddStatus(v ...*VulnStatus) *VulnStatusIndexCreate {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return vsic.AddStatuIDs(ids...)
}

// Mutation returns the VulnStatusIndexMutation object of the builder.
func (vsic *VulnStatusIndexCreate) Mutation() *VulnStatusIndexMutation {
	return vsic.mutation
}

// Save creates the VulnStatusIndex in the database.
func (vsic *VulnStatusIndexCreate) Save(ctx context.Context) (*VulnStatusIndex, error) {
	var (
		err  error
		node *VulnStatusIndex
	)
	if len(vsic.hooks) == 0 {
		if err = vsic.check(); err != nil {
			return nil, err
		}
		node, err = vsic.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*VulnStatusIndexMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = vsic.check(); err != nil {
				return nil, err
			}
			vsic.mutation = mutation
			if node, err = vsic.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(vsic.hooks) - 1; i >= 0; i-- {
			if vsic.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = vsic.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, vsic.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (vsic *VulnStatusIndexCreate) SaveX(ctx context.Context) *VulnStatusIndex {
	v, err := vsic.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (vsic *VulnStatusIndexCreate) Exec(ctx context.Context) error {
	_, err := vsic.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vsic *VulnStatusIndexCreate) ExecX(ctx context.Context) {
	if err := vsic.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (vsic *VulnStatusIndexCreate) check() error {
	if v, ok := vsic.mutation.ID(); ok {
		if err := vulnstatusindex.IDValidator(v); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`ent: validator failed for field "id": %w`, err)}
		}
	}
	return nil
}

func (vsic *VulnStatusIndexCreate) sqlSave(ctx context.Context) (*VulnStatusIndex, error) {
	_node, _spec := vsic.createSpec()
	if err := sqlgraph.CreateNode(ctx, vsic.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}

func (vsic *VulnStatusIndexCreate) createSpec() (*VulnStatusIndex, *sqlgraph.CreateSpec) {
	var (
		_node = &VulnStatusIndex{config: vsic.config}
		_spec = &sqlgraph.CreateSpec{
			Table: vulnstatusindex.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: vulnstatusindex.FieldID,
			},
		}
	)
	_spec.OnConflict = vsic.conflict
	if id, ok := vsic.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if nodes := vsic.mutation.LatestIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   vulnstatusindex.LatestTable,
			Columns: []string{vulnstatusindex.LatestColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: vulnstatus.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.vuln_status_index_latest = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := vsic.mutation.StatusIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   vulnstatusindex.StatusTable,
			Columns: []string{vulnstatusindex.StatusColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: vulnstatus.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.VulnStatusIndex.Create().
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
//
func (vsic *VulnStatusIndexCreate) OnConflict(opts ...sql.ConflictOption) *VulnStatusIndexUpsertOne {
	vsic.conflict = opts
	return &VulnStatusIndexUpsertOne{
		create: vsic,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//  client.VulnStatusIndex.Create().
//      OnConflict(sql.ConflictColumns(columns...)).
//      Exec(ctx)
//
func (vsic *VulnStatusIndexCreate) OnConflictColumns(columns ...string) *VulnStatusIndexUpsertOne {
	vsic.conflict = append(vsic.conflict, sql.ConflictColumns(columns...))
	return &VulnStatusIndexUpsertOne{
		create: vsic,
	}
}

type (
	// VulnStatusIndexUpsertOne is the builder for "upsert"-ing
	//  one VulnStatusIndex node.
	VulnStatusIndexUpsertOne struct {
		create *VulnStatusIndexCreate
	}

	// VulnStatusIndexUpsert is the "OnConflict" setter.
	VulnStatusIndexUpsert struct {
		*sql.UpdateSet
	}
)

// UpdateNewValues updates the fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//  client.VulnStatusIndex.Create().
//      OnConflict(sql.ResolveWithNewValues()).
//      Exec(ctx)
//
func (u *VulnStatusIndexUpsertOne) UpdateNewValues() *VulnStatusIndexUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//  client.VulnStatusIndex.Create().
//      OnConflict(sql.ResolveWithIgnore()).
//      Exec(ctx)
//
func (u *VulnStatusIndexUpsertOne) Ignore() *VulnStatusIndexUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *VulnStatusIndexUpsertOne) DoNothing() *VulnStatusIndexUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the VulnStatusIndexCreate.OnConflict
// documentation for more info.
func (u *VulnStatusIndexUpsertOne) Update(set func(*VulnStatusIndexUpsert)) *VulnStatusIndexUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&VulnStatusIndexUpsert{UpdateSet: update})
	}))
	return u
}

// Exec executes the query.
func (u *VulnStatusIndexUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for VulnStatusIndexCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *VulnStatusIndexUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *VulnStatusIndexUpsertOne) ID(ctx context.Context) (id string, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *VulnStatusIndexUpsertOne) IDX(ctx context.Context) string {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// VulnStatusIndexCreateBulk is the builder for creating many VulnStatusIndex entities in bulk.
type VulnStatusIndexCreateBulk struct {
	config
	builders []*VulnStatusIndexCreate
	conflict []sql.ConflictOption
}

// Save creates the VulnStatusIndex entities in the database.
func (vsicb *VulnStatusIndexCreateBulk) Save(ctx context.Context) ([]*VulnStatusIndex, error) {
	specs := make([]*sqlgraph.CreateSpec, len(vsicb.builders))
	nodes := make([]*VulnStatusIndex, len(vsicb.builders))
	mutators := make([]Mutator, len(vsicb.builders))
	for i := range vsicb.builders {
		func(i int, root context.Context) {
			builder := vsicb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*VulnStatusIndexMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, vsicb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = vsicb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, vsicb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, vsicb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (vsicb *VulnStatusIndexCreateBulk) SaveX(ctx context.Context) []*VulnStatusIndex {
	v, err := vsicb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (vsicb *VulnStatusIndexCreateBulk) Exec(ctx context.Context) error {
	_, err := vsicb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vsicb *VulnStatusIndexCreateBulk) ExecX(ctx context.Context) {
	if err := vsicb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.VulnStatusIndex.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
//
func (vsicb *VulnStatusIndexCreateBulk) OnConflict(opts ...sql.ConflictOption) *VulnStatusIndexUpsertBulk {
	vsicb.conflict = opts
	return &VulnStatusIndexUpsertBulk{
		create: vsicb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//  client.VulnStatusIndex.Create().
//      OnConflict(sql.ConflictColumns(columns...)).
//      Exec(ctx)
//
func (vsicb *VulnStatusIndexCreateBulk) OnConflictColumns(columns ...string) *VulnStatusIndexUpsertBulk {
	vsicb.conflict = append(vsicb.conflict, sql.ConflictColumns(columns...))
	return &VulnStatusIndexUpsertBulk{
		create: vsicb,
	}
}

// VulnStatusIndexUpsertBulk is the builder for "upsert"-ing
// a bulk of VulnStatusIndex nodes.
type VulnStatusIndexUpsertBulk struct {
	create *VulnStatusIndexCreateBulk
}

// UpdateNewValues updates the fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//  client.VulnStatusIndex.Create().
//      OnConflict(sql.ResolveWithNewValues()).
//      Exec(ctx)
//
func (u *VulnStatusIndexUpsertBulk) UpdateNewValues() *VulnStatusIndexUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//  client.VulnStatusIndex.Create().
//      OnConflict(sql.ResolveWithIgnore()).
//      Exec(ctx)
//
func (u *VulnStatusIndexUpsertBulk) Ignore() *VulnStatusIndexUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *VulnStatusIndexUpsertBulk) DoNothing() *VulnStatusIndexUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the VulnStatusIndexCreateBulk.OnConflict
// documentation for more info.
func (u *VulnStatusIndexUpsertBulk) Update(set func(*VulnStatusIndexUpsert)) *VulnStatusIndexUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&VulnStatusIndexUpsert{UpdateSet: update})
	}))
	return u
}

// Exec executes the query.
func (u *VulnStatusIndexUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the VulnStatusIndexCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for VulnStatusIndexCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *VulnStatusIndexUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
