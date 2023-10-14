package model

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/types"
)

type GitHubRepo struct {
	RepoID   int64
	Owner    string
	RepoName string
}

func (x *GitHubRepo) Validate() error {
	if x.RepoID == 0 {
		return goerr.Wrap(types.ErrInvalidOption, "repo ID is empty")
	}
	if x.Owner == "" {
		return goerr.Wrap(types.ErrInvalidOption, "owner is empty")
	}
	if x.RepoName == "" {
		return goerr.Wrap(types.ErrInvalidOption, "repo name is empty")
	}

	return nil
}

type GitHubCommit struct {
	GitHubRepo
	CommitID string
}

func (x *GitHubCommit) Validate() error {
	if err := x.GitHubRepo.Validate(); err != nil {
		return err
	}
	if x.CommitID == "" {
		return goerr.Wrap(types.ErrInvalidOption, "commit ID is empty")
	}

	return nil
}
