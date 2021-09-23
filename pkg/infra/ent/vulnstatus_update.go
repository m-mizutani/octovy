// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/ent/predicate"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnstatus"
)

// VulnStatusUpdate is the builder for updating VulnStatus entities.
type VulnStatusUpdate struct {
	config
	hooks    []Hook
	mutation *VulnStatusMutation
}

// Where appends a list predicates to the VulnStatusUpdate builder.
func (vsu *VulnStatusUpdate) Where(ps ...predicate.VulnStatus) *VulnStatusUpdate {
	vsu.mutation.Where(ps...)
	return vsu
}

// SetStatus sets the "status" field.
func (vsu *VulnStatusUpdate) SetStatus(tst types.VulnStatusType) *VulnStatusUpdate {
	vsu.mutation.SetStatus(tst)
	return vsu
}

// SetSource sets the "source" field.
func (vsu *VulnStatusUpdate) SetSource(s string) *VulnStatusUpdate {
	vsu.mutation.SetSource(s)
	return vsu
}

// SetPkgName sets the "pkg_name" field.
func (vsu *VulnStatusUpdate) SetPkgName(s string) *VulnStatusUpdate {
	vsu.mutation.SetPkgName(s)
	return vsu
}

// SetPkgType sets the "pkg_type" field.
func (vsu *VulnStatusUpdate) SetPkgType(tt types.PkgType) *VulnStatusUpdate {
	vsu.mutation.SetPkgType(tt)
	return vsu
}

// SetVulnID sets the "vuln_id" field.
func (vsu *VulnStatusUpdate) SetVulnID(s string) *VulnStatusUpdate {
	vsu.mutation.SetVulnID(s)
	return vsu
}

// SetExpiresAt sets the "expires_at" field.
func (vsu *VulnStatusUpdate) SetExpiresAt(i int64) *VulnStatusUpdate {
	vsu.mutation.ResetExpiresAt()
	vsu.mutation.SetExpiresAt(i)
	return vsu
}

// AddExpiresAt adds i to the "expires_at" field.
func (vsu *VulnStatusUpdate) AddExpiresAt(i int64) *VulnStatusUpdate {
	vsu.mutation.AddExpiresAt(i)
	return vsu
}

// SetCreatedAt sets the "created_at" field.
func (vsu *VulnStatusUpdate) SetCreatedAt(i int64) *VulnStatusUpdate {
	vsu.mutation.ResetCreatedAt()
	vsu.mutation.SetCreatedAt(i)
	return vsu
}

// AddCreatedAt adds i to the "created_at" field.
func (vsu *VulnStatusUpdate) AddCreatedAt(i int64) *VulnStatusUpdate {
	vsu.mutation.AddCreatedAt(i)
	return vsu
}

// SetComment sets the "comment" field.
func (vsu *VulnStatusUpdate) SetComment(s string) *VulnStatusUpdate {
	vsu.mutation.SetComment(s)
	return vsu
}

// Mutation returns the VulnStatusMutation object of the builder.
func (vsu *VulnStatusUpdate) Mutation() *VulnStatusMutation {
	return vsu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (vsu *VulnStatusUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(vsu.hooks) == 0 {
		if err = vsu.check(); err != nil {
			return 0, err
		}
		affected, err = vsu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*VulnStatusMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = vsu.check(); err != nil {
				return 0, err
			}
			vsu.mutation = mutation
			affected, err = vsu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(vsu.hooks) - 1; i >= 0; i-- {
			if vsu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = vsu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, vsu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (vsu *VulnStatusUpdate) SaveX(ctx context.Context) int {
	affected, err := vsu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (vsu *VulnStatusUpdate) Exec(ctx context.Context) error {
	_, err := vsu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vsu *VulnStatusUpdate) ExecX(ctx context.Context) {
	if err := vsu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (vsu *VulnStatusUpdate) check() error {
	if v, ok := vsu.mutation.Status(); ok {
		if err := vulnstatus.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf("ent: validator failed for field \"status\": %w", err)}
		}
	}
	if v, ok := vsu.mutation.PkgType(); ok {
		if err := vulnstatus.PkgTypeValidator(v); err != nil {
			return &ValidationError{Name: "pkg_type", err: fmt.Errorf("ent: validator failed for field \"pkg_type\": %w", err)}
		}
	}
	return nil
}

func (vsu *VulnStatusUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   vulnstatus.Table,
			Columns: vulnstatus.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: vulnstatus.FieldID,
			},
		},
	}
	if ps := vsu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := vsu.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: vulnstatus.FieldStatus,
		})
	}
	if value, ok := vsu.mutation.Source(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: vulnstatus.FieldSource,
		})
	}
	if value, ok := vsu.mutation.PkgName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: vulnstatus.FieldPkgName,
		})
	}
	if value, ok := vsu.mutation.PkgType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: vulnstatus.FieldPkgType,
		})
	}
	if value, ok := vsu.mutation.VulnID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: vulnstatus.FieldVulnID,
		})
	}
	if value, ok := vsu.mutation.ExpiresAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: vulnstatus.FieldExpiresAt,
		})
	}
	if value, ok := vsu.mutation.AddedExpiresAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: vulnstatus.FieldExpiresAt,
		})
	}
	if value, ok := vsu.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: vulnstatus.FieldCreatedAt,
		})
	}
	if value, ok := vsu.mutation.AddedCreatedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: vulnstatus.FieldCreatedAt,
		})
	}
	if value, ok := vsu.mutation.Comment(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: vulnstatus.FieldComment,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, vsu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{vulnstatus.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// VulnStatusUpdateOne is the builder for updating a single VulnStatus entity.
type VulnStatusUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *VulnStatusMutation
}

// SetStatus sets the "status" field.
func (vsuo *VulnStatusUpdateOne) SetStatus(tst types.VulnStatusType) *VulnStatusUpdateOne {
	vsuo.mutation.SetStatus(tst)
	return vsuo
}

// SetSource sets the "source" field.
func (vsuo *VulnStatusUpdateOne) SetSource(s string) *VulnStatusUpdateOne {
	vsuo.mutation.SetSource(s)
	return vsuo
}

// SetPkgName sets the "pkg_name" field.
func (vsuo *VulnStatusUpdateOne) SetPkgName(s string) *VulnStatusUpdateOne {
	vsuo.mutation.SetPkgName(s)
	return vsuo
}

// SetPkgType sets the "pkg_type" field.
func (vsuo *VulnStatusUpdateOne) SetPkgType(tt types.PkgType) *VulnStatusUpdateOne {
	vsuo.mutation.SetPkgType(tt)
	return vsuo
}

// SetVulnID sets the "vuln_id" field.
func (vsuo *VulnStatusUpdateOne) SetVulnID(s string) *VulnStatusUpdateOne {
	vsuo.mutation.SetVulnID(s)
	return vsuo
}

// SetExpiresAt sets the "expires_at" field.
func (vsuo *VulnStatusUpdateOne) SetExpiresAt(i int64) *VulnStatusUpdateOne {
	vsuo.mutation.ResetExpiresAt()
	vsuo.mutation.SetExpiresAt(i)
	return vsuo
}

// AddExpiresAt adds i to the "expires_at" field.
func (vsuo *VulnStatusUpdateOne) AddExpiresAt(i int64) *VulnStatusUpdateOne {
	vsuo.mutation.AddExpiresAt(i)
	return vsuo
}

// SetCreatedAt sets the "created_at" field.
func (vsuo *VulnStatusUpdateOne) SetCreatedAt(i int64) *VulnStatusUpdateOne {
	vsuo.mutation.ResetCreatedAt()
	vsuo.mutation.SetCreatedAt(i)
	return vsuo
}

// AddCreatedAt adds i to the "created_at" field.
func (vsuo *VulnStatusUpdateOne) AddCreatedAt(i int64) *VulnStatusUpdateOne {
	vsuo.mutation.AddCreatedAt(i)
	return vsuo
}

// SetComment sets the "comment" field.
func (vsuo *VulnStatusUpdateOne) SetComment(s string) *VulnStatusUpdateOne {
	vsuo.mutation.SetComment(s)
	return vsuo
}

// Mutation returns the VulnStatusMutation object of the builder.
func (vsuo *VulnStatusUpdateOne) Mutation() *VulnStatusMutation {
	return vsuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (vsuo *VulnStatusUpdateOne) Select(field string, fields ...string) *VulnStatusUpdateOne {
	vsuo.fields = append([]string{field}, fields...)
	return vsuo
}

// Save executes the query and returns the updated VulnStatus entity.
func (vsuo *VulnStatusUpdateOne) Save(ctx context.Context) (*VulnStatus, error) {
	var (
		err  error
		node *VulnStatus
	)
	if len(vsuo.hooks) == 0 {
		if err = vsuo.check(); err != nil {
			return nil, err
		}
		node, err = vsuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*VulnStatusMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = vsuo.check(); err != nil {
				return nil, err
			}
			vsuo.mutation = mutation
			node, err = vsuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(vsuo.hooks) - 1; i >= 0; i-- {
			if vsuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = vsuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, vsuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (vsuo *VulnStatusUpdateOne) SaveX(ctx context.Context) *VulnStatus {
	node, err := vsuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (vsuo *VulnStatusUpdateOne) Exec(ctx context.Context) error {
	_, err := vsuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vsuo *VulnStatusUpdateOne) ExecX(ctx context.Context) {
	if err := vsuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (vsuo *VulnStatusUpdateOne) check() error {
	if v, ok := vsuo.mutation.Status(); ok {
		if err := vulnstatus.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf("ent: validator failed for field \"status\": %w", err)}
		}
	}
	if v, ok := vsuo.mutation.PkgType(); ok {
		if err := vulnstatus.PkgTypeValidator(v); err != nil {
			return &ValidationError{Name: "pkg_type", err: fmt.Errorf("ent: validator failed for field \"pkg_type\": %w", err)}
		}
	}
	return nil
}

func (vsuo *VulnStatusUpdateOne) sqlSave(ctx context.Context) (_node *VulnStatus, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   vulnstatus.Table,
			Columns: vulnstatus.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: vulnstatus.FieldID,
			},
		},
	}
	id, ok := vsuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing VulnStatus.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := vsuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, vulnstatus.FieldID)
		for _, f := range fields {
			if !vulnstatus.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != vulnstatus.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := vsuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := vsuo.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: vulnstatus.FieldStatus,
		})
	}
	if value, ok := vsuo.mutation.Source(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: vulnstatus.FieldSource,
		})
	}
	if value, ok := vsuo.mutation.PkgName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: vulnstatus.FieldPkgName,
		})
	}
	if value, ok := vsuo.mutation.PkgType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: vulnstatus.FieldPkgType,
		})
	}
	if value, ok := vsuo.mutation.VulnID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: vulnstatus.FieldVulnID,
		})
	}
	if value, ok := vsuo.mutation.ExpiresAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: vulnstatus.FieldExpiresAt,
		})
	}
	if value, ok := vsuo.mutation.AddedExpiresAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: vulnstatus.FieldExpiresAt,
		})
	}
	if value, ok := vsuo.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: vulnstatus.FieldCreatedAt,
		})
	}
	if value, ok := vsuo.mutation.AddedCreatedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: vulnstatus.FieldCreatedAt,
		})
	}
	if value, ok := vsuo.mutation.Comment(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: vulnstatus.FieldComment,
		})
	}
	_node = &VulnStatus{config: vsuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, vsuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{vulnstatus.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}