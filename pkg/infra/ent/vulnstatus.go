// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/ent/user"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnstatus"
)

// VulnStatus is the model entity for the VulnStatus schema.
type VulnStatus struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Status holds the value of the "status" field.
	Status types.VulnStatusType `json:"status,omitempty"`
	// Source holds the value of the "source" field.
	Source string `json:"source,omitempty"`
	// PkgName holds the value of the "pkg_name" field.
	PkgName string `json:"pkg_name,omitempty"`
	// PkgType holds the value of the "pkg_type" field.
	PkgType string `json:"pkg_type,omitempty"`
	// VulnID holds the value of the "vuln_id" field.
	VulnID string `json:"vuln_id,omitempty"`
	// ExpiresAt holds the value of the "expires_at" field.
	ExpiresAt int64 `json:"expires_at,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt int64 `json:"created_at,omitempty"`
	// Comment holds the value of the "comment" field.
	Comment string `json:"comment,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the VulnStatusQuery when eager-loading is set.
	Edges                    VulnStatusEdges `json:"edges"`
	user_edited_status       *int
	vuln_status_author       *int
	vuln_status_index_status *string
}

// VulnStatusEdges holds the relations/edges for other nodes in the graph.
type VulnStatusEdges struct {
	// Author holds the value of the author edge.
	Author *User `json:"author,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// AuthorOrErr returns the Author value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e VulnStatusEdges) AuthorOrErr() (*User, error) {
	if e.loadedTypes[0] {
		if e.Author == nil {
			// The edge author was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.Author, nil
	}
	return nil, &NotLoadedError{edge: "author"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*VulnStatus) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case vulnstatus.FieldID, vulnstatus.FieldExpiresAt, vulnstatus.FieldCreatedAt:
			values[i] = new(sql.NullInt64)
		case vulnstatus.FieldStatus, vulnstatus.FieldSource, vulnstatus.FieldPkgName, vulnstatus.FieldPkgType, vulnstatus.FieldVulnID, vulnstatus.FieldComment:
			values[i] = new(sql.NullString)
		case vulnstatus.ForeignKeys[0]: // user_edited_status
			values[i] = new(sql.NullInt64)
		case vulnstatus.ForeignKeys[1]: // vuln_status_author
			values[i] = new(sql.NullInt64)
		case vulnstatus.ForeignKeys[2]: // vuln_status_index_status
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type VulnStatus", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the VulnStatus fields.
func (vs *VulnStatus) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case vulnstatus.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			vs.ID = int(value.Int64)
		case vulnstatus.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				vs.Status = types.VulnStatusType(value.String)
			}
		case vulnstatus.FieldSource:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field source", values[i])
			} else if value.Valid {
				vs.Source = value.String
			}
		case vulnstatus.FieldPkgName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field pkg_name", values[i])
			} else if value.Valid {
				vs.PkgName = value.String
			}
		case vulnstatus.FieldPkgType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field pkg_type", values[i])
			} else if value.Valid {
				vs.PkgType = value.String
			}
		case vulnstatus.FieldVulnID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field vuln_id", values[i])
			} else if value.Valid {
				vs.VulnID = value.String
			}
		case vulnstatus.FieldExpiresAt:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field expires_at", values[i])
			} else if value.Valid {
				vs.ExpiresAt = value.Int64
			}
		case vulnstatus.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				vs.CreatedAt = value.Int64
			}
		case vulnstatus.FieldComment:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field comment", values[i])
			} else if value.Valid {
				vs.Comment = value.String
			}
		case vulnstatus.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field user_edited_status", value)
			} else if value.Valid {
				vs.user_edited_status = new(int)
				*vs.user_edited_status = int(value.Int64)
			}
		case vulnstatus.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field vuln_status_author", value)
			} else if value.Valid {
				vs.vuln_status_author = new(int)
				*vs.vuln_status_author = int(value.Int64)
			}
		case vulnstatus.ForeignKeys[2]:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field vuln_status_index_status", values[i])
			} else if value.Valid {
				vs.vuln_status_index_status = new(string)
				*vs.vuln_status_index_status = value.String
			}
		}
	}
	return nil
}

// QueryAuthor queries the "author" edge of the VulnStatus entity.
func (vs *VulnStatus) QueryAuthor() *UserQuery {
	return (&VulnStatusClient{config: vs.config}).QueryAuthor(vs)
}

// Update returns a builder for updating this VulnStatus.
// Note that you need to call VulnStatus.Unwrap() before calling this method if this VulnStatus
// was returned from a transaction, and the transaction was committed or rolled back.
func (vs *VulnStatus) Update() *VulnStatusUpdateOne {
	return (&VulnStatusClient{config: vs.config}).UpdateOne(vs)
}

// Unwrap unwraps the VulnStatus entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (vs *VulnStatus) Unwrap() *VulnStatus {
	tx, ok := vs.config.driver.(*txDriver)
	if !ok {
		panic("ent: VulnStatus is not a transactional entity")
	}
	vs.config.driver = tx.drv
	return vs
}

// String implements the fmt.Stringer.
func (vs *VulnStatus) String() string {
	var builder strings.Builder
	builder.WriteString("VulnStatus(")
	builder.WriteString(fmt.Sprintf("id=%v", vs.ID))
	builder.WriteString(", status=")
	builder.WriteString(fmt.Sprintf("%v", vs.Status))
	builder.WriteString(", source=")
	builder.WriteString(vs.Source)
	builder.WriteString(", pkg_name=")
	builder.WriteString(vs.PkgName)
	builder.WriteString(", pkg_type=")
	builder.WriteString(vs.PkgType)
	builder.WriteString(", vuln_id=")
	builder.WriteString(vs.VulnID)
	builder.WriteString(", expires_at=")
	builder.WriteString(fmt.Sprintf("%v", vs.ExpiresAt))
	builder.WriteString(", created_at=")
	builder.WriteString(fmt.Sprintf("%v", vs.CreatedAt))
	builder.WriteString(", comment=")
	builder.WriteString(vs.Comment)
	builder.WriteByte(')')
	return builder.String()
}

// VulnStatusSlice is a parsable slice of VulnStatus.
type VulnStatusSlice []*VulnStatus

func (vs VulnStatusSlice) config(cfg config) {
	for _i := range vs {
		vs[_i].config = cfg
	}
}
