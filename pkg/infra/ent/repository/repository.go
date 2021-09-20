// Code generated by entc, DO NOT EDIT.

package repository

const (
	// Label holds the string label denoting the repository type in the database.
	Label = "repository"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldOwner holds the string denoting the owner field in the database.
	FieldOwner = "owner"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldInstallID holds the string denoting the install_id field in the database.
	FieldInstallID = "install_id"
	// FieldURL holds the string denoting the url field in the database.
	FieldURL = "url"
	// FieldAvatarURL holds the string denoting the avatar_url field in the database.
	FieldAvatarURL = "avatar_url"
	// FieldDefaultBranch holds the string denoting the default_branch field in the database.
	FieldDefaultBranch = "default_branch"
	// EdgeScan holds the string denoting the scan edge name in mutations.
	EdgeScan = "scan"
	// EdgeStatus holds the string denoting the status edge name in mutations.
	EdgeStatus = "status"
	// Table holds the table name of the repository in the database.
	Table = "repositories"
	// ScanTable is the table that holds the scan relation/edge. The primary key declared below.
	ScanTable = "repository_scan"
	// ScanInverseTable is the table name for the Scan entity.
	// It exists in this package in order to avoid circular dependency with the "scan" package.
	ScanInverseTable = "scans"
	// StatusTable is the table that holds the status relation/edge.
	StatusTable = "vuln_status"
	// StatusInverseTable is the table name for the VulnStatus entity.
	// It exists in this package in order to avoid circular dependency with the "vulnstatus" package.
	StatusInverseTable = "vuln_status"
	// StatusColumn is the table column denoting the status relation/edge.
	StatusColumn = "repository_status"
)

// Columns holds all SQL columns for repository fields.
var Columns = []string{
	FieldID,
	FieldOwner,
	FieldName,
	FieldInstallID,
	FieldURL,
	FieldAvatarURL,
	FieldDefaultBranch,
}

var (
	// ScanPrimaryKey and ScanColumn2 are the table columns denoting the
	// primary key for the scan relation (M2M).
	ScanPrimaryKey = []string{"repository_id", "scan_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}
