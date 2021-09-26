// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/octovy/pkg/infra/ent/predicate"
	"github.com/m-mizutani/octovy/pkg/infra/ent/user"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnstatus"
)

// VulnStatusQuery is the builder for querying VulnStatus entities.
type VulnStatusQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.VulnStatus
	// eager-loading edges.
	withAuthor *UserQuery
	withFKs    bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the VulnStatusQuery builder.
func (vsq *VulnStatusQuery) Where(ps ...predicate.VulnStatus) *VulnStatusQuery {
	vsq.predicates = append(vsq.predicates, ps...)
	return vsq
}

// Limit adds a limit step to the query.
func (vsq *VulnStatusQuery) Limit(limit int) *VulnStatusQuery {
	vsq.limit = &limit
	return vsq
}

// Offset adds an offset step to the query.
func (vsq *VulnStatusQuery) Offset(offset int) *VulnStatusQuery {
	vsq.offset = &offset
	return vsq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (vsq *VulnStatusQuery) Unique(unique bool) *VulnStatusQuery {
	vsq.unique = &unique
	return vsq
}

// Order adds an order step to the query.
func (vsq *VulnStatusQuery) Order(o ...OrderFunc) *VulnStatusQuery {
	vsq.order = append(vsq.order, o...)
	return vsq
}

// QueryAuthor chains the current query on the "author" edge.
func (vsq *VulnStatusQuery) QueryAuthor() *UserQuery {
	query := &UserQuery{config: vsq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := vsq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := vsq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(vulnstatus.Table, vulnstatus.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, vulnstatus.AuthorTable, vulnstatus.AuthorColumn),
		)
		fromU = sqlgraph.SetNeighbors(vsq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first VulnStatus entity from the query.
// Returns a *NotFoundError when no VulnStatus was found.
func (vsq *VulnStatusQuery) First(ctx context.Context) (*VulnStatus, error) {
	nodes, err := vsq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{vulnstatus.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (vsq *VulnStatusQuery) FirstX(ctx context.Context) *VulnStatus {
	node, err := vsq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first VulnStatus ID from the query.
// Returns a *NotFoundError when no VulnStatus ID was found.
func (vsq *VulnStatusQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = vsq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{vulnstatus.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (vsq *VulnStatusQuery) FirstIDX(ctx context.Context) int {
	id, err := vsq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single VulnStatus entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when exactly one VulnStatus entity is not found.
// Returns a *NotFoundError when no VulnStatus entities are found.
func (vsq *VulnStatusQuery) Only(ctx context.Context) (*VulnStatus, error) {
	nodes, err := vsq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{vulnstatus.Label}
	default:
		return nil, &NotSingularError{vulnstatus.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (vsq *VulnStatusQuery) OnlyX(ctx context.Context) *VulnStatus {
	node, err := vsq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only VulnStatus ID in the query.
// Returns a *NotSingularError when exactly one VulnStatus ID is not found.
// Returns a *NotFoundError when no entities are found.
func (vsq *VulnStatusQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = vsq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{vulnstatus.Label}
	default:
		err = &NotSingularError{vulnstatus.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (vsq *VulnStatusQuery) OnlyIDX(ctx context.Context) int {
	id, err := vsq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of VulnStatusSlice.
func (vsq *VulnStatusQuery) All(ctx context.Context) ([]*VulnStatus, error) {
	if err := vsq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return vsq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (vsq *VulnStatusQuery) AllX(ctx context.Context) []*VulnStatus {
	nodes, err := vsq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of VulnStatus IDs.
func (vsq *VulnStatusQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := vsq.Select(vulnstatus.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (vsq *VulnStatusQuery) IDsX(ctx context.Context) []int {
	ids, err := vsq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (vsq *VulnStatusQuery) Count(ctx context.Context) (int, error) {
	if err := vsq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return vsq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (vsq *VulnStatusQuery) CountX(ctx context.Context) int {
	count, err := vsq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (vsq *VulnStatusQuery) Exist(ctx context.Context) (bool, error) {
	if err := vsq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return vsq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (vsq *VulnStatusQuery) ExistX(ctx context.Context) bool {
	exist, err := vsq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the VulnStatusQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (vsq *VulnStatusQuery) Clone() *VulnStatusQuery {
	if vsq == nil {
		return nil
	}
	return &VulnStatusQuery{
		config:     vsq.config,
		limit:      vsq.limit,
		offset:     vsq.offset,
		order:      append([]OrderFunc{}, vsq.order...),
		predicates: append([]predicate.VulnStatus{}, vsq.predicates...),
		withAuthor: vsq.withAuthor.Clone(),
		// clone intermediate query.
		sql:  vsq.sql.Clone(),
		path: vsq.path,
	}
}

// WithAuthor tells the query-builder to eager-load the nodes that are connected to
// the "author" edge. The optional arguments are used to configure the query builder of the edge.
func (vsq *VulnStatusQuery) WithAuthor(opts ...func(*UserQuery)) *VulnStatusQuery {
	query := &UserQuery{config: vsq.config}
	for _, opt := range opts {
		opt(query)
	}
	vsq.withAuthor = query
	return vsq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Status types.VulnStatusType `json:"status,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.VulnStatus.Query().
//		GroupBy(vulnstatus.FieldStatus).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (vsq *VulnStatusQuery) GroupBy(field string, fields ...string) *VulnStatusGroupBy {
	group := &VulnStatusGroupBy{config: vsq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := vsq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return vsq.sqlQuery(ctx), nil
	}
	return group
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Status types.VulnStatusType `json:"status,omitempty"`
//	}
//
//	client.VulnStatus.Query().
//		Select(vulnstatus.FieldStatus).
//		Scan(ctx, &v)
//
func (vsq *VulnStatusQuery) Select(fields ...string) *VulnStatusSelect {
	vsq.fields = append(vsq.fields, fields...)
	return &VulnStatusSelect{VulnStatusQuery: vsq}
}

func (vsq *VulnStatusQuery) prepareQuery(ctx context.Context) error {
	for _, f := range vsq.fields {
		if !vulnstatus.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if vsq.path != nil {
		prev, err := vsq.path(ctx)
		if err != nil {
			return err
		}
		vsq.sql = prev
	}
	return nil
}

func (vsq *VulnStatusQuery) sqlAll(ctx context.Context) ([]*VulnStatus, error) {
	var (
		nodes       = []*VulnStatus{}
		withFKs     = vsq.withFKs
		_spec       = vsq.querySpec()
		loadedTypes = [1]bool{
			vsq.withAuthor != nil,
		}
	)
	if vsq.withAuthor != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, vulnstatus.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &VulnStatus{config: vsq.config}
		nodes = append(nodes, node)
		return node.scanValues(columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if err := sqlgraph.QueryNodes(ctx, vsq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := vsq.withAuthor; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*VulnStatus)
		for i := range nodes {
			if nodes[i].vuln_status_author == nil {
				continue
			}
			fk := *nodes[i].vuln_status_author
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(user.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "vuln_status_author" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Author = n
			}
		}
	}

	return nodes, nil
}

func (vsq *VulnStatusQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := vsq.querySpec()
	return sqlgraph.CountNodes(ctx, vsq.driver, _spec)
}

func (vsq *VulnStatusQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := vsq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (vsq *VulnStatusQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   vulnstatus.Table,
			Columns: vulnstatus.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: vulnstatus.FieldID,
			},
		},
		From:   vsq.sql,
		Unique: true,
	}
	if unique := vsq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := vsq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, vulnstatus.FieldID)
		for i := range fields {
			if fields[i] != vulnstatus.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := vsq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := vsq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := vsq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := vsq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (vsq *VulnStatusQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(vsq.driver.Dialect())
	t1 := builder.Table(vulnstatus.Table)
	columns := vsq.fields
	if len(columns) == 0 {
		columns = vulnstatus.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if vsq.sql != nil {
		selector = vsq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	for _, p := range vsq.predicates {
		p(selector)
	}
	for _, p := range vsq.order {
		p(selector)
	}
	if offset := vsq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := vsq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// VulnStatusGroupBy is the group-by builder for VulnStatus entities.
type VulnStatusGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (vsgb *VulnStatusGroupBy) Aggregate(fns ...AggregateFunc) *VulnStatusGroupBy {
	vsgb.fns = append(vsgb.fns, fns...)
	return vsgb
}

// Scan applies the group-by query and scans the result into the given value.
func (vsgb *VulnStatusGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := vsgb.path(ctx)
	if err != nil {
		return err
	}
	vsgb.sql = query
	return vsgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (vsgb *VulnStatusGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := vsgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (vsgb *VulnStatusGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(vsgb.fields) > 1 {
		return nil, errors.New("ent: VulnStatusGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := vsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (vsgb *VulnStatusGroupBy) StringsX(ctx context.Context) []string {
	v, err := vsgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (vsgb *VulnStatusGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = vsgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{vulnstatus.Label}
	default:
		err = fmt.Errorf("ent: VulnStatusGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (vsgb *VulnStatusGroupBy) StringX(ctx context.Context) string {
	v, err := vsgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (vsgb *VulnStatusGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(vsgb.fields) > 1 {
		return nil, errors.New("ent: VulnStatusGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := vsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (vsgb *VulnStatusGroupBy) IntsX(ctx context.Context) []int {
	v, err := vsgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (vsgb *VulnStatusGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = vsgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{vulnstatus.Label}
	default:
		err = fmt.Errorf("ent: VulnStatusGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (vsgb *VulnStatusGroupBy) IntX(ctx context.Context) int {
	v, err := vsgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (vsgb *VulnStatusGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(vsgb.fields) > 1 {
		return nil, errors.New("ent: VulnStatusGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := vsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (vsgb *VulnStatusGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := vsgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (vsgb *VulnStatusGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = vsgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{vulnstatus.Label}
	default:
		err = fmt.Errorf("ent: VulnStatusGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (vsgb *VulnStatusGroupBy) Float64X(ctx context.Context) float64 {
	v, err := vsgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (vsgb *VulnStatusGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(vsgb.fields) > 1 {
		return nil, errors.New("ent: VulnStatusGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := vsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (vsgb *VulnStatusGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := vsgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (vsgb *VulnStatusGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = vsgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{vulnstatus.Label}
	default:
		err = fmt.Errorf("ent: VulnStatusGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (vsgb *VulnStatusGroupBy) BoolX(ctx context.Context) bool {
	v, err := vsgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (vsgb *VulnStatusGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range vsgb.fields {
		if !vulnstatus.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := vsgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := vsgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (vsgb *VulnStatusGroupBy) sqlQuery() *sql.Selector {
	selector := vsgb.sql.Select()
	aggregation := make([]string, 0, len(vsgb.fns))
	for _, fn := range vsgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(vsgb.fields)+len(vsgb.fns))
		for _, f := range vsgb.fields {
			columns = append(columns, selector.C(f))
		}
		for _, c := range aggregation {
			columns = append(columns, c)
		}
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(vsgb.fields...)...)
}

// VulnStatusSelect is the builder for selecting fields of VulnStatus entities.
type VulnStatusSelect struct {
	*VulnStatusQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (vss *VulnStatusSelect) Scan(ctx context.Context, v interface{}) error {
	if err := vss.prepareQuery(ctx); err != nil {
		return err
	}
	vss.sql = vss.VulnStatusQuery.sqlQuery(ctx)
	return vss.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (vss *VulnStatusSelect) ScanX(ctx context.Context, v interface{}) {
	if err := vss.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (vss *VulnStatusSelect) Strings(ctx context.Context) ([]string, error) {
	if len(vss.fields) > 1 {
		return nil, errors.New("ent: VulnStatusSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := vss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (vss *VulnStatusSelect) StringsX(ctx context.Context) []string {
	v, err := vss.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (vss *VulnStatusSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = vss.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{vulnstatus.Label}
	default:
		err = fmt.Errorf("ent: VulnStatusSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (vss *VulnStatusSelect) StringX(ctx context.Context) string {
	v, err := vss.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (vss *VulnStatusSelect) Ints(ctx context.Context) ([]int, error) {
	if len(vss.fields) > 1 {
		return nil, errors.New("ent: VulnStatusSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := vss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (vss *VulnStatusSelect) IntsX(ctx context.Context) []int {
	v, err := vss.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (vss *VulnStatusSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = vss.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{vulnstatus.Label}
	default:
		err = fmt.Errorf("ent: VulnStatusSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (vss *VulnStatusSelect) IntX(ctx context.Context) int {
	v, err := vss.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (vss *VulnStatusSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(vss.fields) > 1 {
		return nil, errors.New("ent: VulnStatusSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := vss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (vss *VulnStatusSelect) Float64sX(ctx context.Context) []float64 {
	v, err := vss.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (vss *VulnStatusSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = vss.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{vulnstatus.Label}
	default:
		err = fmt.Errorf("ent: VulnStatusSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (vss *VulnStatusSelect) Float64X(ctx context.Context) float64 {
	v, err := vss.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (vss *VulnStatusSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(vss.fields) > 1 {
		return nil, errors.New("ent: VulnStatusSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := vss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (vss *VulnStatusSelect) BoolsX(ctx context.Context) []bool {
	v, err := vss.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (vss *VulnStatusSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = vss.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{vulnstatus.Label}
	default:
		err = fmt.Errorf("ent: VulnStatusSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (vss *VulnStatusSelect) BoolX(ctx context.Context) bool {
	v, err := vss.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (vss *VulnStatusSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := vss.sql.Query()
	if err := vss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
