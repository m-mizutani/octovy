package usecase

import (
	"testing"

	gh "github.com/google/go-github/v39/github"
	"github.com/stretchr/testify/require"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/github"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/infra/trivy"
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
	RegisterRepository(ctx *model.Context, repo *ent.Repository) (*ent.Repository, error)
	UpdateVulnStatus(ctx *model.Context, req *model.UpdateVulnStatusRequest) (*ent.VulnStatus, error)
	LookupScanReport(ctx *model.Context, scanID string) (*ent.Scan, error)
	GetRepositories(ctx *model.Context) ([]*ent.Repository, error)
	GetVulnerabilities(ctx *model.Context, offset, limit int64) ([]*ent.Vulnerability, error)
	GetVulnerabilityCount(ctx *model.Context) (int, error)
	GetVulnerability(ctx *model.Context, vulnID string) (*model.RespVulnerability, error)
	CreateVulnerability(ctx *model.Context, vuln *ent.Vulnerability) error

	// Severity
	CreateSeverity(ctx *model.Context, label string) (*ent.Severity, error)
	DeleteSeverity(ctx *model.Context, id int) error
	GetSeverities(ctx *model.Context) ([]*ent.Severity, error)
	UpdateSeverity(ctx *model.Context, id int, label string) error
	AssignSeverity(ctx *model.Context, vulnID string, id int) error

	// Handle GitHub App Webhook event
	HandleGitHubPushEvent(ctx *model.Context, event *gh.PushEvent) error
	HandleGitHubPullReqEvent(ctx *model.Context, event *gh.PullRequestEvent) error
	HandleGitHubInstallationEvent(ctx *model.Context, event *gh.InstallationEvent) error
	VerifyGitHubSecret(sigSHA256 string, body []byte) error

	// Auth
	CreateAuthState(ctx *model.Context) (string, error)
	AuthGitHubUser(ctx *model.Context, code, state string) (*ent.User, error)
	LookupUser(ctx *model.Context, userID int) (*ent.User, error)
	CreateSession(ctx *model.Context, user *ent.User) (*ent.Session, error)
	ValidateSession(ctx *model.Context, ssnID string) (*ent.Session, error)
	RevokeSession(ctx *model.Context, token string) error

	// Error handling
	HandleError(ctx *model.Context, err error)

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

func NewTest(t *testing.T) Interface {
	uc := New(&model.Config{}).(*usecase)

	dbClient := db.NewMock(t)
	ghClient := github.NewMock()
	newGitHubApp, _ := githubapp.NewMock()
	util := infra.NewUtils()
	trivyClient := trivy.NewMock()

	uc.disableInvokeThread = true
	uc.infra = &infra.Interfaces{
		DB:           dbClient,
		GitHub:       ghClient,
		NewGitHubApp: newGitHubApp,
		Trivy:        trivyClient,
		Utils:        util,
	}

	uc.testErrorHandler = func(err error) {
		require.NoError(t, err)
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

	if x.config.TrivyPath != "" {
		x.infra.Trivy.SetPath(x.config.TrivyPath)
	}

	if err := x.infra.DB.Open(x.config.DBType, x.config.DBConfig); err != nil {
		x.HandleError(model.NewContext(), err)
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
