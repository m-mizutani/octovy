package model

import "github.com/aquasecurity/trivy-db/pkg/types"

type PkgType string

const (
	PkgRubyGems PkgType = "rubygems"
	PkgNPM      PkgType = "npm"
	PkgGoModule PkgType = "gomod"
	PkgPyPI     PkgType = "pypi"
)

type PackageRecord struct {
	Detected ScanTarget
	// File path of lock file
	Source string
	Package
	ScannedAt int64
	Removed   bool
}

type Package struct {
	Type            PkgType
	Name            string
	Version         string
	Vulnerabilities []string
}

type PackageSource struct {
	Source   string
	Packages []*Package
}

type Vulnerability struct {
	VulnID         string
	Detail         types.Vulnerability
	FirstSeenAt    int64
	LastModifiedAt int64
}
