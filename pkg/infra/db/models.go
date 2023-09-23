// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type TargetClass string

const (
	TargetClassOsPkgs   TargetClass = "os-pkgs"
	TargetClassLangPkgs TargetClass = "lang-pkgs"
)

func (e *TargetClass) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TargetClass(s)
	case string:
		*e = TargetClass(s)
	default:
		return fmt.Errorf("unsupported scan type for TargetClass: %T", src)
	}
	return nil
}

type NullTargetClass struct {
	TargetClass TargetClass
	Valid       bool // Valid is true if TargetClass is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTargetClass) Scan(value interface{}) error {
	if value == nil {
		ns.TargetClass, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TargetClass.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTargetClass) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TargetClass), nil
}

type MetaGithubRepository struct {
	ID              uuid.UUID
	ScanID          uuid.UUID
	Owner           string
	RepoName        string
	CommitID        string
	Branch          sql.NullString
	IsDefaultBranch sql.NullBool
	BaseCommitID    sql.NullString
	PullRequestID   sql.NullInt32
	PageSeq         sql.NullInt32
}

type Package struct {
	ID         string
	TargetType string
	Name       string
	Version    string
}

type Result struct {
	ID         uuid.UUID
	ScanID     uuid.UUID
	Target     string
	TargetType string
	Class      TargetClass
}

type ResultPackage struct {
	ID       uuid.UUID
	ResultID uuid.UUID
	PkgID    string
}

type ResultVulnerability struct {
	ID           uuid.UUID
	ResultID     uuid.UUID
	VulnID       string
	PkgID        string
	FixedVersion sql.NullString
	PrimaryUrl   sql.NullString
}

type Scan struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	ArtifactName string
	ArtifactType string
	PageSeq      sql.NullInt32
}

type Vulnerability struct {
	ID             string
	Title          string
	Severity       string
	PublishedAt    sql.NullTime
	LastModifiedAt sql.NullTime
	Data           pqtype.NullRawMessage
	PageSeq        sql.NullInt32
}
