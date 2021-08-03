package model

import "github.com/m-mizutani/goerr"

type GitHubRepo struct {
	Owner    string
	RepoName string
}

func (x *GitHubRepo) IsValid() error {
	if x.Owner == "" {
		return goerr.Wrap(ErrInvalidValue, "Owner is not set")
	}
	if x.RepoName == "" {
		return goerr.Wrap(ErrInvalidValue, "RepoName is not set")
	}

	return nil
}

type GitHubBranch struct {
	GitHubRepo
	Branch string
}

func (x *GitHubBranch) IsValid() error {
	if x.Branch == "" {
		return goerr.Wrap(ErrInvalidValue, "Branch is not set")
	}
	if err := x.GitHubRepo.IsValid(); err != nil {
		return err
	}

	return nil
}

type GitHubCommit struct {
	GitHubRepo
	CommitID string
}

type Owner struct {
	Name      string
	AvatarURL string
}

type Repository struct {
	GitHubRepo
	URL            string
	OwnerAvatarURL string
	DefaultBranch  string
	Branch         Branch
	InstallID      int64
}

type Branch struct {
	GitHubBranch
	LastScannedAt int64
	ReportSummary ScanReportSummary
}
