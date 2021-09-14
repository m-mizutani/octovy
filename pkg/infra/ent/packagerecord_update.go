// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/ent/packagerecord"
	"github.com/m-mizutani/octovy/pkg/infra/ent/predicate"
	"github.com/m-mizutani/octovy/pkg/infra/ent/scan"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnerability"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnstatus"
)

// PackageRecordUpdate is the builder for updating PackageRecord entities.
type PackageRecordUpdate struct {
	config
	hooks    []Hook
	mutation *PackageRecordMutation
}

// Where appends a list predicates to the PackageRecordUpdate builder.
func (pru *PackageRecordUpdate) Where(ps ...predicate.PackageRecord) *PackageRecordUpdate {
	pru.mutation.Where(ps...)
	return pru
}

// SetType sets the "type" field.
func (pru *PackageRecordUpdate) SetType(tt types.PkgType) *PackageRecordUpdate {
	pru.mutation.SetType(tt)
	return pru
}

// AddScanIDs adds the "scan" edge to the Scan entity by IDs.
func (pru *PackageRecordUpdate) AddScanIDs(ids ...string) *PackageRecordUpdate {
	pru.mutation.AddScanIDs(ids...)
	return pru
}

// AddScan adds the "scan" edges to the Scan entity.
func (pru *PackageRecordUpdate) AddScan(s ...*Scan) *PackageRecordUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return pru.AddScanIDs(ids...)
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the Vulnerability entity by IDs.
func (pru *PackageRecordUpdate) AddVulnerabilityIDs(ids ...string) *PackageRecordUpdate {
	pru.mutation.AddVulnerabilityIDs(ids...)
	return pru
}

// AddVulnerabilities adds the "vulnerabilities" edges to the Vulnerability entity.
func (pru *PackageRecordUpdate) AddVulnerabilities(v ...*Vulnerability) *PackageRecordUpdate {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return pru.AddVulnerabilityIDs(ids...)
}

// AddStatuIDs adds the "status" edge to the VulnStatus entity by IDs.
func (pru *PackageRecordUpdate) AddStatuIDs(ids ...string) *PackageRecordUpdate {
	pru.mutation.AddStatuIDs(ids...)
	return pru
}

// AddStatus adds the "status" edges to the VulnStatus entity.
func (pru *PackageRecordUpdate) AddStatus(v ...*VulnStatus) *PackageRecordUpdate {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return pru.AddStatuIDs(ids...)
}

// Mutation returns the PackageRecordMutation object of the builder.
func (pru *PackageRecordUpdate) Mutation() *PackageRecordMutation {
	return pru.mutation
}

// ClearScan clears all "scan" edges to the Scan entity.
func (pru *PackageRecordUpdate) ClearScan() *PackageRecordUpdate {
	pru.mutation.ClearScan()
	return pru
}

// RemoveScanIDs removes the "scan" edge to Scan entities by IDs.
func (pru *PackageRecordUpdate) RemoveScanIDs(ids ...string) *PackageRecordUpdate {
	pru.mutation.RemoveScanIDs(ids...)
	return pru
}

// RemoveScan removes "scan" edges to Scan entities.
func (pru *PackageRecordUpdate) RemoveScan(s ...*Scan) *PackageRecordUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return pru.RemoveScanIDs(ids...)
}

// ClearVulnerabilities clears all "vulnerabilities" edges to the Vulnerability entity.
func (pru *PackageRecordUpdate) ClearVulnerabilities() *PackageRecordUpdate {
	pru.mutation.ClearVulnerabilities()
	return pru
}

// RemoveVulnerabilityIDs removes the "vulnerabilities" edge to Vulnerability entities by IDs.
func (pru *PackageRecordUpdate) RemoveVulnerabilityIDs(ids ...string) *PackageRecordUpdate {
	pru.mutation.RemoveVulnerabilityIDs(ids...)
	return pru
}

// RemoveVulnerabilities removes "vulnerabilities" edges to Vulnerability entities.
func (pru *PackageRecordUpdate) RemoveVulnerabilities(v ...*Vulnerability) *PackageRecordUpdate {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return pru.RemoveVulnerabilityIDs(ids...)
}

// ClearStatus clears all "status" edges to the VulnStatus entity.
func (pru *PackageRecordUpdate) ClearStatus() *PackageRecordUpdate {
	pru.mutation.ClearStatus()
	return pru
}

// RemoveStatuIDs removes the "status" edge to VulnStatus entities by IDs.
func (pru *PackageRecordUpdate) RemoveStatuIDs(ids ...string) *PackageRecordUpdate {
	pru.mutation.RemoveStatuIDs(ids...)
	return pru
}

// RemoveStatus removes "status" edges to VulnStatus entities.
func (pru *PackageRecordUpdate) RemoveStatus(v ...*VulnStatus) *PackageRecordUpdate {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return pru.RemoveStatuIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pru *PackageRecordUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(pru.hooks) == 0 {
		if err = pru.check(); err != nil {
			return 0, err
		}
		affected, err = pru.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PackageRecordMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = pru.check(); err != nil {
				return 0, err
			}
			pru.mutation = mutation
			affected, err = pru.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(pru.hooks) - 1; i >= 0; i-- {
			if pru.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = pru.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, pru.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (pru *PackageRecordUpdate) SaveX(ctx context.Context) int {
	affected, err := pru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pru *PackageRecordUpdate) Exec(ctx context.Context) error {
	_, err := pru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pru *PackageRecordUpdate) ExecX(ctx context.Context) {
	if err := pru.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pru *PackageRecordUpdate) check() error {
	if v, ok := pru.mutation.GetType(); ok {
		if err := packagerecord.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf("ent: validator failed for field \"type\": %w", err)}
		}
	}
	return nil
}

func (pru *PackageRecordUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   packagerecord.Table,
			Columns: packagerecord.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: packagerecord.FieldID,
			},
		},
	}
	if ps := pru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pru.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: packagerecord.FieldType,
		})
	}
	if pru.mutation.ScanCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   packagerecord.ScanTable,
			Columns: packagerecord.ScanPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: scan.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pru.mutation.RemovedScanIDs(); len(nodes) > 0 && !pru.mutation.ScanCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   packagerecord.ScanTable,
			Columns: packagerecord.ScanPrimaryKey,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pru.mutation.ScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   packagerecord.ScanTable,
			Columns: packagerecord.ScanPrimaryKey,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pru.mutation.VulnerabilitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   packagerecord.VulnerabilitiesTable,
			Columns: packagerecord.VulnerabilitiesPrimaryKey,
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
	if nodes := pru.mutation.RemovedVulnerabilitiesIDs(); len(nodes) > 0 && !pru.mutation.VulnerabilitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   packagerecord.VulnerabilitiesTable,
			Columns: packagerecord.VulnerabilitiesPrimaryKey,
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
	if nodes := pru.mutation.VulnerabilitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   packagerecord.VulnerabilitiesTable,
			Columns: packagerecord.VulnerabilitiesPrimaryKey,
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
	if pru.mutation.StatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   packagerecord.StatusTable,
			Columns: []string{packagerecord.StatusColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: vulnstatus.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pru.mutation.RemovedStatusIDs(); len(nodes) > 0 && !pru.mutation.StatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   packagerecord.StatusTable,
			Columns: []string{packagerecord.StatusColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: vulnstatus.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pru.mutation.StatusIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   packagerecord.StatusTable,
			Columns: []string{packagerecord.StatusColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: vulnstatus.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{packagerecord.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// PackageRecordUpdateOne is the builder for updating a single PackageRecord entity.
type PackageRecordUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *PackageRecordMutation
}

// SetType sets the "type" field.
func (pruo *PackageRecordUpdateOne) SetType(tt types.PkgType) *PackageRecordUpdateOne {
	pruo.mutation.SetType(tt)
	return pruo
}

// AddScanIDs adds the "scan" edge to the Scan entity by IDs.
func (pruo *PackageRecordUpdateOne) AddScanIDs(ids ...string) *PackageRecordUpdateOne {
	pruo.mutation.AddScanIDs(ids...)
	return pruo
}

// AddScan adds the "scan" edges to the Scan entity.
func (pruo *PackageRecordUpdateOne) AddScan(s ...*Scan) *PackageRecordUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return pruo.AddScanIDs(ids...)
}

// AddVulnerabilityIDs adds the "vulnerabilities" edge to the Vulnerability entity by IDs.
func (pruo *PackageRecordUpdateOne) AddVulnerabilityIDs(ids ...string) *PackageRecordUpdateOne {
	pruo.mutation.AddVulnerabilityIDs(ids...)
	return pruo
}

// AddVulnerabilities adds the "vulnerabilities" edges to the Vulnerability entity.
func (pruo *PackageRecordUpdateOne) AddVulnerabilities(v ...*Vulnerability) *PackageRecordUpdateOne {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return pruo.AddVulnerabilityIDs(ids...)
}

// AddStatuIDs adds the "status" edge to the VulnStatus entity by IDs.
func (pruo *PackageRecordUpdateOne) AddStatuIDs(ids ...string) *PackageRecordUpdateOne {
	pruo.mutation.AddStatuIDs(ids...)
	return pruo
}

// AddStatus adds the "status" edges to the VulnStatus entity.
func (pruo *PackageRecordUpdateOne) AddStatus(v ...*VulnStatus) *PackageRecordUpdateOne {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return pruo.AddStatuIDs(ids...)
}

// Mutation returns the PackageRecordMutation object of the builder.
func (pruo *PackageRecordUpdateOne) Mutation() *PackageRecordMutation {
	return pruo.mutation
}

// ClearScan clears all "scan" edges to the Scan entity.
func (pruo *PackageRecordUpdateOne) ClearScan() *PackageRecordUpdateOne {
	pruo.mutation.ClearScan()
	return pruo
}

// RemoveScanIDs removes the "scan" edge to Scan entities by IDs.
func (pruo *PackageRecordUpdateOne) RemoveScanIDs(ids ...string) *PackageRecordUpdateOne {
	pruo.mutation.RemoveScanIDs(ids...)
	return pruo
}

// RemoveScan removes "scan" edges to Scan entities.
func (pruo *PackageRecordUpdateOne) RemoveScan(s ...*Scan) *PackageRecordUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return pruo.RemoveScanIDs(ids...)
}

// ClearVulnerabilities clears all "vulnerabilities" edges to the Vulnerability entity.
func (pruo *PackageRecordUpdateOne) ClearVulnerabilities() *PackageRecordUpdateOne {
	pruo.mutation.ClearVulnerabilities()
	return pruo
}

// RemoveVulnerabilityIDs removes the "vulnerabilities" edge to Vulnerability entities by IDs.
func (pruo *PackageRecordUpdateOne) RemoveVulnerabilityIDs(ids ...string) *PackageRecordUpdateOne {
	pruo.mutation.RemoveVulnerabilityIDs(ids...)
	return pruo
}

// RemoveVulnerabilities removes "vulnerabilities" edges to Vulnerability entities.
func (pruo *PackageRecordUpdateOne) RemoveVulnerabilities(v ...*Vulnerability) *PackageRecordUpdateOne {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return pruo.RemoveVulnerabilityIDs(ids...)
}

// ClearStatus clears all "status" edges to the VulnStatus entity.
func (pruo *PackageRecordUpdateOne) ClearStatus() *PackageRecordUpdateOne {
	pruo.mutation.ClearStatus()
	return pruo
}

// RemoveStatuIDs removes the "status" edge to VulnStatus entities by IDs.
func (pruo *PackageRecordUpdateOne) RemoveStatuIDs(ids ...string) *PackageRecordUpdateOne {
	pruo.mutation.RemoveStatuIDs(ids...)
	return pruo
}

// RemoveStatus removes "status" edges to VulnStatus entities.
func (pruo *PackageRecordUpdateOne) RemoveStatus(v ...*VulnStatus) *PackageRecordUpdateOne {
	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return pruo.RemoveStatuIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (pruo *PackageRecordUpdateOne) Select(field string, fields ...string) *PackageRecordUpdateOne {
	pruo.fields = append([]string{field}, fields...)
	return pruo
}

// Save executes the query and returns the updated PackageRecord entity.
func (pruo *PackageRecordUpdateOne) Save(ctx context.Context) (*PackageRecord, error) {
	var (
		err  error
		node *PackageRecord
	)
	if len(pruo.hooks) == 0 {
		if err = pruo.check(); err != nil {
			return nil, err
		}
		node, err = pruo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PackageRecordMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = pruo.check(); err != nil {
				return nil, err
			}
			pruo.mutation = mutation
			node, err = pruo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(pruo.hooks) - 1; i >= 0; i-- {
			if pruo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = pruo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, pruo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (pruo *PackageRecordUpdateOne) SaveX(ctx context.Context) *PackageRecord {
	node, err := pruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (pruo *PackageRecordUpdateOne) Exec(ctx context.Context) error {
	_, err := pruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pruo *PackageRecordUpdateOne) ExecX(ctx context.Context) {
	if err := pruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pruo *PackageRecordUpdateOne) check() error {
	if v, ok := pruo.mutation.GetType(); ok {
		if err := packagerecord.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf("ent: validator failed for field \"type\": %w", err)}
		}
	}
	return nil
}

func (pruo *PackageRecordUpdateOne) sqlSave(ctx context.Context) (_node *PackageRecord, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   packagerecord.Table,
			Columns: packagerecord.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: packagerecord.FieldID,
			},
		},
	}
	id, ok := pruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing PackageRecord.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := pruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, packagerecord.FieldID)
		for _, f := range fields {
			if !packagerecord.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != packagerecord.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := pruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pruo.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: packagerecord.FieldType,
		})
	}
	if pruo.mutation.ScanCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   packagerecord.ScanTable,
			Columns: packagerecord.ScanPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: scan.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pruo.mutation.RemovedScanIDs(); len(nodes) > 0 && !pruo.mutation.ScanCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   packagerecord.ScanTable,
			Columns: packagerecord.ScanPrimaryKey,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pruo.mutation.ScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   packagerecord.ScanTable,
			Columns: packagerecord.ScanPrimaryKey,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pruo.mutation.VulnerabilitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   packagerecord.VulnerabilitiesTable,
			Columns: packagerecord.VulnerabilitiesPrimaryKey,
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
	if nodes := pruo.mutation.RemovedVulnerabilitiesIDs(); len(nodes) > 0 && !pruo.mutation.VulnerabilitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   packagerecord.VulnerabilitiesTable,
			Columns: packagerecord.VulnerabilitiesPrimaryKey,
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
	if nodes := pruo.mutation.VulnerabilitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   packagerecord.VulnerabilitiesTable,
			Columns: packagerecord.VulnerabilitiesPrimaryKey,
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
	if pruo.mutation.StatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   packagerecord.StatusTable,
			Columns: []string{packagerecord.StatusColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: vulnstatus.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pruo.mutation.RemovedStatusIDs(); len(nodes) > 0 && !pruo.mutation.StatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   packagerecord.StatusTable,
			Columns: []string{packagerecord.StatusColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: vulnstatus.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pruo.mutation.StatusIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   packagerecord.StatusTable,
			Columns: []string{packagerecord.StatusColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: vulnstatus.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &PackageRecord{config: pruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, pruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{packagerecord.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
