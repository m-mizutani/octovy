package db

import (
	"context"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/ent/repository"
)

func (x *Client) CreateRepo(ctx context.Context, repo *ent.Repository) (*ent.Repository, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	repoID, err := x.client.Repository.Query().
		Where(repository.Owner(repo.Owner)).
		Where(repository.Name(repo.Name)).
		FirstID(ctx)

	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, goerr.Wrap(err)
		}

		newRepo, err := x.client.Repository.Create().
			SetName(repo.Name).
			SetOwner(repo.Owner).
			Save(ctx)
		if err != nil {
			return nil, goerr.Wrap(err)
		}
		repoID = newRepo.ID
	}

	q := x.client.Repository.UpdateOneID(repoID).
		SetInstallID(repo.InstallID).
		SetURL(repo.URL)
	if repo.DefaultBranch != nil {
		q = q.SetDefaultBranch(*repo.DefaultBranch)
	}
	if repo.AvatarURL != nil {
		q = q.SetAvatarURL(*repo.AvatarURL)
	}

	if _, err := q.Save(ctx); err != nil {
		return nil, goerr.Wrap(err)
	}

	updated, err := x.client.Repository.Get(ctx, repoID)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	return updated, nil
}
