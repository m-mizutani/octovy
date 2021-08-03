package interfaces

import (
	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

type Usecases interface {
	ScanRepository(req *model.ScanRepositoryRequest) error
	SendScanRequest(req *model.ScanRepositoryRequest) error
	RecvScanRequest() *model.ScanRepositoryRequest

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

	UpdateVulnStatus(response *model.VulnStatus) error

	LookupScanReport(reportID string) (*model.ScanReportResponse, error)

	UpdateTrivyDB() error

	HandleGitHubPushEvent(event *github.PushEvent) error
	HandleGitHubPullReqEvent(event *github.PullRequestEvent) error
	HandleGitHubInstallationEvent(event *github.InstallationEvent) error

	GetGitHubAppClientID() (string, error)
	CreateAuthState() (string, error)
	AuthGitHubUser(code, state string) (*model.User, error)
	LookupUser(userID string) (*model.User, error)
	CreateSession(user *model.User) (*model.Session, error)
	ValidateSession(token string) (*model.Session, error)
	RevokeSession(token string) error

	GetOctovyMetadata() *model.Metadata
}
