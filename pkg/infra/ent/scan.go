// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/m-mizutani/octovy/pkg/infra/ent/scan"
)

// Scan is the model entity for the Scan schema.
type Scan struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// Branch holds the value of the "branch" field.
	Branch string `json:"branch,omitempty"`
	// CommitID holds the value of the "commit_id" field.
	CommitID string `json:"commit_id,omitempty"`
	// RequestedAt holds the value of the "requested_at" field.
	RequestedAt int64 `json:"requested_at,omitempty"`
	// ScannedAt holds the value of the "scanned_at" field.
	ScannedAt int64 `json:"scanned_at,omitempty"`
	// CheckID holds the value of the "check_id" field.
	CheckID int64 `json:"check_id,omitempty"`
	// PullRequestTarget holds the value of the "pull_request_target" field.
	PullRequestTarget string `json:"pull_request_target,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ScanQuery when eager-loading is set.
	Edges           ScanEdges `json:"edges"`
	repository_main *int
}

// ScanEdges holds the relations/edges for other nodes in the graph.
type ScanEdges struct {
	// Repository holds the value of the repository edge.
	Repository []*Repository `json:"repository,omitempty"`
	// Packages holds the value of the packages edge.
	Packages []*PackageRecord `json:"packages,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// RepositoryOrErr returns the Repository value or an error if the edge
// was not loaded in eager-loading.
func (e ScanEdges) RepositoryOrErr() ([]*Repository, error) {
	if e.loadedTypes[0] {
		return e.Repository, nil
	}
	return nil, &NotLoadedError{edge: "repository"}
}

// PackagesOrErr returns the Packages value or an error if the edge
// was not loaded in eager-loading.
func (e ScanEdges) PackagesOrErr() ([]*PackageRecord, error) {
	if e.loadedTypes[1] {
		return e.Packages, nil
	}
	return nil, &NotLoadedError{edge: "packages"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Scan) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case scan.FieldRequestedAt, scan.FieldScannedAt, scan.FieldCheckID:
			values[i] = new(sql.NullInt64)
		case scan.FieldID, scan.FieldBranch, scan.FieldCommitID, scan.FieldPullRequestTarget:
			values[i] = new(sql.NullString)
		case scan.ForeignKeys[0]: // repository_main
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Scan", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Scan fields.
func (s *Scan) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case scan.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				s.ID = value.String
			}
		case scan.FieldBranch:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field branch", values[i])
			} else if value.Valid {
				s.Branch = value.String
			}
		case scan.FieldCommitID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field commit_id", values[i])
			} else if value.Valid {
				s.CommitID = value.String
			}
		case scan.FieldRequestedAt:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field requested_at", values[i])
			} else if value.Valid {
				s.RequestedAt = value.Int64
			}
		case scan.FieldScannedAt:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field scanned_at", values[i])
			} else if value.Valid {
				s.ScannedAt = value.Int64
			}
		case scan.FieldCheckID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field check_id", values[i])
			} else if value.Valid {
				s.CheckID = value.Int64
			}
		case scan.FieldPullRequestTarget:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field pull_request_target", values[i])
			} else if value.Valid {
				s.PullRequestTarget = value.String
			}
		case scan.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field repository_main", value)
			} else if value.Valid {
				s.repository_main = new(int)
				*s.repository_main = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryRepository queries the "repository" edge of the Scan entity.
func (s *Scan) QueryRepository() *RepositoryQuery {
	return (&ScanClient{config: s.config}).QueryRepository(s)
}

// QueryPackages queries the "packages" edge of the Scan entity.
func (s *Scan) QueryPackages() *PackageRecordQuery {
	return (&ScanClient{config: s.config}).QueryPackages(s)
}

// Update returns a builder for updating this Scan.
// Note that you need to call Scan.Unwrap() before calling this method if this Scan
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Scan) Update() *ScanUpdateOne {
	return (&ScanClient{config: s.config}).UpdateOne(s)
}

// Unwrap unwraps the Scan entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (s *Scan) Unwrap() *Scan {
	tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("ent: Scan is not a transactional entity")
	}
	s.config.driver = tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Scan) String() string {
	var builder strings.Builder
	builder.WriteString("Scan(")
	builder.WriteString(fmt.Sprintf("id=%v", s.ID))
	builder.WriteString(", branch=")
	builder.WriteString(s.Branch)
	builder.WriteString(", commit_id=")
	builder.WriteString(s.CommitID)
	builder.WriteString(", requested_at=")
	builder.WriteString(fmt.Sprintf("%v", s.RequestedAt))
	builder.WriteString(", scanned_at=")
	builder.WriteString(fmt.Sprintf("%v", s.ScannedAt))
	builder.WriteString(", check_id=")
	builder.WriteString(fmt.Sprintf("%v", s.CheckID))
	builder.WriteString(", pull_request_target=")
	builder.WriteString(s.PullRequestTarget)
	builder.WriteByte(')')
	return builder.String()
}

// Scans is a parsable slice of Scan.
type Scans []*Scan

func (s Scans) config(cfg config) {
	for _i := range s {
		s[_i].config = cfg
	}
}
