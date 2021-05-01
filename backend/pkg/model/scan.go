package model

type ScanTarget struct {
	GitHubBranch

	// Ref presents commitID or branch name
	Ref       string
	UpdatedAt int64
}

type PackageVersion struct {
	Name    string
	Version string
}

type PackageSource struct {
	Source   string
	PkgType  PkgType
	Packages []*PackageVersion
}

type ScanResult struct {
	Target  ScanTarget
	Sources []*PackageSource
}
