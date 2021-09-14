package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/ent/scan"
)

func (x *Client) PutPackages(ctx context.Context, packages []*ent.PackageRecord, vulnIDs []string) ([]*ent.PackageRecord, error) {
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
			AddVulnerabilityIDs(vulnIDs...)
	}
	added, err := x.client.PackageRecord.CreateBulk(pkgBuilder...).Save(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return added, nil
}

func (x *Client) PutScan(ctx context.Context, scan *ent.Scan, repo *ent.Repository, packages []*ent.PackageRecord) (*ent.Scan, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	scan, err := x.client.Scan.Create().
		SetID(uuid.NewString()).
		SetCommitID(scan.CommitID).
		SetBranch(scan.Branch).
		SetRequestedAt(scan.RequestedAt).
		SetCheckID(scan.CheckID).
		SetPullRequestTarget(scan.PullRequestTarget).
		AddRepository(repo).
		AddPackages(packages...).
		Save(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return scan, nil
}

func (x *Client) GetScan(ctx context.Context, id string) (*ent.Scan, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	got, err := x.client.Scan.Query().Where(scan.ID(id)).
		WithRepository().
		WithPackages(func(prq *ent.PackageRecordQuery) {
			prq.WithStatus().WithVulnerabilities()
		}).
		Only(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return got, nil
}

func (x *Client) GetLatestScan(ctx context.Context, owner, repoName, branch string) (*ent.Scan, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	return nil, nil
}
