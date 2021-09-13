package db

import (
	"context"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

func (x *Client) GetBranch(ctx context.Context, key *BranchKey) (*ent.Branch, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	branchID, err := x.client.Branch.Create().
		SetKey(key.Owner + "/" + key.RepoName + ":" + key.Branch).
		SetRepoName(key.RepoName).
		SetOwner(key.Owner).
		SetName(key.Branch).
		OnConflict().
		Ignore().ID(ctx)

	if err != nil {
		return nil, goerr.Wrap(err)
	}

	branch, err := x.client.Branch.Get(ctx, branchID)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	return branch, nil
}
