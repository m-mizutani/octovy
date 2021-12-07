package db

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/ent/repository"
	"github.com/m-mizutani/octovy/pkg/infra/ent/scan"
)

func (x *Client) CreateRepo(ctx *model.Context, repo *ent.Repository) (*ent.Repository, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}
	logger.With("repo", repo).Trace("starting CreateRepo")

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
	logger.With("updated", updated).Trace("done CreateRepo")

	return updated, nil
}

func (x *Client) GetRepositories(ctx *model.Context) ([]*ent.Repository, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	resp, err := x.client.Repository.Query().
		WithStatus().
		WithLabels().
		WithLatest(func(sq *ent.ScanQuery) {
			sq.WithPackages(func(prq *ent.PackageRecordQuery) {
				prq.WithVulnerabilities()
			})
		}).All(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return resp, nil
}

func (x *Client) GetRepositoriesWithVuln(ctx *model.Context, vulnID string) ([]*ent.Repository, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	repos, err := x.client.Repository.Query().
		WithStatus().
		WithLatest(func(sq *ent.ScanQuery) {
			sq.WithPackages()
		}).All(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	hasVulnID := func(vulnIDs []string) bool {
		for i := range vulnIDs {
			if vulnIDs[i] == vulnID {
				return true
			}
		}
		return false
	}

	var resp []*ent.Repository
	for _, repo := range repos {
		if repo.Edges.Latest == nil {
			continue
		}

		for _, pkg := range repo.Edges.Latest.Edges.Packages {
			if hasVulnID(pkg.VulnIds) {
				resp = append(resp, repo)
				break
			}
		}
	}

	return resp, nil
}

func (x *Client) GetRepository(ctx *model.Context, repo *model.GitHubRepo) (*ent.Repository, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	resp, err := x.client.Repository.Query().
		Where(repository.Owner(repo.Owner)).
		Where(repository.Name(repo.Name)).
		WithLabels().
		WithStatus(func(vsiq *ent.VulnStatusIndexQuery) {
			vsiq.WithLatest(func(vsq *ent.VulnStatusQuery) {
				vsq.WithAuthor()
			})
		}).
		First(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return resp, nil
}

func (x *Client) GetRepositoryScan(ctx *model.Context, req *model.GetRepoScanRequest) ([]*ent.Scan, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	resp, err := x.client.Repository.Query().
		Where(repository.Owner(req.Owner)).
		Where(repository.Name(req.Name)).
		WithScan(func(sq *ent.ScanQuery) {
			sq.Order(ent.Desc(scan.FieldScannedAt)).
				Offset(req.Offset).
				Limit(req.Limit).
				WithPackages()
		}).
		WithStatus().
		All(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return resp[0].Edges.Scan, nil
}
