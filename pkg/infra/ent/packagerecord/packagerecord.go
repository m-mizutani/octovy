// Code generated by entc, DO NOT EDIT.

package packagerecord

const (
	// Label holds the string label denoting the packagerecord type in the database.
	Label = "package_record"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldSource holds the string denoting the source field in the database.
	FieldSource = "source"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldVersion holds the string denoting the version field in the database.
	FieldVersion = "version"
	// FieldVulnIds holds the string denoting the vuln_ids field in the database.
	FieldVulnIds = "vuln_ids"
	// EdgeScan holds the string denoting the scan edge name in mutations.
	EdgeScan = "scan"
	// EdgeVulnerabilities holds the string denoting the vulnerabilities edge name in mutations.
	EdgeVulnerabilities = "vulnerabilities"
	// Table holds the table name of the packagerecord in the database.
	Table = "package_records"
	// ScanTable is the table that holds the scan relation/edge. The primary key declared below.
	ScanTable = "scan_packages"
	// ScanInverseTable is the table name for the Scan entity.
	// It exists in this package in order to avoid circular dependency with the "scan" package.
	ScanInverseTable = "scans"
	// VulnerabilitiesTable is the table that holds the vulnerabilities relation/edge. The primary key declared below.
	VulnerabilitiesTable = "package_record_vulnerabilities"
	// VulnerabilitiesInverseTable is the table name for the Vulnerability entity.
	// It exists in this package in order to avoid circular dependency with the "vulnerability" package.
	VulnerabilitiesInverseTable = "vulnerabilities"
)

// Columns holds all SQL columns for packagerecord fields.
var Columns = []string{
	FieldID,
	FieldType,
	FieldSource,
	FieldName,
	FieldVersion,
	FieldVulnIds,
}

var (
	// ScanPrimaryKey and ScanColumn2 are the table columns denoting the
	// primary key for the scan relation (M2M).
	ScanPrimaryKey = []string{"scan_id", "package_record_id"}
	// VulnerabilitiesPrimaryKey and VulnerabilitiesColumn2 are the table columns denoting the
	// primary key for the vulnerabilities relation (M2M).
	VulnerabilitiesPrimaryKey = []string{"package_record_id", "vulnerability_id"}
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
