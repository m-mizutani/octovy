package githubapp

import "github.com/m-mizutani/octovy/pkg/domain/types"

type Client interface {
	GetContents() error
}

func New(id types.GitHubAppID, key types.GitHubAppPrivateKey) Client {
	return &clientImpl{}
}

type clientImpl struct{}

func (x *clientImpl) GetContents() error {
	return nil
}
