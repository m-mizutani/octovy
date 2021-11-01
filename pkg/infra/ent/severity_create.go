// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/octovy/pkg/infra/ent/severity"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnerability"
)

// SeverityCreate is the builder for creating a Severity entity.
type SeverityCreate struct {
	config
	mutation *SeverityMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetLabel sets the "label" field.
func (sc *SeverityCreate) SetLabel(s string) *SeverityCreate {
	sc.mutation.SetLabel(s)
	return sc
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the Vulnerability entity by IDs.
func (sc *SeverityCreate) AddVulnerabilityIDs(ids ...string) *SeverityCreate {
	sc.mutation.AddVulnerabilityIDs(ids...)
	return sc
}

// AddVulnerabilities adds the "vulnerabilities" edges to the Vulnerability entity.
func (sc *SeverityCreate) AddVulnerabilities(v ...*Vulnerability) *SeverityCreate {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return sc.AddVulnerabilityIDs(ids...)
}

// Mutation returns the SeverityMutation object of the builder.
func (sc *SeverityCreate) Mutation() *SeverityMutation {
	return sc.mutation
}

// Save creates the Severity in the database.
func (sc *SeverityCreate) Save(ctx context.Context) (*Severity, error) {
	var (
		err  error
		node *Severity
	)
	if len(sc.hooks) == 0 {
		if err = sc.check(); err != nil {
			return nil, err
		}
		node, err = sc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SeverityMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = sc.check(); err != nil {
				return nil, err
			}
			sc.mutation = mutation
			if node, err = sc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(sc.hooks) - 1; i >= 0; i-- {
			if sc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = sc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, sc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SeverityCreate) SaveX(ctx context.Context) *Severity {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SeverityCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SeverityCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *SeverityCreate) check() error {
	if _, ok := sc.mutation.Label(); !ok {
		return &ValidationError{Name: "label", err: errors.New(`ent: missing required field "label"`)}
	}
	if v, ok := sc.mutation.Label(); ok {
		if err := severity.LabelValidator(v); err != nil {
			return &ValidationError{Name: "label", err: fmt.Errorf(`ent: validator failed for field "label": %w`, err)}
		}
	}
	return nil
}

func (sc *SeverityCreate) sqlSave(ctx context.Context) (*Severity, error) {
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (sc *SeverityCreate) createSpec() (*Severity, *sqlgraph.CreateSpec) {
	var (
		_node = &Severity{config: sc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: severity.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: severity.FieldID,
			},
		}
	)
	_spec.OnConflict = sc.conflict
	if value, ok := sc.mutation.Label(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: severity.FieldLabel,
		})
		_node.Label = value
	}
	if nodes := sc.mutation.VulnerabilitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   severity.VulnerabilitiesTable,
			Columns: []string{severity.VulnerabilitiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: vulnerability.FieldID,
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
//	client.Severity.Create().
//		SetLabel(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SeverityUpsert) {
//			SetLabel(v+v).
//		}).
//		Exec(ctx)
//
func (sc *SeverityCreate) OnConflict(opts ...sql.ConflictOption) *SeverityUpsertOne {
	sc.conflict = opts
	return &SeverityUpsertOne{
		create: sc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//  client.Severity.Create().
//      OnConflict(sql.ConflictColumns(columns...)).
//      Exec(ctx)
//
func (sc *SeverityCreate) OnConflictColumns(columns ...string) *SeverityUpsertOne {
	sc.conflict = append(sc.conflict, sql.ConflictColumns(columns...))
	return &SeverityUpsertOne{
		create: sc,
	}
}

type (
	// SeverityUpsertOne is the builder for "upsert"-ing
	//  one Severity node.
	SeverityUpsertOne struct {
		create *SeverityCreate
	}

	// SeverityUpsert is the "OnConflict" setter.
	SeverityUpsert struct {
		*sql.UpdateSet
	}
)

// SetLabel sets the "label" field.
func (u *SeverityUpsert) SetLabel(v string) *SeverityUpsert {
	u.Set(severity.FieldLabel, v)
	return u
}

// UpdateLabel sets the "label" field to the value that was provided on create.
func (u *SeverityUpsert) UpdateLabel() *SeverityUpsert {
	u.SetExcluded(severity.FieldLabel)
	return u
}

// UpdateNewValues updates the fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//  client.Severity.Create().
//      OnConflict(sql.ResolveWithNewValues()).
//      Exec(ctx)
//
func (u *SeverityUpsertOne) UpdateNewValues() *SeverityUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//  client.Severity.Create().
//      OnConflict(sql.ResolveWithIgnore()).
//      Exec(ctx)
//
func (u *SeverityUpsertOne) Ignore() *SeverityUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SeverityUpsertOne) DoNothing() *SeverityUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SeverityCreate.OnConflict
// documentation for more info.
func (u *SeverityUpsertOne) Update(set func(*SeverityUpsert)) *SeverityUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SeverityUpsert{UpdateSet: update})
	}))
	return u
}

// SetLabel sets the "label" field.
func (u *SeverityUpsertOne) SetLabel(v string) *SeverityUpsertOne {
	return u.Update(func(s *SeverityUpsert) {
		s.SetLabel(v)
	})
}

// UpdateLabel sets the "label" field to the value that was provided on create.
func (u *SeverityUpsertOne) UpdateLabel() *SeverityUpsertOne {
	return u.Update(func(s *SeverityUpsert) {
		s.UpdateLabel()
	})
}

// Exec executes the query.
func (u *SeverityUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SeverityCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SeverityUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *SeverityUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *SeverityUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// SeverityCreateBulk is the builder for creating many Severity entities in bulk.
type SeverityCreateBulk struct {
	config
	builders []*SeverityCreate
	conflict []sql.ConflictOption
}

// Save creates the Severity entities in the database.
func (scb *SeverityCreateBulk) Save(ctx context.Context) ([]*Severity, error) {
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Severity, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SeverityMutation)
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
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = scb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
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
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *SeverityCreateBulk) SaveX(ctx context.Context) []*Severity {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SeverityCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SeverityCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Severity.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SeverityUpsert) {
//			SetLabel(v+v).
//		}).
//		Exec(ctx)
//
func (scb *SeverityCreateBulk) OnConflict(opts ...sql.ConflictOption) *SeverityUpsertBulk {
	scb.conflict = opts
	return &SeverityUpsertBulk{
		create: scb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//  client.Severity.Create().
//      OnConflict(sql.ConflictColumns(columns...)).
//      Exec(ctx)
//
func (scb *SeverityCreateBulk) OnConflictColumns(columns ...string) *SeverityUpsertBulk {
	scb.conflict = append(scb.conflict, sql.ConflictColumns(columns...))
	return &SeverityUpsertBulk{
		create: scb,
	}
}

// SeverityUpsertBulk is the builder for "upsert"-ing
// a bulk of Severity nodes.
type SeverityUpsertBulk struct {
	create *SeverityCreateBulk
}

// UpdateNewValues updates the fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//  client.Severity.Create().
//      OnConflict(sql.ResolveWithNewValues()).
//      Exec(ctx)
//
func (u *SeverityUpsertBulk) UpdateNewValues() *SeverityUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//  client.Severity.Create().
//      OnConflict(sql.ResolveWithIgnore()).
//      Exec(ctx)
//
func (u *SeverityUpsertBulk) Ignore() *SeverityUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SeverityUpsertBulk) DoNothing() *SeverityUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SeverityCreateBulk.OnConflict
// documentation for more info.
func (u *SeverityUpsertBulk) Update(set func(*SeverityUpsert)) *SeverityUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SeverityUpsert{UpdateSet: update})
	}))
	return u
}

// SetLabel sets the "label" field.
func (u *SeverityUpsertBulk) SetLabel(v string) *SeverityUpsertBulk {
	return u.Update(func(s *SeverityUpsert) {
		s.SetLabel(v)
	})
}

// UpdateLabel sets the "label" field to the value that was provided on create.
func (u *SeverityUpsertBulk) UpdateLabel() *SeverityUpsertBulk {
	return u.Update(func(s *SeverityUpsert) {
		s.UpdateLabel()
	})
}

// Exec executes the query.
func (u *SeverityUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the SeverityCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SeverityCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SeverityUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
