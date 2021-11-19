// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/octovy/pkg/infra/ent/repolabel"
	"github.com/m-mizutani/octovy/pkg/infra/ent/repository"
	"github.com/m-mizutani/octovy/pkg/infra/ent/scan"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnstatusindex"
)

// RepositoryCreate is the builder for creating a Repository entity.
type RepositoryCreate struct {
	config
	mutation *RepositoryMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetOwner sets the "owner" field.
func (rc *RepositoryCreate) SetOwner(s string) *RepositoryCreate {
	rc.mutation.SetOwner(s)
	return rc
}

// SetName sets the "name" field.
func (rc *RepositoryCreate) SetName(s string) *RepositoryCreate {
	rc.mutation.SetName(s)
	return rc
}

// SetInstallID sets the "install_id" field.
func (rc *RepositoryCreate) SetInstallID(i int64) *RepositoryCreate {
	rc.mutation.SetInstallID(i)
	return rc
}

// SetNillableInstallID sets the "install_id" field if the given value is not nil.
func (rc *RepositoryCreate) SetNillableInstallID(i *int64) *RepositoryCreate {
	if i != nil {
		rc.SetInstallID(*i)
	}
	return rc
}

// SetURL sets the "url" field.
func (rc *RepositoryCreate) SetURL(s string) *RepositoryCreate {
	rc.mutation.SetURL(s)
	return rc
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (rc *RepositoryCreate) SetNillableURL(s *string) *RepositoryCreate {
	if s != nil {
		rc.SetURL(*s)
	}
	return rc
}

// SetAvatarURL sets the "avatar_url" field.
func (rc *RepositoryCreate) SetAvatarURL(s string) *RepositoryCreate {
	rc.mutation.SetAvatarURL(s)
	return rc
}

// SetNillableAvatarURL sets the "avatar_url" field if the given value is not nil.
func (rc *RepositoryCreate) SetNillableAvatarURL(s *string) *RepositoryCreate {
	if s != nil {
		rc.SetAvatarURL(*s)
	}
	return rc
}

// SetDefaultBranch sets the "default_branch" field.
func (rc *RepositoryCreate) SetDefaultBranch(s string) *RepositoryCreate {
	rc.mutation.SetDefaultBranch(s)
	return rc
}

// SetNillableDefaultBranch sets the "default_branch" field if the given value is not nil.
func (rc *RepositoryCreate) SetNillableDefaultBranch(s *string) *RepositoryCreate {
	if s != nil {
		rc.SetDefaultBranch(*s)
	}
	return rc
}

// AddScanIDs adds the "scan" edge to the Scan entity by IDs.
func (rc *RepositoryCreate) AddScanIDs(ids ...string) *RepositoryCreate {
	rc.mutation.AddScanIDs(ids...)
	return rc
}

// AddScan adds the "scan" edges to the Scan entity.
func (rc *RepositoryCreate) AddScan(s ...*Scan) *RepositoryCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return rc.AddScanIDs(ids...)
}

// AddMainIDs adds the "main" edge to the Scan entity by IDs.
func (rc *RepositoryCreate) AddMainIDs(ids ...string) *RepositoryCreate {
	rc.mutation.AddMainIDs(ids...)
	return rc
}

// AddMain adds the "main" edges to the Scan entity.
func (rc *RepositoryCreate) AddMain(s ...*Scan) *RepositoryCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return rc.AddMainIDs(ids...)
}

// SetLatestID sets the "latest" edge to the Scan entity by ID.
func (rc *RepositoryCreate) SetLatestID(id string) *RepositoryCreate {
	rc.mutation.SetLatestID(id)
	return rc
}

// SetNillableLatestID sets the "latest" edge to the Scan entity by ID if the given value is not nil.
func (rc *RepositoryCreate) SetNillableLatestID(id *string) *RepositoryCreate {
	if id != nil {
		rc = rc.SetLatestID(*id)
	}
	return rc
}

// SetLatest sets the "latest" edge to the Scan entity.
func (rc *RepositoryCreate) SetLatest(s *Scan) *RepositoryCreate {
	return rc.SetLatestID(s.ID)
}

// AddStatuIDs adds the "status" edge to the VulnStatusIndex entity by IDs.
func (rc *RepositoryCreate) AddStatuIDs(ids ...string) *RepositoryCreate {
	rc.mutation.AddStatuIDs(ids...)
	return rc
}

// AddStatus adds the "status" edges to the VulnStatusIndex entity.
func (rc *RepositoryCreate) AddStatus(v ...*VulnStatusIndex) *RepositoryCreate {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return rc.AddStatuIDs(ids...)
}

// AddLabelIDs adds the "labels" edge to the RepoLabel entity by IDs.
func (rc *RepositoryCreate) AddLabelIDs(ids ...int) *RepositoryCreate {
	rc.mutation.AddLabelIDs(ids...)
	return rc
}

// AddLabels adds the "labels" edges to the RepoLabel entity.
func (rc *RepositoryCreate) AddLabels(r ...*RepoLabel) *RepositoryCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rc.AddLabelIDs(ids...)
}

// Mutation returns the RepositoryMutation object of the builder.
func (rc *RepositoryCreate) Mutation() *RepositoryMutation {
	return rc.mutation
}

// Save creates the Repository in the database.
func (rc *RepositoryCreate) Save(ctx context.Context) (*Repository, error) {
	var (
		err  error
		node *Repository
	)
	if len(rc.hooks) == 0 {
		if err = rc.check(); err != nil {
			return nil, err
		}
		node, err = rc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*RepositoryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = rc.check(); err != nil {
				return nil, err
			}
			rc.mutation = mutation
			if node, err = rc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(rc.hooks) - 1; i >= 0; i-- {
			if rc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = rc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (rc *RepositoryCreate) SaveX(ctx context.Context) *Repository {
	v, err := rc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rc *RepositoryCreate) Exec(ctx context.Context) error {
	_, err := rc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rc *RepositoryCreate) ExecX(ctx context.Context) {
	if err := rc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rc *RepositoryCreate) check() error {
	if _, ok := rc.mutation.Owner(); !ok {
		return &ValidationError{Name: "owner", err: errors.New(`ent: missing required field "owner"`)}
	}
	if _, ok := rc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "name"`)}
	}
	return nil
}

func (rc *RepositoryCreate) sqlSave(ctx context.Context) (*Repository, error) {
	_node, _spec := rc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (rc *RepositoryCreate) createSpec() (*Repository, *sqlgraph.CreateSpec) {
	var (
		_node = &Repository{config: rc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: repository.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: repository.FieldID,
			},
		}
	)
	_spec.OnConflict = rc.conflict
	if value, ok := rc.mutation.Owner(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: repository.FieldOwner,
		})
		_node.Owner = value
	}
	if value, ok := rc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: repository.FieldName,
		})
		_node.Name = value
	}
	if value, ok := rc.mutation.InstallID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: repository.FieldInstallID,
		})
		_node.InstallID = value
	}
	if value, ok := rc.mutation.URL(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: repository.FieldURL,
		})
		_node.URL = value
	}
	if value, ok := rc.mutation.AvatarURL(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: repository.FieldAvatarURL,
		})
		_node.AvatarURL = &value
	}
	if value, ok := rc.mutation.DefaultBranch(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: repository.FieldDefaultBranch,
		})
		_node.DefaultBranch = &value
	}
	if nodes := rc.mutation.ScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   repository.ScanTable,
			Columns: repository.ScanPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: scan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.MainIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   repository.MainTable,
			Columns: []string{repository.MainColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: scan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.LatestIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   repository.LatestTable,
			Columns: []string{repository.LatestColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: scan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.repository_latest = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.StatusIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   repository.StatusTable,
			Columns: []string{repository.StatusColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: vulnstatusindex.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.LabelsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   repository.LabelsTable,
			Columns: repository.LabelsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: repolabel.FieldID,
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
//	client.Repository.Create().
//		SetOwner(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.RepositoryUpsert) {
//			SetOwner(v+v).
//		}).
//		Exec(ctx)
//
func (rc *RepositoryCreate) OnConflict(opts ...sql.ConflictOption) *RepositoryUpsertOne {
	rc.conflict = opts
	return &RepositoryUpsertOne{
		create: rc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Repository.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (rc *RepositoryCreate) OnConflictColumns(columns ...string) *RepositoryUpsertOne {
	rc.conflict = append(rc.conflict, sql.ConflictColumns(columns...))
	return &RepositoryUpsertOne{
		create: rc,
	}
}

type (
	// RepositoryUpsertOne is the builder for "upsert"-ing
	//  one Repository node.
	RepositoryUpsertOne struct {
		create *RepositoryCreate
	}

	// RepositoryUpsert is the "OnConflict" setter.
	RepositoryUpsert struct {
		*sql.UpdateSet
	}
)

// SetOwner sets the "owner" field.
func (u *RepositoryUpsert) SetOwner(v string) *RepositoryUpsert {
	u.Set(repository.FieldOwner, v)
	return u
}

// UpdateOwner sets the "owner" field to the value that was provided on create.
func (u *RepositoryUpsert) UpdateOwner() *RepositoryUpsert {
	u.SetExcluded(repository.FieldOwner)
	return u
}

// SetName sets the "name" field.
func (u *RepositoryUpsert) SetName(v string) *RepositoryUpsert {
	u.Set(repository.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *RepositoryUpsert) UpdateName() *RepositoryUpsert {
	u.SetExcluded(repository.FieldName)
	return u
}

// SetInstallID sets the "install_id" field.
func (u *RepositoryUpsert) SetInstallID(v int64) *RepositoryUpsert {
	u.Set(repository.FieldInstallID, v)
	return u
}

// UpdateInstallID sets the "install_id" field to the value that was provided on create.
func (u *RepositoryUpsert) UpdateInstallID() *RepositoryUpsert {
	u.SetExcluded(repository.FieldInstallID)
	return u
}

// ClearInstallID clears the value of the "install_id" field.
func (u *RepositoryUpsert) ClearInstallID() *RepositoryUpsert {
	u.SetNull(repository.FieldInstallID)
	return u
}

// SetURL sets the "url" field.
func (u *RepositoryUpsert) SetURL(v string) *RepositoryUpsert {
	u.Set(repository.FieldURL, v)
	return u
}

// UpdateURL sets the "url" field to the value that was provided on create.
func (u *RepositoryUpsert) UpdateURL() *RepositoryUpsert {
	u.SetExcluded(repository.FieldURL)
	return u
}

// ClearURL clears the value of the "url" field.
func (u *RepositoryUpsert) ClearURL() *RepositoryUpsert {
	u.SetNull(repository.FieldURL)
	return u
}

// SetAvatarURL sets the "avatar_url" field.
func (u *RepositoryUpsert) SetAvatarURL(v string) *RepositoryUpsert {
	u.Set(repository.FieldAvatarURL, v)
	return u
}

// UpdateAvatarURL sets the "avatar_url" field to the value that was provided on create.
func (u *RepositoryUpsert) UpdateAvatarURL() *RepositoryUpsert {
	u.SetExcluded(repository.FieldAvatarURL)
	return u
}

// ClearAvatarURL clears the value of the "avatar_url" field.
func (u *RepositoryUpsert) ClearAvatarURL() *RepositoryUpsert {
	u.SetNull(repository.FieldAvatarURL)
	return u
}

// SetDefaultBranch sets the "default_branch" field.
func (u *RepositoryUpsert) SetDefaultBranch(v string) *RepositoryUpsert {
	u.Set(repository.FieldDefaultBranch, v)
	return u
}

// UpdateDefaultBranch sets the "default_branch" field to the value that was provided on create.
func (u *RepositoryUpsert) UpdateDefaultBranch() *RepositoryUpsert {
	u.SetExcluded(repository.FieldDefaultBranch)
	return u
}

// ClearDefaultBranch clears the value of the "default_branch" field.
func (u *RepositoryUpsert) ClearDefaultBranch() *RepositoryUpsert {
	u.SetNull(repository.FieldDefaultBranch)
	return u
}

// UpdateNewValues updates the fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Repository.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
//
func (u *RepositoryUpsertOne) UpdateNewValues() *RepositoryUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//  client.Repository.Create().
//      OnConflict(sql.ResolveWithIgnore()).
//      Exec(ctx)
//
func (u *RepositoryUpsertOne) Ignore() *RepositoryUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *RepositoryUpsertOne) DoNothing() *RepositoryUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the RepositoryCreate.OnConflict
// documentation for more info.
func (u *RepositoryUpsertOne) Update(set func(*RepositoryUpsert)) *RepositoryUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&RepositoryUpsert{UpdateSet: update})
	}))
	return u
}

// SetOwner sets the "owner" field.
func (u *RepositoryUpsertOne) SetOwner(v string) *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.SetOwner(v)
	})
}

// UpdateOwner sets the "owner" field to the value that was provided on create.
func (u *RepositoryUpsertOne) UpdateOwner() *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.UpdateOwner()
	})
}

// SetName sets the "name" field.
func (u *RepositoryUpsertOne) SetName(v string) *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *RepositoryUpsertOne) UpdateName() *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.UpdateName()
	})
}

// SetInstallID sets the "install_id" field.
func (u *RepositoryUpsertOne) SetInstallID(v int64) *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.SetInstallID(v)
	})
}

// UpdateInstallID sets the "install_id" field to the value that was provided on create.
func (u *RepositoryUpsertOne) UpdateInstallID() *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.UpdateInstallID()
	})
}

// ClearInstallID clears the value of the "install_id" field.
func (u *RepositoryUpsertOne) ClearInstallID() *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.ClearInstallID()
	})
}

// SetURL sets the "url" field.
func (u *RepositoryUpsertOne) SetURL(v string) *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.SetURL(v)
	})
}

// UpdateURL sets the "url" field to the value that was provided on create.
func (u *RepositoryUpsertOne) UpdateURL() *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.UpdateURL()
	})
}

// ClearURL clears the value of the "url" field.
func (u *RepositoryUpsertOne) ClearURL() *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.ClearURL()
	})
}

// SetAvatarURL sets the "avatar_url" field.
func (u *RepositoryUpsertOne) SetAvatarURL(v string) *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.SetAvatarURL(v)
	})
}

// UpdateAvatarURL sets the "avatar_url" field to the value that was provided on create.
func (u *RepositoryUpsertOne) UpdateAvatarURL() *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.UpdateAvatarURL()
	})
}

// ClearAvatarURL clears the value of the "avatar_url" field.
func (u *RepositoryUpsertOne) ClearAvatarURL() *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.ClearAvatarURL()
	})
}

// SetDefaultBranch sets the "default_branch" field.
func (u *RepositoryUpsertOne) SetDefaultBranch(v string) *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.SetDefaultBranch(v)
	})
}

// UpdateDefaultBranch sets the "default_branch" field to the value that was provided on create.
func (u *RepositoryUpsertOne) UpdateDefaultBranch() *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.UpdateDefaultBranch()
	})
}

// ClearDefaultBranch clears the value of the "default_branch" field.
func (u *RepositoryUpsertOne) ClearDefaultBranch() *RepositoryUpsertOne {
	return u.Update(func(s *RepositoryUpsert) {
		s.ClearDefaultBranch()
	})
}

// Exec executes the query.
func (u *RepositoryUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for RepositoryCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *RepositoryUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *RepositoryUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *RepositoryUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// RepositoryCreateBulk is the builder for creating many Repository entities in bulk.
type RepositoryCreateBulk struct {
	config
	builders []*RepositoryCreate
	conflict []sql.ConflictOption
}

// Save creates the Repository entities in the database.
func (rcb *RepositoryCreateBulk) Save(ctx context.Context) ([]*Repository, error) {
	specs := make([]*sqlgraph.CreateSpec, len(rcb.builders))
	nodes := make([]*Repository, len(rcb.builders))
	mutators := make([]Mutator, len(rcb.builders))
	for i := range rcb.builders {
		func(i int, root context.Context) {
			builder := rcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*RepositoryMutation)
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
					_, err = mutators[i+1].Mutate(root, rcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = rcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, rcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rcb *RepositoryCreateBulk) SaveX(ctx context.Context) []*Repository {
	v, err := rcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rcb *RepositoryCreateBulk) Exec(ctx context.Context) error {
	_, err := rcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rcb *RepositoryCreateBulk) ExecX(ctx context.Context) {
	if err := rcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Repository.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.RepositoryUpsert) {
//			SetOwner(v+v).
//		}).
//		Exec(ctx)
//
func (rcb *RepositoryCreateBulk) OnConflict(opts ...sql.ConflictOption) *RepositoryUpsertBulk {
	rcb.conflict = opts
	return &RepositoryUpsertBulk{
		create: rcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Repository.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (rcb *RepositoryCreateBulk) OnConflictColumns(columns ...string) *RepositoryUpsertBulk {
	rcb.conflict = append(rcb.conflict, sql.ConflictColumns(columns...))
	return &RepositoryUpsertBulk{
		create: rcb,
	}
}

// RepositoryUpsertBulk is the builder for "upsert"-ing
// a bulk of Repository nodes.
type RepositoryUpsertBulk struct {
	create *RepositoryCreateBulk
}

// UpdateNewValues updates the fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Repository.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
//
func (u *RepositoryUpsertBulk) UpdateNewValues() *RepositoryUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Repository.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
//
func (u *RepositoryUpsertBulk) Ignore() *RepositoryUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *RepositoryUpsertBulk) DoNothing() *RepositoryUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the RepositoryCreateBulk.OnConflict
// documentation for more info.
func (u *RepositoryUpsertBulk) Update(set func(*RepositoryUpsert)) *RepositoryUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&RepositoryUpsert{UpdateSet: update})
	}))
	return u
}

// SetOwner sets the "owner" field.
func (u *RepositoryUpsertBulk) SetOwner(v string) *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.SetOwner(v)
	})
}

// UpdateOwner sets the "owner" field to the value that was provided on create.
func (u *RepositoryUpsertBulk) UpdateOwner() *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.UpdateOwner()
	})
}

// SetName sets the "name" field.
func (u *RepositoryUpsertBulk) SetName(v string) *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *RepositoryUpsertBulk) UpdateName() *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.UpdateName()
	})
}

// SetInstallID sets the "install_id" field.
func (u *RepositoryUpsertBulk) SetInstallID(v int64) *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.SetInstallID(v)
	})
}

// UpdateInstallID sets the "install_id" field to the value that was provided on create.
func (u *RepositoryUpsertBulk) UpdateInstallID() *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.UpdateInstallID()
	})
}

// ClearInstallID clears the value of the "install_id" field.
func (u *RepositoryUpsertBulk) ClearInstallID() *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.ClearInstallID()
	})
}

// SetURL sets the "url" field.
func (u *RepositoryUpsertBulk) SetURL(v string) *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.SetURL(v)
	})
}

// UpdateURL sets the "url" field to the value that was provided on create.
func (u *RepositoryUpsertBulk) UpdateURL() *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.UpdateURL()
	})
}

// ClearURL clears the value of the "url" field.
func (u *RepositoryUpsertBulk) ClearURL() *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.ClearURL()
	})
}

// SetAvatarURL sets the "avatar_url" field.
func (u *RepositoryUpsertBulk) SetAvatarURL(v string) *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.SetAvatarURL(v)
	})
}

// UpdateAvatarURL sets the "avatar_url" field to the value that was provided on create.
func (u *RepositoryUpsertBulk) UpdateAvatarURL() *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.UpdateAvatarURL()
	})
}

// ClearAvatarURL clears the value of the "avatar_url" field.
func (u *RepositoryUpsertBulk) ClearAvatarURL() *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.ClearAvatarURL()
	})
}

// SetDefaultBranch sets the "default_branch" field.
func (u *RepositoryUpsertBulk) SetDefaultBranch(v string) *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.SetDefaultBranch(v)
	})
}

// UpdateDefaultBranch sets the "default_branch" field to the value that was provided on create.
func (u *RepositoryUpsertBulk) UpdateDefaultBranch() *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.UpdateDefaultBranch()
	})
}

// ClearDefaultBranch clears the value of the "default_branch" field.
func (u *RepositoryUpsertBulk) ClearDefaultBranch() *RepositoryUpsertBulk {
	return u.Update(func(s *RepositoryUpsert) {
		s.ClearDefaultBranch()
	})
}

// Exec executes the query.
func (u *RepositoryUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the RepositoryCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for RepositoryCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *RepositoryUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
