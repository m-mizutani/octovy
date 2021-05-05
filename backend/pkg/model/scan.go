package model

import "github.com/m-mizutani/goerr"

type ScanRepositoryRequest struct {
	ScanTarget
	InstallID int64
}

func (x *ScanRepositoryRequest) IsValid() error {
	if x.Branch == "" {
		return goerr.Wrap(ErrInvalidScanRequest, "Branch is empty")
	}
	if x.Owner == "" {
		return goerr.Wrap(ErrInvalidScanRequest, "Owner is empty")
	}
	if x.RepoName == "" {
		return goerr.Wrap(ErrInvalidScanRequest, "RepoName is empty")
	}
	if x.CommitID == "" {
		return goerr.Wrap(ErrInvalidScanRequest, "Ref is empty")
	}
	if x.InstallID == 0 {
		return goerr.Wrap(ErrInvalidScanRequest, "InstallID must not be 0")
	}

	return nil
}

type ScanTarget struct {
	GitHubBranch
	CommitID    string
	UpdatedAt   int64
	RequestedAt int64
}

type ScanResult struct {
	Target      ScanTarget
	ScannedAt   int64
	Sources     []*PackageSource
	TrivyDBMeta TrivyDBMeta
}
