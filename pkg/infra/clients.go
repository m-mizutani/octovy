package infra

import "github.com/m-mizutani/octovy/pkg/infra/githubapp"

type Clients struct {
	githubApp githubapp.Client
}

type Option func(*Clients)

func New(options ...Option) *Clients {
	client := &Clients{}

	return client
}

func (x *Clients) GitHubApp() githubapp.Client {
	return x.githubApp
}

func WithGitHubApp(client githubapp.Client) Option {
	return func(x *Clients) {
		x.githubApp = client
	}
}
