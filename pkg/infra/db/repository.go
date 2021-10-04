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
	logger.Debug().Interface("repo", repo).Send()

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

	q := x.client.Repository.UpdateOneID(repoID)
	if repo.InstallID != 0 {
		q = q.SetInstallID(repo.InstallID)
	}
	if repo.URL != "" {
		q = q.SetURL(repo.URL)
	}
	if repo.DefaultBranch != nil {
		q = q.SetDefaultBranch(*repo.DefaultBranch)
	}
	if repo.AvatarURL != nil {
		q = q.SetAvatarURL(*repo.AvatarURL)
	}

	updated, err := q.Save(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	logger.Debug().Interface("updated", updated).Msg("done CreateRepo")

	return updated, nil
}

func (x *Client) GetRepositories(ctx context.Context) ([]*ent.Repository, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	allRepos, err := x.client.Repository.Query().All(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	resp := make([]*ent.Repository, len(allRepos))
	for i := range allRepos {
		got, err := x.client.Repository.Query().
			Where(repository.ID(allRepos[i].ID)).
			WithMain(func(sq *ent.ScanQuery) {
				sq.Order(ent.Desc("scanned_at")).Limit(1).
					WithPackages(func(prq *ent.PackageRecordQuery) {
						prq.WithVulnerabilities()
					})
			}).
			WithScan(func(sq *ent.ScanQuery) {
				sq.Order(ent.Desc("scanned_at")).Limit(1).
					WithPackages(func(prq *ent.PackageRecordQuery) {
						prq.WithVulnerabilities()
					})
			}).
			All(ctx)
		if err != nil {
			return nil, goerr.Wrap(err)
		}
		resp[i] = got[0]
	}

	return resp, nil
}
