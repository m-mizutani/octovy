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
	// EdgeMain holds the string denoting the main edge name in mutations.
	EdgeMain = "main"
	// EdgeLatest holds the string denoting the latest edge name in mutations.
	EdgeLatest = "latest"
	// EdgeStatus holds the string denoting the status edge name in mutations.
	EdgeStatus = "status"
	// EdgeLabels holds the string denoting the labels edge name in mutations.
	EdgeLabels = "labels"
	// Table holds the table name of the repository in the database.
	Table = "repositories"
	// ScanTable is the table that holds the scan relation/edge. The primary key declared below.
	ScanTable = "repository_scan"
	// ScanInverseTable is the table name for the Scan entity.
	// It exists in this package in order to avoid circular dependency with the "scan" package.
	ScanInverseTable = "scans"
	// MainTable is the table that holds the main relation/edge.
	MainTable = "scans"
	// MainInverseTable is the table name for the Scan entity.
	// It exists in this package in order to avoid circular dependency with the "scan" package.
	MainInverseTable = "scans"
	// MainColumn is the table column denoting the main relation/edge.
	MainColumn = "repository_main"
	// LatestTable is the table that holds the latest relation/edge.
	LatestTable = "repositories"
	// LatestInverseTable is the table name for the Scan entity.
	// It exists in this package in order to avoid circular dependency with the "scan" package.
	LatestInverseTable = "scans"
	// LatestColumn is the table column denoting the latest relation/edge.
	LatestColumn = "repository_latest"
	// StatusTable is the table that holds the status relation/edge.
	StatusTable = "vuln_status_indexes"
	// StatusInverseTable is the table name for the VulnStatusIndex entity.
	// It exists in this package in order to avoid circular dependency with the "vulnstatusindex" package.
	StatusInverseTable = "vuln_status_indexes"
	// StatusColumn is the table column denoting the status relation/edge.
	StatusColumn = "repository_status"
	// LabelsTable is the table that holds the labels relation/edge. The primary key declared below.
	LabelsTable = "repository_labels"
	// LabelsInverseTable is the table name for the RepoLabel entity.
	// It exists in this package in order to avoid circular dependency with the "repolabel" package.
	LabelsInverseTable = "repo_labels"
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

// ForeignKeys holds the SQL foreign-keys that are owned by the "repositories"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"repository_latest",
}

var (
	// ScanPrimaryKey and ScanColumn2 are the table columns denoting the
	// primary key for the scan relation (M2M).
	ScanPrimaryKey = []string{"repository_id", "scan_id"}
	// LabelsPrimaryKey and LabelsColumn2 are the table columns denoting the
	// primary key for the labels relation (M2M).
	LabelsPrimaryKey = []string{"repository_id", "repo_label_id"}
)

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
