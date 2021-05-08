package model

import "github.com/m-mizutani/goerr"

type ScanRepositoryRequest struct {
	ScanTarget
	InstallID int64
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

type ScanTarget struct {
	GitHubBranch
	CommitID    string
	UpdatedAt   int64
	RequestedAt int64
}

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

type ScanResult struct {
	Target      ScanTarget
	ScannedAt   int64
	Sources     []*PackageSource
	TrivyDBMeta TrivyDBMeta
}

func (x *ScanResult) IsValid() error {
	if err := x.Target.IsValid(); err != nil {
		return err
	}
	if x.ScannedAt == 0 {
		return goerr.Wrap(ErrInvalidInputValues, "ScannedAt is not set")
	}

	return nil
}
