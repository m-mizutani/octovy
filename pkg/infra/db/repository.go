package db

import (
	"context"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

func (x *Client) CreateRepo(ctx context.Context, repo *ent.Repository) (*ent.Repository, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	branchID, err := x.client.Repository.Create().
		SetName(repo.Name).
		SetOwner(repo.Owner).
		SetInstallID(repo.InstallID).
		OnConflict().
		Ignore().ID(ctx)

	if err != nil {
		return nil, goerr.Wrap(err)
	}

	branch, err := x.client.Repository.Get(ctx, branchID)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	return branch, nil
}
