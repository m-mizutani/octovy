package types

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
