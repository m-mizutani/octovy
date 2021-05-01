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
	if x.Ref == "" {
		return goerr.Wrap(ErrInvalidScanRequest, "Ref is empty")
	}
	if x.InstallID == 0 {
		return goerr.Wrap(ErrInvalidScanRequest, "InstallID must not be 0")
	}

	return nil
}
