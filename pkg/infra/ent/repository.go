// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/m-mizutani/octovy/pkg/infra/ent/repository"
	"github.com/m-mizutani/octovy/pkg/infra/ent/scan"
)

// Repository is the model entity for the Repository schema.
type Repository struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Owner holds the value of the "owner" field.
	Owner string `json:"owner,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// InstallID holds the value of the "install_id" field.
	InstallID int64 `json:"install_id,omitempty"`
	// URL holds the value of the "url" field.
	URL string `json:"url,omitempty"`
	// AvatarURL holds the value of the "avatar_url" field.
	AvatarURL *string `json:"avatar_url,omitempty"`
	// DefaultBranch holds the value of the "default_branch" field.
	DefaultBranch *string `json:"default_branch,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the RepositoryQuery when eager-loading is set.
	Edges             RepositoryEdges `json:"edges"`
	repository_latest *string
}

// RepositoryEdges holds the relations/edges for other nodes in the graph.
type RepositoryEdges struct {
	// Scan holds the value of the scan edge.
	Scan []*Scan `json:"scan,omitempty"`
	// Main holds the value of the main edge.
	Main []*Scan `json:"main,omitempty"`
	// Latest holds the value of the latest edge.
	Latest *Scan `json:"latest,omitempty"`
	// Report holds the value of the report edge.
	Report []*Report `json:"report,omitempty"`
	// LatestReport holds the value of the latest_report edge.
	LatestReport []*Report `json:"latest_report,omitempty"`
	// Status holds the value of the status edge.
	Status []*VulnStatusIndex `json:"status,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [6]bool
}

// ScanOrErr returns the Scan value or an error if the edge
// was not loaded in eager-loading.
func (e RepositoryEdges) ScanOrErr() ([]*Scan, error) {
	if e.loadedTypes[0] {
		return e.Scan, nil
	}
	return nil, &NotLoadedError{edge: "scan"}
}

// MainOrErr returns the Main value or an error if the edge
// was not loaded in eager-loading.
func (e RepositoryEdges) MainOrErr() ([]*Scan, error) {
	if e.loadedTypes[1] {
		return e.Main, nil
	}
	return nil, &NotLoadedError{edge: "main"}
}

// LatestOrErr returns the Latest value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e RepositoryEdges) LatestOrErr() (*Scan, error) {
	if e.loadedTypes[2] {
		if e.Latest == nil {
			// The edge latest was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: scan.Label}
		}
		return e.Latest, nil
	}
	return nil, &NotLoadedError{edge: "latest"}
}

// ReportOrErr returns the Report value or an error if the edge
// was not loaded in eager-loading.
func (e RepositoryEdges) ReportOrErr() ([]*Report, error) {
	if e.loadedTypes[3] {
		return e.Report, nil
	}
	return nil, &NotLoadedError{edge: "report"}
}

// LatestReportOrErr returns the LatestReport value or an error if the edge
// was not loaded in eager-loading.
func (e RepositoryEdges) LatestReportOrErr() ([]*Report, error) {
	if e.loadedTypes[4] {
		return e.LatestReport, nil
	}
	return nil, &NotLoadedError{edge: "latest_report"}
}

// StatusOrErr returns the Status value or an error if the edge
// was not loaded in eager-loading.
func (e RepositoryEdges) StatusOrErr() ([]*VulnStatusIndex, error) {
	if e.loadedTypes[5] {
		return e.Status, nil
	}
	return nil, &NotLoadedError{edge: "status"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Repository) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case repository.FieldID, repository.FieldInstallID:
			values[i] = new(sql.NullInt64)
		case repository.FieldOwner, repository.FieldName, repository.FieldURL, repository.FieldAvatarURL, repository.FieldDefaultBranch:
			values[i] = new(sql.NullString)
		case repository.ForeignKeys[0]: // repository_latest
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Repository", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Repository fields.
func (r *Repository) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case repository.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			r.ID = int(value.Int64)
		case repository.FieldOwner:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field owner", values[i])
			} else if value.Valid {
				r.Owner = value.String
			}
		case repository.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				r.Name = value.String
			}
		case repository.FieldInstallID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field install_id", values[i])
			} else if value.Valid {
				r.InstallID = value.Int64
			}
		case repository.FieldURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field url", values[i])
			} else if value.Valid {
				r.URL = value.String
			}
		case repository.FieldAvatarURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field avatar_url", values[i])
			} else if value.Valid {
				r.AvatarURL = new(string)
				*r.AvatarURL = value.String
			}
		case repository.FieldDefaultBranch:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field default_branch", values[i])
			} else if value.Valid {
				r.DefaultBranch = new(string)
				*r.DefaultBranch = value.String
			}
		case repository.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field repository_latest", values[i])
			} else if value.Valid {
				r.repository_latest = new(string)
				*r.repository_latest = value.String
			}
		}
	}
	return nil
}

// QueryScan queries the "scan" edge of the Repository entity.
func (r *Repository) QueryScan() *ScanQuery {
	return (&RepositoryClient{config: r.config}).QueryScan(r)
}

// QueryMain queries the "main" edge of the Repository entity.
func (r *Repository) QueryMain() *ScanQuery {
	return (&RepositoryClient{config: r.config}).QueryMain(r)
}

// QueryLatest queries the "latest" edge of the Repository entity.
func (r *Repository) QueryLatest() *ScanQuery {
	return (&RepositoryClient{config: r.config}).QueryLatest(r)
}

// QueryReport queries the "report" edge of the Repository entity.
func (r *Repository) QueryReport() *ReportQuery {
	return (&RepositoryClient{config: r.config}).QueryReport(r)
}

// QueryLatestReport queries the "latest_report" edge of the Repository entity.
func (r *Repository) QueryLatestReport() *ReportQuery {
	return (&RepositoryClient{config: r.config}).QueryLatestReport(r)
}

// QueryStatus queries the "status" edge of the Repository entity.
func (r *Repository) QueryStatus() *VulnStatusIndexQuery {
	return (&RepositoryClient{config: r.config}).QueryStatus(r)
}

// Update returns a builder for updating this Repository.
// Note that you need to call Repository.Unwrap() before calling this method if this Repository
// was returned from a transaction, and the transaction was committed or rolled back.
func (r *Repository) Update() *RepositoryUpdateOne {
	return (&RepositoryClient{config: r.config}).UpdateOne(r)
}

// Unwrap unwraps the Repository entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (r *Repository) Unwrap() *Repository {
	tx, ok := r.config.driver.(*txDriver)
	if !ok {
		panic("ent: Repository is not a transactional entity")
	}
	r.config.driver = tx.drv
	return r
}

// String implements the fmt.Stringer.
func (r *Repository) String() string {
	var builder strings.Builder
	builder.WriteString("Repository(")
	builder.WriteString(fmt.Sprintf("id=%v", r.ID))
	builder.WriteString(", owner=")
	builder.WriteString(r.Owner)
	builder.WriteString(", name=")
	builder.WriteString(r.Name)
	builder.WriteString(", install_id=")
	builder.WriteString(fmt.Sprintf("%v", r.InstallID))
	builder.WriteString(", url=")
	builder.WriteString(r.URL)
	if v := r.AvatarURL; v != nil {
		builder.WriteString(", avatar_url=")
		builder.WriteString(*v)
	}
	if v := r.DefaultBranch; v != nil {
		builder.WriteString(", default_branch=")
		builder.WriteString(*v)
	}
	builder.WriteByte(')')
	return builder.String()
}

// Repositories is a parsable slice of Repository.
type Repositories []*Repository

func (r Repositories) config(cfg config) {
	for _i := range r {
		r[_i].config = cfg
	}
}
