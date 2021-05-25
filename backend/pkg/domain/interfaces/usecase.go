package interfaces

import (
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

type Usecases interface {
	ScanRepository(req *model.ScanRepositoryRequest) error
	SendScanRequest(req *model.ScanRepositoryRequest) error

	FeedbackScanResult(req *model.FeedbackRequest) error

	RegisterRepository(repo *model.Repository) error
	PutNewRepository(repo *model.Repository) (bool, error)
	UpdateRepositoryDefaultBranch(repo *model.GitHubRepo, branch string) error

	FindOwners() ([]*model.Owner, error)
	FindRepos() ([]*model.Repository, error)
	FindReposByOwner(owner string) ([]*model.Repository, error)
	FindReposByFullName(owner, name string) (*model.Repository, error)
	LookupBranch(branch *model.GitHubBranch) (*model.Branch, error)
	FindPkgs(pkgType model.PkgType, name string) ([]*model.PackageRecord, error)
	FindPkgsByRepo(branch *model.GitHubBranch) ([]*model.PackageRecord, error)
	FindVulnerability(vulnID string) (*model.Vulnerability, error)
	FindPackageRecordsByBranch(*model.GitHubBranch) ([]*model.PackageRecord, error)
	FindPackageRecordsByName(pkgType model.PkgType, pkgName string) ([]*model.PackageRecord, error)

	LookupScanReport(reportID string) (*model.ScanReport, error)

	UpdateTrivyDB() error
}
