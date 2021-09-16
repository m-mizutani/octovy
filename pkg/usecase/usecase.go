package usecase

import (
	"github.com/google/go-github/v39/github"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/utils"
)

var logger = utils.Logger

type Interface interface {
	// Scan
	SendScanRequest(req *model.ScanRepositoryRequest) error
	RunScanThread() error

	//
	RegisterRepository(repo *ent.Repository) (*ent.Repository, error)

	UpdateVulnStatus(req *model.UpdateVulnStatusRequest) error
	LookupScanReport(reportID string) (*ent.Scan, error)

	HandleGitHubPushEvent(event *github.PushEvent) error
	HandleGitHubPullReqEvent(event *github.PullRequestEvent) error
	HandleGitHubInstallationEvent(event *github.InstallationEvent) error

	// Auth
	GetGitHubAppClientID() (string, error)
	CreateAuthState() (string, error)
	AuthGitHubUser(code, state string) (*ent.User, error)
	LookupUser(userID string) (*ent.User, error)
	CreateSession(user *ent.User) (*ent.Session, error)
	ValidateSession(token string) (*ent.Session, error)
	RevokeSession(token string) error

	// Config
	FrontendURL() string
}

func New(cfg *model.Config) Interface {
	uc := &usecase{
		config: cfg,
		infra:  infra.New(),
	}

	return uc
}

type usecase struct {
	config *model.Config
	infra  infra.Interfaces
}

func (x *usecase) RegisterRepository(repo *ent.Repository) (*ent.Repository, error) {
	panic("not implemented") // TODO: Implement
}

func (x *usecase) UpdateRepositoryDefaultBranch(repo *model.GitHubRepo, branch string) error {
	panic("not implemented") // TODO: Implement
}

func (x *usecase) UpdateVulnStatus(req *model.UpdateVulnStatusRequest) error {
	panic("not implemented") // TODO: Implement
}

func (x *usecase) LookupScanReport(reportID string) (*ent.Scan, error) {
	panic("not implemented") // TODO: Implement
}

// Auth
func (x *usecase) GetGitHubAppClientID() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (x *usecase) CreateAuthState() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (x *usecase) AuthGitHubUser(code string, state string) (*ent.User, error) {
	panic("not implemented") // TODO: Implement
}

func (x *usecase) LookupUser(userID string) (*ent.User, error) {
	panic("not implemented") // TODO: Implement
}

func (x *usecase) CreateSession(user *ent.User) (*ent.Session, error) {
	panic("not implemented") // TODO: Implement
}

func (x *usecase) ValidateSession(token string) (*ent.Session, error) {
	panic("not implemented") // TODO: Implement
}

func (x *usecase) RevokeSession(token string) error {
	panic("not implemented") // TODO: Implement
}

func (x *usecase) FrontendURL() string {
	return x.config.FrontendURL
}
