// Code generated by entc, DO NOT EDIT.

package checkrule

const (
	// Label holds the string label denoting the checkrule type in the database.
	Label = "check_rule"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldResult holds the string denoting the result field in the database.
	FieldResult = "result"
	// EdgeSeverity holds the string denoting the severity edge name in mutations.
	EdgeSeverity = "severity"
	// Table holds the table name of the checkrule in the database.
	Table = "check_rules"
	// SeverityTable is the table that holds the severity relation/edge.
	SeverityTable = "check_rules"
	// SeverityInverseTable is the table name for the Severity entity.
	// It exists in this package in order to avoid circular dependency with the "severity" package.
	SeverityInverseTable = "severities"
	// SeverityColumn is the table column denoting the severity relation/edge.
	SeverityColumn = "check_rule_severity"
)

// Columns holds all SQL columns for checkrule fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldResult,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "check_rules"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"check_rule_severity",
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
