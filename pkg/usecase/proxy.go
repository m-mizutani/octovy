package usecase

import (
	"context"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

func (x *usecase) RegisterRepository(ctx context.Context, repo *ent.Repository) (*ent.Repository, error) {
	if !x.initialized {
		panic("usecase is not initialized")
	}

	return x.infra.DB.CreateRepo(ctx, repo)
}

func (x *usecase) UpdateVulnStatus(ctx context.Context, req *model.UpdateVulnStatusRequest) error {
	if !x.initialized {
		panic("usecase is not initialized")
	}

	tgt, err := x.infra.DB.CreateRepo(ctx, &ent.Repository{
		Owner: req.Owner,
		Name:  req.RepoName,
	})
	if err != nil {
		return err
	}

	status := &ent.VulnStatus{
		Status:    req.Status,
		Source:    req.Source,
		PkgName:   req.PkgName,
		PkgType:   req.PkgType,
		ExpiresAt: req.ExpiresAt,
		CreatedAt: x.infra.Utils.Now().Unix(),
		VulnID:    req.VulnID,
		Comment:   req.Comment,
	}

	if err := x.infra.DB.PutVulnStatus(ctx, tgt, status, req.UserID); err != nil {
		return err
	}

	return nil
}

func (x *usecase) LookupScanReport(ctx context.Context, scanID string) (*ent.Scan, error) {
	if !x.initialized {
		panic("usecase is not initialized")
	}

	return x.infra.DB.GetScan(ctx, scanID)
}

func (x *usecase) GetRepositories(ctx context.Context) ([]*ent.Repository, error) {
	if !x.initialized {
		panic("usecase is not initialized")
	}

	return x.infra.DB.GetRepositories(ctx)
}

func (x *usecase) GetVulnerabilities(ctx context.Context, offset, limit int64) ([]*ent.Vulnerability, error) {
	if !x.initialized {
		panic("usecase is not initialized")
	}

	return x.infra.DB.GetLatestVulnerabilities(ctx, int(offset), int(limit))
}

func (x *usecase) GetVulnerabilityCount(ctx context.Context) (int, error) {
	if !x.initialized {
		panic("usecase is not initialized")
	}

	return x.infra.DB.GetVulnerabilityCount(ctx)
}

func (x *usecase) GetVulnerability(ctx context.Context, vulnID string) (*model.RespVulnerability, error) {
	vuln, err := x.infra.DB.GetVulnerability(ctx, vulnID)
	if err != nil {
		return nil, err
	}
	if vuln == nil {
		return nil, nil
	}

	repos, err := x.infra.DB.GetRepositoriesWithVuln(ctx, vulnID)
	if err != nil {
		return nil, err
	}

	return &model.RespVulnerability{
		Vulnerability: vuln,
		Affected:      repos,
	}, nil
}
