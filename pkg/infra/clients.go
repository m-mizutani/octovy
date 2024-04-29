package infra

import (
	"net/http"

	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/infra/gh"
	"github.com/m-mizutani/octovy/pkg/infra/trivy"
)

type Clients struct {
	githubApp   gh.Client
	httpClient  HTTPClient
	trivyClient trivy.Client
	bqClient    interfaces.BigQuery
	fsClient    interfaces.Firestore
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
func (x *Clients) BigQuery() interfaces.BigQuery {
	return x.bqClient
}
func (x *Clients) Firestore() interfaces.Firestore {
	return x.fsClient
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

func WithBigQuery(client interfaces.BigQuery) Option {
	return func(x *Clients) {
		x.bqClient = client
	}
}

func WithFirestore(client interfaces.Firestore) Option {
	return func(x *Clients) {
		x.fsClient = client
	}
}
