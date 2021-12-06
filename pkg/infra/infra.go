package infra

import (
	"io/ioutil"
	"time"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/github"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/infra/opa"
	"github.com/m-mizutani/octovy/pkg/infra/policy"
	"github.com/m-mizutani/octovy/pkg/infra/trivy"
)

type Config struct {
	DBType   string
	DBConfig string `zlog:"secret"`

	GitHubAppID         int64
	GitHubAppPrivateKey string `zlog:"secret"`
	GitHubAppClientID   string
	GitHubAppSecret     string `zlog:"secret"`

	CheckPolicyData string

	OPA opa.Config

	TrivyPath string
}

type Clients struct {
	DB           db.Interface
	GitHub       github.Interface
	NewGitHubApp githubapp.Factory
	CheckPolicy  policy.Check
	Trivy        trivy.Interface
	OPAClient    opa.Interface
	Utils        *Utils

	config *Config
}

type Utils struct {
	Now      func() time.Time
	ReadFile func(fname string) ([]byte, error)
}

func NewUtils() *Utils {
	return &Utils{
		Now:      time.Now,
		ReadFile: ioutil.ReadFile,
	}
}

func New(cfg *Config) (*Clients, error) {
	githubClient, err := github.New(cfg.GitHubAppClientID, cfg.GitHubAppSecret)
	if err != nil {
		return nil, err
	}

	clients := &Clients{
		DB:           db.New(),
		GitHub:       githubClient,
		NewGitHubApp: githubapp.New(cfg.GitHubAppID, []byte(cfg.GitHubAppPrivateKey)),
		Trivy:        trivy.New(),
		Utils:        NewUtils(),

		config: cfg,
	}

	if clients.config.TrivyPath != "" {
		clients.Trivy.SetPath(clients.config.TrivyPath)
	}

	if err := clients.DB.Open(clients.config.DBType, clients.config.DBConfig); err != nil {
		return nil, goerr.Wrap(err)
	}

	if clients.config.CheckPolicyData != "" {
		check, err := policy.NewCheck(clients.config.CheckPolicyData)
		if err != nil {
			return nil, err
		}

		clients.CheckPolicy = check
	}

	if clients.config.OPA.BaseURL != "" {
		client, err := opa.New(&clients.config.OPA)
		if err != nil {
			return nil, err
		}
		clients.OPAClient = client
	}

	return clients, nil
}

func (x *Clients) GitHubAppClientID() string {
	return x.config.GitHubAppClientID
}

func (x *Clients) GitHubAppID() string {
	return x.config.GitHubAppClientID
}
