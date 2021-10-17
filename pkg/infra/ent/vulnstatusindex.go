// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnstatus"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnstatusindex"
)

// VulnStatusIndex is the model entity for the VulnStatusIndex schema.
type VulnStatusIndex struct {
	config
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the VulnStatusIndexQuery when eager-loading is set.
	Edges                    VulnStatusIndexEdges `json:"edges"`
	repository_status        *int
	vuln_status_index_latest *int
}

// VulnStatusIndexEdges holds the relations/edges for other nodes in the graph.
type VulnStatusIndexEdges struct {
	// Latest holds the value of the latest edge.
	Latest *VulnStatus `json:"latest,omitempty"`
	// Status holds the value of the status edge.
	Status []*VulnStatus `json:"status,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// LatestOrErr returns the Latest value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e VulnStatusIndexEdges) LatestOrErr() (*VulnStatus, error) {
	if e.loadedTypes[0] {
		if e.Latest == nil {
			// The edge latest was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: vulnstatus.Label}
		}
		return e.Latest, nil
	}
	return nil, &NotLoadedError{edge: "latest"}
}

// StatusOrErr returns the Status value or an error if the edge
// was not loaded in eager-loading.
func (e VulnStatusIndexEdges) StatusOrErr() ([]*VulnStatus, error) {
	if e.loadedTypes[1] {
		return e.Status, nil
	}
	return nil, &NotLoadedError{edge: "status"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*VulnStatusIndex) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case vulnstatusindex.FieldID:
			values[i] = new(sql.NullString)
		case vulnstatusindex.ForeignKeys[0]: // repository_status
			values[i] = new(sql.NullInt64)
		case vulnstatusindex.ForeignKeys[1]: // vuln_status_index_latest
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type VulnStatusIndex", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the VulnStatusIndex fields.
func (vsi *VulnStatusIndex) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case vulnstatusindex.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				vsi.ID = value.String
			}
		case vulnstatusindex.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field repository_status", value)
			} else if value.Valid {
				vsi.repository_status = new(int)
				*vsi.repository_status = int(value.Int64)
			}
		case vulnstatusindex.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field vuln_status_index_latest", value)
			} else if value.Valid {
				vsi.vuln_status_index_latest = new(int)
				*vsi.vuln_status_index_latest = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryLatest queries the "latest" edge of the VulnStatusIndex entity.
func (vsi *VulnStatusIndex) QueryLatest() *VulnStatusQuery {
	return (&VulnStatusIndexClient{config: vsi.config}).QueryLatest(vsi)
}

// QueryStatus queries the "status" edge of the VulnStatusIndex entity.
func (vsi *VulnStatusIndex) QueryStatus() *VulnStatusQuery {
	return (&VulnStatusIndexClient{config: vsi.config}).QueryStatus(vsi)
}

// Update returns a builder for updating this VulnStatusIndex.
// Note that you need to call VulnStatusIndex.Unwrap() before calling this method if this VulnStatusIndex
// was returned from a transaction, and the transaction was committed or rolled back.
func (vsi *VulnStatusIndex) Update() *VulnStatusIndexUpdateOne {
	return (&VulnStatusIndexClient{config: vsi.config}).UpdateOne(vsi)
}

// Unwrap unwraps the VulnStatusIndex entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (vsi *VulnStatusIndex) Unwrap() *VulnStatusIndex {
	tx, ok := vsi.config.driver.(*txDriver)
	if !ok {
		panic("ent: VulnStatusIndex is not a transactional entity")
	}
	vsi.config.driver = tx.drv
	return vsi
}

// String implements the fmt.Stringer.
func (vsi *VulnStatusIndex) String() string {
	var builder strings.Builder
	builder.WriteString("VulnStatusIndex(")
	builder.WriteString(fmt.Sprintf("id=%v", vsi.ID))
	builder.WriteByte(')')
	return builder.String()
}

// VulnStatusIndexes is a parsable slice of VulnStatusIndex.
type VulnStatusIndexes []*VulnStatusIndex

func (vsi VulnStatusIndexes) config(cfg config) {
	for _i := range vsi {
		vsi[_i].config = cfg
	}
}
