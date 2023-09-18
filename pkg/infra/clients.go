package infra

import (
	"database/sql"
	"net/http"

	gh "github.com/m-mizutani/octovy/pkg/infra/gh"
	"github.com/m-mizutani/octovy/pkg/infra/trivy"
)

type Clients struct {
	githubApp   gh.Client
	httpClient  HTTPClient
	trivyClient trivy.Client
	dbClient    *sql.DB
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Option func(*Clients)

func New(options ...Option) *Clients {
	client := &Clients{
		httpClient:  http.DefaultClient,
		trivyClient: trivy.New("trivy"),
	}

	for _, opt := range options {
		opt(client)
	}

	return client
}

func (x *Clients) GitHubApp() gh.Client {
	return x.githubApp
}
func (x *Clients) HTTPClient() HTTPClient {
	return x.httpClient
}
func (x *Clients) Trivy() trivy.Client {
	return x.trivyClient
}
func (x *Clients) DB() *sql.DB {
	return x.dbClient
}

func WithGitHubApp(client gh.Client) Option {
	return func(x *Clients) {
		x.githubApp = client
	}
}

func WithHTTPClient(client HTTPClient) Option {
	return func(x *Clients) {
		x.httpClient = client
	}
}

func WithTrivy(client trivy.Client) Option {
	return func(x *Clients) {
		x.trivyClient = client
	}
}

func WithDB(client *sql.DB) Option {
	return func(x *Clients) {
		x.dbClient = client
	}
}
