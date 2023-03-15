// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/octovy/pkg/infra/ent/predicate"
	"github.com/m-mizutani/octovy/pkg/infra/ent/severity"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnerability"
)

// SeverityUpdate is the builder for updating Severity entities.
type SeverityUpdate struct {
	config
	hooks    []Hook
	mutation *SeverityMutation
}

// Where appends a list predicates to the SeverityUpdate builder.
func (su *SeverityUpdate) Where(ps ...predicate.Severity) *SeverityUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetLabel sets the "label" field.
func (su *SeverityUpdate) SetLabel(s string) *SeverityUpdate {
	su.mutation.SetLabel(s)
	return su
}

// SetColor sets the "color" field.
func (su *SeverityUpdate) SetColor(s string) *SeverityUpdate {
	su.mutation.SetColor(s)
	return su
}

// SetNillableColor sets the "color" field if the given value is not nil.
func (su *SeverityUpdate) SetNillableColor(s *string) *SeverityUpdate {
	if s != nil {
		su.SetColor(*s)
	}
	return su
}

// ClearColor clears the value of the "color" field.
func (su *SeverityUpdate) ClearColor() *SeverityUpdate {
	su.mutation.ClearColor()
	return su
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the Vulnerability entity by IDs.
func (su *SeverityUpdate) AddVulnerabilityIDs(ids ...string) *SeverityUpdate {
	su.mutation.AddVulnerabilityIDs(ids...)
	return su
}

// AddVulnerabilities adds the "vulnerabilities" edges to the Vulnerability entity.
func (su *SeverityUpdate) AddVulnerabilities(v ...*Vulnerability) *SeverityUpdate {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return su.AddVulnerabilityIDs(ids...)
}

// Mutation returns the SeverityMutation object of the builder.
func (su *SeverityUpdate) Mutation() *SeverityMutation {
	return su.mutation
}

// ClearVulnerabilities clears all "vulnerabilities" edges to the Vulnerability entity.
func (su *SeverityUpdate) ClearVulnerabilities() *SeverityUpdate {
	su.mutation.ClearVulnerabilities()
	return su
}

// RemoveVulnerabilityIDs removes the "vulnerabilities" edge to Vulnerability entities by IDs.
func (su *SeverityUpdate) RemoveVulnerabilityIDs(ids ...string) *SeverityUpdate {
	su.mutation.RemoveVulnerabilityIDs(ids...)
	return su
}

// RemoveVulnerabilities removes "vulnerabilities" edges to Vulnerability entities.
func (su *SeverityUpdate) RemoveVulnerabilities(v ...*Vulnerability) *SeverityUpdate {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return su.RemoveVulnerabilityIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *SeverityUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(su.hooks) == 0 {
		if err = su.check(); err != nil {
			return 0, err
		}
		affected, err = su.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SeverityMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = su.check(); err != nil {
				return 0, err
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
func (su *SeverityUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *SeverityUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *SeverityUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (su *SeverityUpdate) check() error {
	if v, ok := su.mutation.Label(); ok {
		if err := severity.LabelValidator(v); err != nil {
			return &ValidationError{Name: "label", err: fmt.Errorf("ent: validator failed for field \"label\": %w", err)}
		}
	}
	return nil
}

func (su *SeverityUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   severity.Table,
			Columns: severity.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: severity.FieldID,
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
	if value, ok := su.mutation.Label(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: severity.FieldLabel,
		})
	}
	if value, ok := su.mutation.Color(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: severity.FieldColor,
		})
	}
	if su.mutation.ColorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: severity.FieldColor,
		})
	}
	if su.mutation.VulnerabilitiesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedVulnerabilitiesIDs(); len(nodes) > 0 && !su.mutation.VulnerabilitiesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.VulnerabilitiesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{severity.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// SeverityUpdateOne is the builder for updating a single Severity entity.
type SeverityUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SeverityMutation
}

// SetLabel sets the "label" field.
func (suo *SeverityUpdateOne) SetLabel(s string) *SeverityUpdateOne {
	suo.mutation.SetLabel(s)
	return suo
}

// SetColor sets the "color" field.
func (suo *SeverityUpdateOne) SetColor(s string) *SeverityUpdateOne {
	suo.mutation.SetColor(s)
	return suo
}

// SetNillableColor sets the "color" field if the given value is not nil.
func (suo *SeverityUpdateOne) SetNillableColor(s *string) *SeverityUpdateOne {
	if s != nil {
		suo.SetColor(*s)
	}
	return suo
}

// ClearColor clears the value of the "color" field.
func (suo *SeverityUpdateOne) ClearColor() *SeverityUpdateOne {
	suo.mutation.ClearColor()
	return suo
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the Vulnerability entity by IDs.
func (suo *SeverityUpdateOne) AddVulnerabilityIDs(ids ...string) *SeverityUpdateOne {
	suo.mutation.AddVulnerabilityIDs(ids...)
	return suo
}

// AddVulnerabilities adds the "vulnerabilities" edges to the Vulnerability entity.
func (suo *SeverityUpdateOne) AddVulnerabilities(v ...*Vulnerability) *SeverityUpdateOne {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return suo.AddVulnerabilityIDs(ids...)
}

// Mutation returns the SeverityMutation object of the builder.
func (suo *SeverityUpdateOne) Mutation() *SeverityMutation {
	return suo.mutation
}

// ClearVulnerabilities clears all "vulnerabilities" edges to the Vulnerability entity.
func (suo *SeverityUpdateOne) ClearVulnerabilities() *SeverityUpdateOne {
	suo.mutation.ClearVulnerabilities()
	return suo
}

// RemoveVulnerabilityIDs removes the "vulnerabilities" edge to Vulnerability entities by IDs.
func (suo *SeverityUpdateOne) RemoveVulnerabilityIDs(ids ...string) *SeverityUpdateOne {
	suo.mutation.RemoveVulnerabilityIDs(ids...)
	return suo
}

// RemoveVulnerabilities removes "vulnerabilities" edges to Vulnerability entities.
func (suo *SeverityUpdateOne) RemoveVulnerabilities(v ...*Vulnerability) *SeverityUpdateOne {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return suo.RemoveVulnerabilityIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *SeverityUpdateOne) Select(field string, fields ...string) *SeverityUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Severity entity.
func (suo *SeverityUpdateOne) Save(ctx context.Context) (*Severity, error) {
	var (
		err  error
		node *Severity
	)
	if len(suo.hooks) == 0 {
		if err = suo.check(); err != nil {
			return nil, err
		}
		node, err = suo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SeverityMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = suo.check(); err != nil {
				return nil, err
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
func (suo *SeverityUpdateOne) SaveX(ctx context.Context) *Severity {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *SeverityUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *SeverityUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (suo *SeverityUpdateOne) check() error {
	if v, ok := suo.mutation.Label(); ok {
		if err := severity.LabelValidator(v); err != nil {
			return &ValidationError{Name: "label", err: fmt.Errorf("ent: validator failed for field \"label\": %w", err)}
		}
	}
	return nil
}

func (suo *SeverityUpdateOne) sqlSave(ctx context.Context) (_node *Severity, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   severity.Table,
			Columns: severity.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: severity.FieldID,
			},
		},
	}
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Severity.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, severity.FieldID)
		for _, f := range fields {
			if !severity.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != severity.FieldID {
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
	if value, ok := suo.mutation.Label(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: severity.FieldLabel,
		})
	}
	if value, ok := suo.mutation.Color(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: severity.FieldColor,
		})
	}
	if suo.mutation.ColorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: severity.FieldColor,
		})
	}
	if suo.mutation.VulnerabilitiesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedVulnerabilitiesIDs(); len(nodes) > 0 && !suo.mutation.VulnerabilitiesCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.VulnerabilitiesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Severity{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{severity.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}