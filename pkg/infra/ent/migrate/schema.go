// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AuthStateCachesColumns holds the columns for the "auth_state_caches" table.
	AuthStateCachesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString},
		{Name: "expires_at", Type: field.TypeInt64},
	}
	// AuthStateCachesTable holds the schema information for the "auth_state_caches" table.
	AuthStateCachesTable = &schema.Table{
		Name:       "auth_state_caches",
		Columns:    AuthStateCachesColumns,
		PrimaryKey: []*schema.Column{AuthStateCachesColumns[0]},
	}
	// PackageRecordsColumns holds the columns for the "package_records" table.
	PackageRecordsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "type", Type: field.TypeEnum, Enums: []string{"rubygems", "npm", "gomod", "pypi"}},
		{Name: "source", Type: field.TypeString},
		{Name: "name", Type: field.TypeString},
		{Name: "version", Type: field.TypeString},
		{Name: "vuln_ids", Type: field.TypeJSON},
	}
	// PackageRecordsTable holds the schema information for the "package_records" table.
	PackageRecordsTable = &schema.Table{
		Name:       "package_records",
		Columns:    PackageRecordsColumns,
		PrimaryKey: []*schema.Column{PackageRecordsColumns[0]},
	}
	// RepositoriesColumns holds the columns for the "repositories" table.
	RepositoriesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "owner", Type: field.TypeString},
		{Name: "name", Type: field.TypeString},
		{Name: "install_id", Type: field.TypeInt64, Nullable: true},
		{Name: "url", Type: field.TypeString, Nullable: true},
		{Name: "avatar_url", Type: field.TypeString, Nullable: true},
		{Name: "default_branch", Type: field.TypeString, Nullable: true},
		{Name: "repository_latest", Type: field.TypeString, Nullable: true},
	}
	// RepositoriesTable holds the schema information for the "repositories" table.
	RepositoriesTable = &schema.Table{
		Name:       "repositories",
		Columns:    RepositoriesColumns,
		PrimaryKey: []*schema.Column{RepositoriesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "repositories_scans_latest",
				Columns:    []*schema.Column{RepositoriesColumns[7]},
				RefColumns: []*schema.Column{ScansColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// ScansColumns holds the columns for the "scans" table.
	ScansColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "branch", Type: field.TypeString},
		{Name: "commit_id", Type: field.TypeString},
		{Name: "requested_at", Type: field.TypeInt64},
		{Name: "scanned_at", Type: field.TypeInt64},
		{Name: "check_id", Type: field.TypeInt64, Nullable: true},
		{Name: "pull_request_target", Type: field.TypeString, Nullable: true},
		{Name: "repository_main", Type: field.TypeInt, Nullable: true},
	}
	// ScansTable holds the schema information for the "scans" table.
	ScansTable = &schema.Table{
		Name:       "scans",
		Columns:    ScansColumns,
		PrimaryKey: []*schema.Column{ScansColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "scans_repositories_main",
				Columns:    []*schema.Column{ScansColumns[7]},
				RefColumns: []*schema.Column{RepositoriesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// SessionsColumns holds the columns for the "sessions" table.
	SessionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString},
		{Name: "user_id", Type: field.TypeInt},
		{Name: "token", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeInt64},
		{Name: "expires_at", Type: field.TypeInt64},
		{Name: "session_login", Type: field.TypeInt, Nullable: true},
	}
	// SessionsTable holds the schema information for the "sessions" table.
	SessionsTable = &schema.Table{
		Name:       "sessions",
		Columns:    SessionsColumns,
		PrimaryKey: []*schema.Column{SessionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "sessions_users_login",
				Columns:    []*schema.Column{SessionsColumns[5]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "github_id", Type: field.TypeInt64, Unique: true},
		{Name: "login", Type: field.TypeString},
		{Name: "name", Type: field.TypeString},
		{Name: "avatar_url", Type: field.TypeString},
		{Name: "url", Type: field.TypeString},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// VulnStatusColumns holds the columns for the "vuln_status" table.
	VulnStatusColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"none", "snoozed", "mitigated", "unaffected", "fixed"}},
		{Name: "source", Type: field.TypeString},
		{Name: "pkg_name", Type: field.TypeString},
		{Name: "pkg_type", Type: field.TypeEnum, Enums: []string{"rubygems", "npm", "gomod", "pypi"}},
		{Name: "vuln_id", Type: field.TypeString},
		{Name: "expires_at", Type: field.TypeInt64},
		{Name: "created_at", Type: field.TypeInt64},
		{Name: "comment", Type: field.TypeString},
		{Name: "user_edited_status", Type: field.TypeInt, Nullable: true},
		{Name: "vuln_status_author", Type: field.TypeInt, Nullable: true},
		{Name: "vuln_status_index_status", Type: field.TypeString, Nullable: true},
		{Name: "vulnerability_status", Type: field.TypeString, Nullable: true},
	}
	// VulnStatusTable holds the schema information for the "vuln_status" table.
	VulnStatusTable = &schema.Table{
		Name:       "vuln_status",
		Columns:    VulnStatusColumns,
		PrimaryKey: []*schema.Column{VulnStatusColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "vuln_status_users_edited_status",
				Columns:    []*schema.Column{VulnStatusColumns[9]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "vuln_status_users_author",
				Columns:    []*schema.Column{VulnStatusColumns[10]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "vuln_status_vuln_status_indexes_status",
				Columns:    []*schema.Column{VulnStatusColumns[11]},
				RefColumns: []*schema.Column{VulnStatusIndexesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "vuln_status_vulnerabilities_status",
				Columns:    []*schema.Column{VulnStatusColumns[12]},
				RefColumns: []*schema.Column{VulnerabilitiesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// VulnStatusIndexesColumns holds the columns for the "vuln_status_indexes" table.
	VulnStatusIndexesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "repository_status", Type: field.TypeInt, Nullable: true},
	}
	// VulnStatusIndexesTable holds the schema information for the "vuln_status_indexes" table.
	VulnStatusIndexesTable = &schema.Table{
		Name:       "vuln_status_indexes",
		Columns:    VulnStatusIndexesColumns,
		PrimaryKey: []*schema.Column{VulnStatusIndexesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "vuln_status_indexes_repositories_status",
				Columns:    []*schema.Column{VulnStatusIndexesColumns[1]},
				RefColumns: []*schema.Column{RepositoriesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// VulnerabilitiesColumns holds the columns for the "vulnerabilities" table.
	VulnerabilitiesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "first_seen_at", Type: field.TypeInt64},
		{Name: "last_modified_at", Type: field.TypeInt64},
		{Name: "title", Type: field.TypeString, Nullable: true},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "cwe_id", Type: field.TypeJSON, Nullable: true},
		{Name: "severity", Type: field.TypeString, Nullable: true},
		{Name: "cvss", Type: field.TypeJSON, Nullable: true},
		{Name: "references", Type: field.TypeJSON, Nullable: true},
	}
	// VulnerabilitiesTable holds the schema information for the "vulnerabilities" table.
	VulnerabilitiesTable = &schema.Table{
		Name:       "vulnerabilities",
		Columns:    VulnerabilitiesColumns,
		PrimaryKey: []*schema.Column{VulnerabilitiesColumns[0]},
	}
	// PackageRecordVulnerabilitiesColumns holds the columns for the "package_record_vulnerabilities" table.
	PackageRecordVulnerabilitiesColumns = []*schema.Column{
		{Name: "package_record_id", Type: field.TypeInt},
		{Name: "vulnerability_id", Type: field.TypeString},
	}
	// PackageRecordVulnerabilitiesTable holds the schema information for the "package_record_vulnerabilities" table.
	PackageRecordVulnerabilitiesTable = &schema.Table{
		Name:       "package_record_vulnerabilities",
		Columns:    PackageRecordVulnerabilitiesColumns,
		PrimaryKey: []*schema.Column{PackageRecordVulnerabilitiesColumns[0], PackageRecordVulnerabilitiesColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "package_record_vulnerabilities_package_record_id",
				Columns:    []*schema.Column{PackageRecordVulnerabilitiesColumns[0]},
				RefColumns: []*schema.Column{PackageRecordsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "package_record_vulnerabilities_vulnerability_id",
				Columns:    []*schema.Column{PackageRecordVulnerabilitiesColumns[1]},
				RefColumns: []*schema.Column{VulnerabilitiesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// RepositoryScanColumns holds the columns for the "repository_scan" table.
	RepositoryScanColumns = []*schema.Column{
		{Name: "repository_id", Type: field.TypeInt},
		{Name: "scan_id", Type: field.TypeString},
	}
	// RepositoryScanTable holds the schema information for the "repository_scan" table.
	RepositoryScanTable = &schema.Table{
		Name:       "repository_scan",
		Columns:    RepositoryScanColumns,
		PrimaryKey: []*schema.Column{RepositoryScanColumns[0], RepositoryScanColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "repository_scan_repository_id",
				Columns:    []*schema.Column{RepositoryScanColumns[0]},
				RefColumns: []*schema.Column{RepositoriesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "repository_scan_scan_id",
				Columns:    []*schema.Column{RepositoryScanColumns[1]},
				RefColumns: []*schema.Column{ScansColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// ScanPackagesColumns holds the columns for the "scan_packages" table.
	ScanPackagesColumns = []*schema.Column{
		{Name: "scan_id", Type: field.TypeString},
		{Name: "package_record_id", Type: field.TypeInt},
	}
	// ScanPackagesTable holds the schema information for the "scan_packages" table.
	ScanPackagesTable = &schema.Table{
		Name:       "scan_packages",
		Columns:    ScanPackagesColumns,
		PrimaryKey: []*schema.Column{ScanPackagesColumns[0], ScanPackagesColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "scan_packages_scan_id",
				Columns:    []*schema.Column{ScanPackagesColumns[0]},
				RefColumns: []*schema.Column{ScansColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "scan_packages_package_record_id",
				Columns:    []*schema.Column{ScanPackagesColumns[1]},
				RefColumns: []*schema.Column{PackageRecordsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AuthStateCachesTable,
		PackageRecordsTable,
		RepositoriesTable,
		ScansTable,
		SessionsTable,
		UsersTable,
		VulnStatusTable,
		VulnStatusIndexesTable,
		VulnerabilitiesTable,
		PackageRecordVulnerabilitiesTable,
		RepositoryScanTable,
		ScanPackagesTable,
	}
)

func init() {
	RepositoriesTable.ForeignKeys[0].RefTable = ScansTable
	ScansTable.ForeignKeys[0].RefTable = RepositoriesTable
	SessionsTable.ForeignKeys[0].RefTable = UsersTable
	VulnStatusTable.ForeignKeys[0].RefTable = UsersTable
	VulnStatusTable.ForeignKeys[1].RefTable = UsersTable
	VulnStatusTable.ForeignKeys[2].RefTable = VulnStatusIndexesTable
	VulnStatusTable.ForeignKeys[3].RefTable = VulnerabilitiesTable
	VulnStatusIndexesTable.ForeignKeys[0].RefTable = RepositoriesTable
	PackageRecordVulnerabilitiesTable.ForeignKeys[0].RefTable = PackageRecordsTable
	PackageRecordVulnerabilitiesTable.ForeignKeys[1].RefTable = VulnerabilitiesTable
	RepositoryScanTable.ForeignKeys[0].RefTable = RepositoriesTable
	RepositoryScanTable.ForeignKeys[1].RefTable = ScansTable
	ScanPackagesTable.ForeignKeys[0].RefTable = ScansTable
	ScanPackagesTable.ForeignKeys[1].RefTable = PackageRecordsTable
}
