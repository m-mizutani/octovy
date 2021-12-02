package usecase

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/github"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/infra/trivy"
	"github.com/m-mizutani/octovy/pkg/utils"
)

var logger = utils.Logger

func New(cfg *model.Config, ifs *infra.Clients) *Usecase {
	uc := &Usecase{
		config:    cfg,
		infra:     ifs,
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

	dbClient := db.NewMock(t)
	ghClient := github.NewMock()
	newGitHubApp, _ := githubapp.NewMock()
	util := infra.NewUtils()
	trivyClient := trivy.NewMock()

	uc := New(&model.Config{}, &infra.Clients{
		DB:           dbClient,
		GitHub:       ghClient,
		NewGitHubApp: newGitHubApp,
		Trivy:        trivyClient,
		Utils:        util,
	})
	uc.disableInvokeThread = true

	for _, opt := range options {
		opt(uc)
	}

	return uc
}

type Usecase struct {
	scanQueue chan *model.ScanRepositoryRequest

	config *model.Config
	infra  *infra.Clients

	// Control usecase for test
	testErrorHandler    func(error)
	disableInvokeThread bool
}

func (x *Usecase) Close() {
	x.flushError()
}

func (x *Usecase) FrontendURL() string {
	return x.config.FrontendURL
}
func (x *Usecase) GetGitHubAppClientID() string {
	return x.infra.GitHubAppClientID()
}
func (x *Usecase) WebhookOnly() bool {
	return x.config.WebhookOnly
}
