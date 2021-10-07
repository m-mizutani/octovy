package usecase

import (
	"context"

	"github.com/google/go-github/v39/github"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/utils"
)

var logger = utils.Logger

type Interface interface {
	Init() error
	Shutdown()

	// Scan
	SendScanRequest(req *model.ScanRepositoryRequest) error

	// Invoke thread
	InvokeScanThread()

	// DB access proxy
	RegisterRepository(ctx context.Context, repo *ent.Repository) (*ent.Repository, error)
	UpdateVulnStatus(ctx context.Context, req *model.UpdateVulnStatusRequest) error
	LookupScanReport(ctx context.Context, scanID string) (*ent.Scan, error)
	GetRepositories(ctx context.Context) ([]*ent.Repository, error)
	GetVulnerabilities(ctx context.Context, offset, limit int64) ([]*ent.Vulnerability, error)
	GetVulnerabilityCount(ctx context.Context) (int, error)

	// Handle GitHub App Webhook event
	HandleGitHubPushEvent(ctx context.Context, event *github.PushEvent) error
	HandleGitHubPullReqEvent(ctx context.Context, event *github.PullRequestEvent) error
	HandleGitHubInstallationEvent(ctx context.Context, event *github.InstallationEvent) error
	VerifyGitHubSecret(sigSHA256 string, body []byte) error

	// Auth
	CreateAuthState(ctx context.Context) (string, error)
	AuthGitHubUser(ctx context.Context, code, state string) (*ent.User, error)
	LookupUser(ctx context.Context, userID int) (*ent.User, error)
	CreateSession(ctx context.Context, user *ent.User) (*ent.Session, error)
	ValidateSession(ctx context.Context, ssnID string) (*ent.Session, error)
	RevokeSession(ctx context.Context, token string) error

	// Error handling
	HandleError(err error)

	// Config proxy
	GetGitHubAppClientID() string
	FrontendURL() string
	WebhookOnly() bool
}

func New(cfg *model.Config) Interface {
	uc := &usecase{
		config:    cfg,
		infra:     infra.New(),
		scanQueue: make(chan *model.ScanRepositoryRequest, 1024),
	}

	return uc
}

type usecase struct {
	initialized bool
	scanQueue   chan *model.ScanRepositoryRequest

	config *model.Config
	infra  *infra.Interfaces

	// Control usecase for test
	testErrorHandler    func(error)
	disableInvokeThread bool
}

func (x *usecase) Init() error {
	if err := x.initErrorHandler(); err != nil {
		return err
	}

	if err := x.infra.DB.Open(x.config.DBType, x.config.DBConfig); err != nil {
		x.HandleError(err)
		return goerr.Wrap(err)
	}

	if !x.disableInvokeThread {
		x.InvokeScanThread()
	}

	x.initialized = true
	return nil
}

func (x *usecase) Shutdown() {
	x.flushError()
}

func (x *usecase) FrontendURL() string {
	return x.config.FrontendURL
}
func (x *usecase) GetGitHubAppClientID() string {
	return x.config.GitHubAppClientID
}
func (x *usecase) WebhookOnly() bool {
	return x.config.WebhookOnly
}
