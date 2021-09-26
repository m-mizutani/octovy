package types

type PkgType string

const (
	PkgRubyGems PkgType = "rubygems"
	PkgNPM      PkgType = "npm"
	PkgGoModule PkgType = "gomod"
	PkgPyPI     PkgType = "pypi"
)

func (x PkgType) Values() []string {
	return []string{
		string(PkgRubyGems),
		string(PkgNPM),
		string(PkgGoModule),
		string(PkgPyPI),
	}
}

type VulnStatusType string

const (
	StatusNone       VulnStatusType = "none"
	StatusSnoozed    VulnStatusType = "snoozed"
	StatusMitigated  VulnStatusType = "mitigated"
	StatusUnaffected VulnStatusType = "unaffected"
	StatusFixed      VulnStatusType = "fixed"
)

func (x VulnStatusType) Values() []string {
	return []string{
		string(StatusNone),
		string(StatusSnoozed),
		string(StatusMitigated),
		string(StatusUnaffected),
		string(StatusFixed),
	}
}
