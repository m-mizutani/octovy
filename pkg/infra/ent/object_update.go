// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/octovy/pkg/infra/ent/object"
	"github.com/m-mizutani/octovy/pkg/infra/ent/predicate"
	"github.com/m-mizutani/octovy/pkg/infra/ent/report"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnerability"
)

// ObjectUpdate is the builder for updating Object entities.
type ObjectUpdate struct {
	config
	hooks    []Hook
	mutation *ObjectMutation
}

// Where appends a list predicates to the ObjectUpdate builder.
func (ou *ObjectUpdate) Where(ps ...predicate.Object) *ObjectUpdate {
	ou.mutation.Where(ps...)
	return ou
}

// SetDescription sets the "description" field.
func (ou *ObjectUpdate) SetDescription(s string) *ObjectUpdate {
	ou.mutation.SetDescription(s)
	return ou
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (ou *ObjectUpdate) SetNillableDescription(s *string) *ObjectUpdate {
	if s != nil {
		ou.SetDescription(*s)
	}
	return ou
}

// ClearDescription clears the value of the "description" field.
func (ou *ObjectUpdate) ClearDescription() *ObjectUpdate {
	ou.mutation.ClearDescription()
	return ou
}

// SetVersion sets the "version" field.
func (ou *ObjectUpdate) SetVersion(s string) *ObjectUpdate {
	ou.mutation.SetVersion(s)
	return ou
}

// SetNillableVersion sets the "version" field if the given value is not nil.
func (ou *ObjectUpdate) SetNillableVersion(s *string) *ObjectUpdate {
	if s != nil {
		ou.SetVersion(*s)
	}
	return ou
}

// ClearVersion clears the value of the "version" field.
func (ou *ObjectUpdate) ClearVersion() *ObjectUpdate {
	ou.mutation.ClearVersion()
	return ou
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the Vulnerability entity by IDs.
func (ou *ObjectUpdate) AddVulnerabilityIDs(ids ...string) *ObjectUpdate {
	ou.mutation.AddVulnerabilityIDs(ids...)
	return ou
}

// AddVulnerabilities adds the "vulnerabilities" edges to the Vulnerability entity.
func (ou *ObjectUpdate) AddVulnerabilities(v ...*Vulnerability) *ObjectUpdate {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return ou.AddVulnerabilityIDs(ids...)
}

// AddReportIDs adds the "report" edge to the Report entity by IDs.
func (ou *ObjectUpdate) AddReportIDs(ids ...int) *ObjectUpdate {
	ou.mutation.AddReportIDs(ids...)
	return ou
}

// AddReport adds the "report" edges to the Report entity.
func (ou *ObjectUpdate) AddReport(r ...*Report) *ObjectUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ou.AddReportIDs(ids...)
}

// Mutation returns the ObjectMutation object of the builder.
func (ou *ObjectUpdate) Mutation() *ObjectMutation {
	return ou.mutation
}

// ClearVulnerabilities clears all "vulnerabilities" edges to the Vulnerability entity.
func (ou *ObjectUpdate) ClearVulnerabilities() *ObjectUpdate {
	ou.mutation.ClearVulnerabilities()
	return ou
}

// RemoveVulnerabilityIDs removes the "vulnerabilities" edge to Vulnerability entities by IDs.
func (ou *ObjectUpdate) RemoveVulnerabilityIDs(ids ...string) *ObjectUpdate {
	ou.mutation.RemoveVulnerabilityIDs(ids...)
	return ou
}

// RemoveVulnerabilities removes "vulnerabilities" edges to Vulnerability entities.
func (ou *ObjectUpdate) RemoveVulnerabilities(v ...*Vulnerability) *ObjectUpdate {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return ou.RemoveVulnerabilityIDs(ids...)
}

// ClearReport clears all "report" edges to the Report entity.
func (ou *ObjectUpdate) ClearReport() *ObjectUpdate {
	ou.mutation.ClearReport()
	return ou
}

// RemoveReportIDs removes the "report" edge to Report entities by IDs.
func (ou *ObjectUpdate) RemoveReportIDs(ids ...int) *ObjectUpdate {
	ou.mutation.RemoveReportIDs(ids...)
	return ou
}

// RemoveReport removes "report" edges to Report entities.
func (ou *ObjectUpdate) RemoveReport(r ...*Report) *ObjectUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ou.RemoveReportIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ou *ObjectUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(ou.hooks) == 0 {
		affected, err = ou.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ObjectMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ou.mutation = mutation
			affected, err = ou.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(ou.hooks) - 1; i >= 0; i-- {
			if ou.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ou.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ou.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (ou *ObjectUpdate) SaveX(ctx context.Context) int {
	affected, err := ou.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ou *ObjectUpdate) Exec(ctx context.Context) error {
	_, err := ou.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ou *ObjectUpdate) ExecX(ctx context.Context) {
	if err := ou.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ou *ObjectUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   object.Table,
			Columns: object.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: object.FieldID,
			},
		},
	}
	if ps := ou.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ou.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: object.FieldDescription,
		})
	}
	if ou.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: object.FieldDescription,
		})
	}
	if value, ok := ou.mutation.Version(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: object.FieldVersion,
		})
	}
	if ou.mutation.VersionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: object.FieldVersion,
		})
	}
	if ou.mutation.VulnerabilitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   object.VulnerabilitiesTable,
			Columns: []string{object.VulnerabilitiesColumn},
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
	if nodes := ou.mutation.RemovedVulnerabilitiesIDs(); len(nodes) > 0 && !ou.mutation.VulnerabilitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   object.VulnerabilitiesTable,
			Columns: []string{object.VulnerabilitiesColumn},
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
	if nodes := ou.mutation.VulnerabilitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   object.VulnerabilitiesTable,
			Columns: []string{object.VulnerabilitiesColumn},
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
	if ou.mutation.ReportCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   object.ReportTable,
			Columns: object.ReportPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: report.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ou.mutation.RemovedReportIDs(); len(nodes) > 0 && !ou.mutation.ReportCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   object.ReportTable,
			Columns: object.ReportPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: report.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ou.mutation.ReportIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   object.ReportTable,
			Columns: object.ReportPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: report.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ou.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{object.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// ObjectUpdateOne is the builder for updating a single Object entity.
type ObjectUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ObjectMutation
}

// SetDescription sets the "description" field.
func (ouo *ObjectUpdateOne) SetDescription(s string) *ObjectUpdateOne {
	ouo.mutation.SetDescription(s)
	return ouo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (ouo *ObjectUpdateOne) SetNillableDescription(s *string) *ObjectUpdateOne {
	if s != nil {
		ouo.SetDescription(*s)
	}
	return ouo
}

// ClearDescription clears the value of the "description" field.
func (ouo *ObjectUpdateOne) ClearDescription() *ObjectUpdateOne {
	ouo.mutation.ClearDescription()
	return ouo
}

// SetVersion sets the "version" field.
func (ouo *ObjectUpdateOne) SetVersion(s string) *ObjectUpdateOne {
	ouo.mutation.SetVersion(s)
	return ouo
}

// SetNillableVersion sets the "version" field if the given value is not nil.
func (ouo *ObjectUpdateOne) SetNillableVersion(s *string) *ObjectUpdateOne {
	if s != nil {
		ouo.SetVersion(*s)
	}
	return ouo
}

// ClearVersion clears the value of the "version" field.
func (ouo *ObjectUpdateOne) ClearVersion() *ObjectUpdateOne {
	ouo.mutation.ClearVersion()
	return ouo
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the Vulnerability entity by IDs.
func (ouo *ObjectUpdateOne) AddVulnerabilityIDs(ids ...string) *ObjectUpdateOne {
	ouo.mutation.AddVulnerabilityIDs(ids...)
	return ouo
}

// AddVulnerabilities adds the "vulnerabilities" edges to the Vulnerability entity.
func (ouo *ObjectUpdateOne) AddVulnerabilities(v ...*Vulnerability) *ObjectUpdateOne {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return ouo.AddVulnerabilityIDs(ids...)
}

// AddReportIDs adds the "report" edge to the Report entity by IDs.
func (ouo *ObjectUpdateOne) AddReportIDs(ids ...int) *ObjectUpdateOne {
	ouo.mutation.AddReportIDs(ids...)
	return ouo
}

// AddReport adds the "report" edges to the Report entity.
func (ouo *ObjectUpdateOne) AddReport(r ...*Report) *ObjectUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ouo.AddReportIDs(ids...)
}

// Mutation returns the ObjectMutation object of the builder.
func (ouo *ObjectUpdateOne) Mutation() *ObjectMutation {
	return ouo.mutation
}

// ClearVulnerabilities clears all "vulnerabilities" edges to the Vulnerability entity.
func (ouo *ObjectUpdateOne) ClearVulnerabilities() *ObjectUpdateOne {
	ouo.mutation.ClearVulnerabilities()
	return ouo
}

// RemoveVulnerabilityIDs removes the "vulnerabilities" edge to Vulnerability entities by IDs.
func (ouo *ObjectUpdateOne) RemoveVulnerabilityIDs(ids ...string) *ObjectUpdateOne {
	ouo.mutation.RemoveVulnerabilityIDs(ids...)
	return ouo
}

// RemoveVulnerabilities removes "vulnerabilities" edges to Vulnerability entities.
func (ouo *ObjectUpdateOne) RemoveVulnerabilities(v ...*Vulnerability) *ObjectUpdateOne {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return ouo.RemoveVulnerabilityIDs(ids...)
}

// ClearReport clears all "report" edges to the Report entity.
func (ouo *ObjectUpdateOne) ClearReport() *ObjectUpdateOne {
	ouo.mutation.ClearReport()
	return ouo
}

// RemoveReportIDs removes the "report" edge to Report entities by IDs.
func (ouo *ObjectUpdateOne) RemoveReportIDs(ids ...int) *ObjectUpdateOne {
	ouo.mutation.RemoveReportIDs(ids...)
	return ouo
}

// RemoveReport removes "report" edges to Report entities.
func (ouo *ObjectUpdateOne) RemoveReport(r ...*Report) *ObjectUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ouo.RemoveReportIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ouo *ObjectUpdateOne) Select(field string, fields ...string) *ObjectUpdateOne {
	ouo.fields = append([]string{field}, fields...)
	return ouo
}

// Save executes the query and returns the updated Object entity.
func (ouo *ObjectUpdateOne) Save(ctx context.Context) (*Object, error) {
	var (
		err  error
		node *Object
	)
	if len(ouo.hooks) == 0 {
		node, err = ouo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ObjectMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ouo.mutation = mutation
			node, err = ouo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(ouo.hooks) - 1; i >= 0; i-- {
			if ouo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ouo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ouo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (ouo *ObjectUpdateOne) SaveX(ctx context.Context) *Object {
	node, err := ouo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ouo *ObjectUpdateOne) Exec(ctx context.Context) error {
	_, err := ouo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ouo *ObjectUpdateOne) ExecX(ctx context.Context) {
	if err := ouo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ouo *ObjectUpdateOne) sqlSave(ctx context.Context) (_node *Object, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   object.Table,
			Columns: object.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: object.FieldID,
			},
		},
	}
	id, ok := ouo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Object.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := ouo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, object.FieldID)
		for _, f := range fields {
			if !object.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != object.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ouo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ouo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: object.FieldDescription,
		})
	}
	if ouo.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: object.FieldDescription,
		})
	}
	if value, ok := ouo.mutation.Version(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: object.FieldVersion,
		})
	}
	if ouo.mutation.VersionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: object.FieldVersion,
		})
	}
	if ouo.mutation.VulnerabilitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   object.VulnerabilitiesTable,
			Columns: []string{object.VulnerabilitiesColumn},
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
	if nodes := ouo.mutation.RemovedVulnerabilitiesIDs(); len(nodes) > 0 && !ouo.mutation.VulnerabilitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   object.VulnerabilitiesTable,
			Columns: []string{object.VulnerabilitiesColumn},
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
	if nodes := ouo.mutation.VulnerabilitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   object.VulnerabilitiesTable,
			Columns: []string{object.VulnerabilitiesColumn},
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
	if ouo.mutation.ReportCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   object.ReportTable,
			Columns: object.ReportPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: report.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ouo.mutation.RemovedReportIDs(); len(nodes) > 0 && !ouo.mutation.ReportCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   object.ReportTable,
			Columns: object.ReportPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: report.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ouo.mutation.ReportIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   object.ReportTable,
			Columns: object.ReportPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: report.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Object{config: ouo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ouo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{object.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
