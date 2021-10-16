package db

import (
	"github.com/google/uuid"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/ent/repository"
	"github.com/m-mizutani/octovy/pkg/infra/ent/scan"
)

func (x *Client) PutPackages(ctx *model.Context, packages []*ent.PackageRecord) ([]*ent.PackageRecord, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	pkgBuilder := make([]*ent.PackageRecordCreate, len(packages))
	for i, pkg := range packages {
		pkgBuilder[i] = x.client.PackageRecord.Create().
			SetName(pkg.Name).
			SetSource(pkg.Source).
			SetType(pkg.Type).
			SetVersion(pkg.Version).
			SetVulnIds(pkg.VulnIds).
			AddVulnerabilityIDs(pkg.VulnIds...)
	}
	added, err := x.client.PackageRecord.CreateBulk(pkgBuilder...).Save(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return added, nil
}

func txRollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = goerr.Wrap(rerr).With("original", err)
	}
	return err
}

func (x *Client) PutScan(ctx *model.Context, scan *ent.Scan, repo *ent.Repository, packages []*ent.PackageRecord) (*ent.Scan, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}
	// logger.Debug().Interface("pkg", packages).Send()
	added, err := x.client.Scan.Create().
		SetID(uuid.NewString()).
		SetCommitID(scan.CommitID).
		SetBranch(scan.Branch).
		SetRequestedAt(scan.RequestedAt).
		SetScannedAt(scan.ScannedAt).
		SetCheckID(scan.CheckID).
		SetPullRequestTarget(scan.PullRequestTarget).
		AddRepository(repo).
		AddPackages(packages...).
		Save(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	if repo.DefaultBranch != nil && scan.Branch == *repo.DefaultBranch {
		if err := x.client.Repository.UpdateOneID(repo.ID).AddMainIDs(added.ID).Exec(ctx); err != nil {
			return nil, goerr.Wrap(err)
		}
		if err := x.client.Repository.UpdateOneID(repo.ID).SetLatestID(added.ID).Exec(ctx); err != nil {
			return nil, goerr.Wrap(err)
		}
	}

	return added, nil
}

func (x *Client) GetScan(ctx *model.Context, id string) (*ent.Scan, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	got, err := x.client.Scan.Query().Where(scan.ID(id)).
		WithRepository(func(rq *ent.RepositoryQuery) {
			rq.WithStatus(func(vsiq *ent.VulnStatusIndexQuery) {
				vsiq.WithStatus(func(vsq *ent.VulnStatusQuery) {
					vsq.Order(ent.Desc("created_at")).Limit(1)
				})
			})
		}).
		WithPackages(func(prq *ent.PackageRecordQuery) {
			prq.WithVulnerabilities()
		}).
		Only(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return got, nil
}

func (x *Client) GetLatestScan(ctx *model.Context, branch model.GitHubBranch) (*ent.Scan, error) {
	latest, err := x.getLatestScanEntity(ctx, branch)
	if err != nil {
		return nil, err
	}
	if latest == nil {
		return nil, nil
	}
	return x.GetScan(ctx, latest.ID)
}

func (x *Client) GetLatestScans(ctx *model.Context) ([]*ent.Scan, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	repos, err := x.client.Repository.Query().
		WithMain(func(sq *ent.ScanQuery) {
			sq.Order(ent.Desc("scanned_at")).Limit(1).
				WithPackages(func(prq *ent.PackageRecordQuery) {
					prq.WithVulnerabilities()
				}).
				WithRepository()
		}).
		All(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	var scans []*ent.Scan
	for _, repo := range repos {
		if len(repo.Edges.Main) > 0 {
			scans = append(scans, repo.Edges.Main[0])
		}
	}

	return scans, nil
}

func (x *Client) getLatestScanEntity(ctx *model.Context, branch model.GitHubBranch) (*ent.Scan, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	got, err := x.client.Repository.Query().
		Where(repository.Owner(branch.Owner)).
		Where(repository.Name(branch.RepoName)).
		WithScan(func(sq *ent.ScanQuery) {
			sq.Where(scan.Branch(branch.Branch)).
				Order(ent.Desc("scanned_at")).
				Limit(1)
		}).First(ctx)

	if ent.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, goerr.Wrap(err)
	}

	if len(got.Edges.Scan) != 1 {
		return nil, nil // not found
	}

	return got.Edges.Scan[0], nil
}
