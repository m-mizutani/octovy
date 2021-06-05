package model

import (
	"github.com/m-mizutani/goerr"
)

type ScanRepositoryRequest struct {
	ScanTarget
	InstallID int64
	Feedback  *FeedbackOptions
}

type FeedbackOptions struct {
	PullReqID     *int
	PullReqBranch string
	CheckID       *int64
}

func (x *ScanRepositoryRequest) IsValid() error {
	if err := x.ScanTarget.IsValid(); err != nil {
		return err
	}
	if x.InstallID == 0 {
		return goerr.Wrap(ErrInvalidInputValues, "InstallID must not be 0")
	}

	return nil
}

type FeedbackRequest struct {
	ReportID  string
	InstallID int64
	Options   FeedbackOptions
}

func (x *FeedbackRequest) IsValid() error {
	if x.ReportID == "" {
		return goerr.Wrap(ErrInvalidInputValues, "ReportID must not be empty")
	}
	if x.InstallID == 0 {
		return goerr.Wrap(ErrInvalidInputValues, "InstallID must not be 0")
	}
	if x.Options.PullReqID == nil && x.Options.CheckID == nil {
		return goerr.Wrap(ErrInvalidInputValues, "Either one of PullReqID and CheckSuiteID is required")
	}

	return nil
}

type ScanTarget struct {
	GitHubBranch
	CommitID       string
	UpdatedAt      int64
	RequestedAt    int64
	URL            string
	IsPullRequest  bool
	IsTargetBranch bool
}

// Value to pointer conversion
func Int64(v int64) *int64 { return &v }
func Int(v int) *int       { return &v }

func (x *ScanTarget) IsValid() error {
	if x.Branch == "" {
		return goerr.Wrap(ErrInvalidInputValues, "Branch is empty")
	}
	if x.Owner == "" {
		return goerr.Wrap(ErrInvalidInputValues, "Owner is empty")
	}
	if x.RepoName == "" {
		return goerr.Wrap(ErrInvalidInputValues, "RepoName is empty")
	}
	if x.CommitID == "" {
		return goerr.Wrap(ErrInvalidInputValues, "CommitID is empty")
	}

	return nil
}

type ScanLog struct {
	Target    ScanTarget
	ScannedAt int64
	Summary   ScanReportSummary
}

type ScanReportSummary struct {
	ReportID     string
	PkgTypes     []PkgType
	PkgCount     int64
	VulnCount    int64
	VulnPkgCount int64
}

type ScanReport struct {
	ReportID    string
	Target      ScanTarget
	ScannedAt   int64
	Sources     []*PackageSource
	TrivyDBMeta TrivyDBMeta
}

// ScanReportResponse is for API response of /scan/report
type ScanReportResponse struct {
	ScanReport
	Vulnerabilities map[string]*Vulnerability
}

func (x *ScanReport) IsValid() error {
	if x.ReportID == "" {
		return goerr.Wrap(ErrInvalidInputValues, "ID is not set")
	}
	if x.ScannedAt == 0 {
		return goerr.Wrap(ErrInvalidInputValues, "ScannedAt is not set")
	}
	if err := x.Target.IsValid(); err != nil {
		return err
	}

	return nil
}

func (x *ScanReport) ToLog() *ScanLog {
	summary := ScanReportSummary{
		ReportID: x.ReportID,
	}
	pkgTypes := map[PkgType]struct{}{}

	for _, src := range x.Sources {
		summary.PkgCount += int64(len(src.Packages))
		for _, pkg := range src.Packages {
			pkgTypes[pkg.Type] = struct{}{}
			summary.VulnCount += int64(len(pkg.Vulnerabilities))
			if len(pkg.Vulnerabilities) > 0 {
				summary.VulnPkgCount++
			}
		}
	}

	for pkgType := range pkgTypes {
		summary.PkgTypes = append(summary.PkgTypes, pkgType)
	}

	return &ScanLog{
		Target:    x.Target,
		ScannedAt: x.ScannedAt,
		Summary:   summary,
	}
}

// Vulnerabilities returns VulnID set of in the report
func (x *ScanReport) Vulnerabilities() []string {
	vulnMap := map[string]struct{}{}
	for s := range x.Sources {
		for p := range x.Sources[s].Packages {
			for _, vulnID := range x.Sources[s].Packages[p].Vulnerabilities {
				vulnMap[vulnID] = struct{}{}
			}
		}
	}

	var vulnIDs []string
	for vulnID := range vulnMap {
		vulnIDs = append(vulnIDs, vulnID)
	}
	return vulnIDs
}
