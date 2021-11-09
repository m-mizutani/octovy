package types

import "github.com/m-mizutani/goerr"

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

type GitHubCheckResult string

const (
	CheckFail    GitHubCheckResult = "fail"
	CheckNeutral GitHubCheckResult = "neutral"
	CheckSuccess GitHubCheckResult = "success"
)

func (x GitHubCheckResult) IsValid() error {
	switch x {
	case CheckFail, CheckNeutral, CheckSuccess:
		return nil
	default:
		return goerr.New("invalid check result")
	}
}
