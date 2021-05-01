package model

type PkgType string

const (
	PkgBundler  PkgType = "bundler"
	PkgNPM      PkgType = "npm"
	PkgYarn     PkgType = "yarn"
	PkgGoModule PkgType = "gomod"
	PkgPipenv   PkgType = "pipenv"
)

type Package struct {
	ScanTarget

	// File path of lock file
	Source  string
	PkgType PkgType
	PkgName string
	Version string
}
