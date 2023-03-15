// Code generated by entc, DO NOT EDIT.

package vulnstatusindex

const (
	// Label holds the string label denoting the vulnstatusindex type in the database.
	Label = "vuln_status_index"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// EdgeLatest holds the string denoting the latest edge name in mutations.
	EdgeLatest = "latest"
	// EdgeStatus holds the string denoting the status edge name in mutations.
	EdgeStatus = "status"
	// Table holds the table name of the vulnstatusindex in the database.
	Table = "vuln_status_indexes"
	// LatestTable is the table that holds the latest relation/edge.
	LatestTable = "vuln_status_indexes"
	// LatestInverseTable is the table name for the VulnStatus entity.
	// It exists in this package in order to avoid circular dependency with the "vulnstatus" package.
	LatestInverseTable = "vuln_status"
	// LatestColumn is the table column denoting the latest relation/edge.
	LatestColumn = "vuln_status_index_latest"
	// StatusTable is the table that holds the status relation/edge.
	StatusTable = "vuln_status"
	// StatusInverseTable is the table name for the VulnStatus entity.
	// It exists in this package in order to avoid circular dependency with the "vulnstatus" package.
	StatusInverseTable = "vuln_status"
	// StatusColumn is the table column denoting the status relation/edge.
	StatusColumn = "vuln_status_index_status"
)

// Columns holds all SQL columns for vulnstatusindex fields.
var Columns = []string{
	FieldID,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "vuln_status_indexes"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"repository_status",
	"vuln_status_index_latest",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// IDValidator is a validator for the "id" field. It is called by the builders before save.
	IDValidator func(string) error
)