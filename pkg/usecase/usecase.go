package usecase

import (
	"testing"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/github"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/infra/policy"
	"github.com/m-mizutani/octovy/pkg/infra/trivy"
	"github.com/m-mizutani/octovy/pkg/utils"
)

var logger = utils.Logger

func New(cfg *model.Config) *Usecase {
	uc := &Usecase{
		config:    cfg,
		infra:     infra.New(),
		scanQueue: make(chan *model.ScanRepositoryRequest, 1024),
	}

	return uc
}

type TestOption func(*Usecase)

func OptInjectDB(client *db.Client) TestOption {
	return func(u *Usecase) {
		u.infra.DB = client
	}
}

func OptInjectErrorHandler(f func(error)) TestOption {
	return func(u *Usecase) {
		u.testErrorHandler = f
	}
}

func NewTest(t *testing.T, options ...TestOption) *Usecase {
	uc := New(&model.Config{})

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

	for _, opt := range options {
		opt(uc)
	}

	return uc
}

type Usecase struct {
	initialized bool
	scanQueue   chan *model.ScanRepositoryRequest

	config *model.Config
	infra  *infra.Interfaces

	// Control usecase for test
	testErrorHandler    func(error)
	disableInvokeThread bool
}

func (x *Usecase) Init() error {
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

	if x.config.CheckPolicyData != "" {
		check, err := policy.NewCheck(x.config.CheckPolicyData)
		if err != nil {
			return err
		}

		x.infra.CheckPolicy = check
	}

	x.initialized = true
	return nil
}

func (x *Usecase) Shutdown() {
	x.flushError()
}

func (x *Usecase) FrontendURL() string {
	return x.config.FrontendURL
}
func (x *Usecase) GetGitHubAppClientID() string {
	return x.config.GitHubAppClientID
}
func (x *Usecase) WebhookOnly() bool {
	return x.config.WebhookOnly
}
