package githubapp

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/types"
)

type Client interface {
	GetContents() error
}

func New(id types.GitHubAppID, key types.GitHubAppPrivateKey) (Client, error) {
	if id == "" || key == "" {
		return nil, goerr.Wrap(types.ErrInvalidOption, "GitHub App ID and Private Key are required")
	}

	return &clientImpl{
		id:  id,
		key: key,
	}, nil
}

type clientImpl struct {
	id  types.GitHubAppID
	key types.GitHubAppPrivateKey
}

func (x *clientImpl) GetContents() error {
	return nil
}
