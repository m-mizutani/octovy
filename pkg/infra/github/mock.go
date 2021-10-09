package github

import (
	"context"

	"github.com/google/go-github/v39/github"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func NewMock() *Mock {
	return &Mock{}
}

type Mock struct {
	AuthenticateMock func(ctx context.Context, clientID string, clientSecret string, code string) (*model.GitHubToken, error)
	GetUserMock      func(ctx context.Context, token *model.GitHubToken) (*github.User, error)
}

func (x *Mock) Authenticate(ctx context.Context, clientID string, clientSecret string, code string) (*model.GitHubToken, error) {
	return x.AuthenticateMock(ctx, clientID, clientSecret, code)
}

func (x *Mock) GetUser(ctx context.Context, token *model.GitHubToken) (*github.User, error) {
	return x.GetUserMock(ctx, token)
}
