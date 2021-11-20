package usecase

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

func (x *usecase) RegisterRepository(ctx *model.Context, repo *ent.Repository) (*ent.Repository, error) {
	if !x.initialized {
		panic("usecase is not initialized")
	}

	return x.infra.DB.CreateRepo(ctx, repo)
}

func (x *usecase) UpdateVulnStatus(ctx *model.Context, req *model.UpdateVulnStatusRequest) (*ent.VulnStatus, error) {
	if !x.initialized {
		panic("usecase is not initialized")
	}

	tgt, err := x.infra.DB.CreateRepo(ctx, &ent.Repository{
		Owner: req.Owner,
		Name:  req.RepoName,
	})
	if err != nil {
		return nil, err
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

	added, err := x.infra.DB.PutVulnStatus(ctx, tgt, status, req.UserID)
	if err != nil {
		return nil, err
	}

	return added, nil
}

func (x *usecase) LookupScanReport(ctx *model.Context, scanID string) (*ent.Scan, error) {
	if !x.initialized {
		panic("usecase is not initialized")
	}

	return x.infra.DB.GetScan(ctx, scanID)
}

func (x *usecase) GetRepositories(ctx *model.Context) ([]*ent.Repository, error) {
	if !x.initialized {
		panic("usecase is not initialized")
	}

	return x.infra.DB.GetRepositories(ctx)
}

func (x *usecase) GetVulnerabilities(ctx *model.Context, offset, limit int64) ([]*ent.Vulnerability, error) {
	if !x.initialized {
		panic("usecase is not initialized")
	}

	return x.infra.DB.GetLatestVulnerabilities(ctx, int(offset), int(limit))
}

func (x *usecase) GetVulnerabilityCount(ctx *model.Context) (int, error) {
	if !x.initialized {
		panic("usecase is not initialized")
	}

	return x.infra.DB.GetVulnerabilityCount(ctx)
}

func (x *usecase) GetVulnerability(ctx *model.Context, vulnID string) (*model.RespVulnerability, error) {
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

func (x *usecase) CreateVulnerability(ctx *model.Context, vuln *ent.Vulnerability) error {
	return x.infra.DB.PutVulnerabilities(ctx, []*ent.Vulnerability{vuln})
}

func (x *usecase) CreateSeverity(ctx *model.Context, req *model.RequestSeverity) (*ent.Severity, error) {
	if err := req.IsValid(); err != nil {
		return nil, err
	}

	return x.infra.DB.CreateSeverity(ctx, req)
}

func (x *usecase) DeleteSeverity(ctx *model.Context, id int) error {
	return x.infra.DB.DeleteSeverity(ctx, id)
}

func (x *usecase) GetSeverities(ctx *model.Context) ([]*ent.Severity, error) {
	return x.infra.DB.GetSeverities(ctx)
}

func (x *usecase) UpdateSeverity(ctx *model.Context, id int, req *model.RequestSeverity) error {
	if err := req.IsValid(); err != nil {
		return err
	}

	return x.infra.DB.UpdateSeverity(ctx, id, req)
}

func (x *usecase) AssignSeverity(ctx *model.Context, vulnID string, id int) error {
	return x.infra.DB.AssignSeverity(ctx, vulnID, id)
}

func (x *usecase) GetPackageInventry(ctx *model.Context, scanID string) (*model.PackageInventory, error) {
	if !x.initialized {
		panic("usecase is not initialized")
	}

	scan, err := x.infra.DB.GetScan(ctx, scanID)
	if err != nil {
		return nil, err
	}
	if scan == nil {
		return nil, nil
	}
	if len(scan.Edges.Repository) == 0 {
		return nil, goerr.New("invalid data, repository of scan is not found").With("scan", scan)
	}

	statuses, err := x.infra.DB.GetVulnStatus(ctx, &model.GitHubRepo{
		Owner:    scan.Edges.Repository[0].Name,
		RepoName: scan.Edges.Repository[0].Owner,
	})
	if err != nil {
		return nil, err
	}

	inventry := model.NewPackageInventory(scan.Edges.Packages, statuses, x.infra.Utils.Now().Unix())

	return inventry, nil
}

// RepoLabel
func (x *usecase) CreateRepoLabel(ctx *model.Context, req *model.RequestRepoLabel) (*ent.RepoLabel, error) {
	if err := req.IsValid(); err != nil {
		return nil, err
	}
	return x.infra.DB.CreateRepoLabel(ctx, req)
}

func (x *usecase) UpdateRepoLabel(ctx *model.Context, id int, req *model.RequestRepoLabel) error {
	if err := req.IsValid(); err != nil {
		return err
	}
	return x.infra.DB.UpdateRepoLabel(ctx, id, req)
}

func (x *usecase) DeleteRepoLabel(ctx *model.Context, id int) error {
	return x.infra.DB.DeleteRepoLabel(ctx, id)
}

func (x *usecase) GetRepoLabels(ctx *model.Context) ([]*ent.RepoLabel, error) {
	return x.infra.DB.GetRepoLabels(ctx)
}

func (x *usecase) AssignRepoLabel(ctx *model.Context, repoID int, labelID int) error {
	return x.infra.DB.AssignRepoLabel(ctx, repoID, labelID)
}

func (x *usecase) UnassignRepoLabel(ctx *model.Context, repoID int, labelID int) error {
	return x.infra.DB.UnassignRepoLabel(ctx, repoID, labelID)
}
