package github

import (
	"github.com/google/go-github/v39/github"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func NewMock() *Mock {
	return &Mock{}
}

type Mock struct {
	AuthenticateMock func(ctx *model.Context, clientID string, clientSecret string, code string) (*model.GitHubToken, error)
	GetUserMock      func(ctx *model.Context, token *model.GitHubToken) (*github.User, error)
}

func (x *Mock) Authenticate(ctx *model.Context, clientID string, clientSecret string, code string) (*model.GitHubToken, error) {
	return x.AuthenticateMock(ctx, clientID, clientSecret, code)
}

func (x *Mock) GetUser(ctx *model.Context, token *model.GitHubToken) (*github.User, error) {
	return x.GetUserMock(ctx, token)
}
