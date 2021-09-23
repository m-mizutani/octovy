package github

import (
	"context"
	"io"

	"github.com/google/go-github/v39/github"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func NewMock() *Mock {
	return &Mock{}
}

type Mock struct {
	ListReleasesMock         func(ctx context.Context, owner string, repo string) ([]*github.RepositoryRelease, error)
	DownloadReleaseAssetMock func(ctx context.Context, owner string, repo string, assetID int64) (io.ReadCloser, error)

	AuthenticateMock func(ctx context.Context, clientID string, clientSecret string, code string) (*model.GitHubToken, error)
	GetUserMock      func(ctx context.Context, token *model.GitHubToken) (*github.User, error)
}

func (x *Mock) ListReleases(ctx context.Context, owner string, repo string) ([]*github.RepositoryRelease, error) {
	return x.ListReleasesMock(ctx, owner, repo)
}

func (x *Mock) DownloadReleaseAsset(ctx context.Context, owner string, repo string, assetID int64) (io.ReadCloser, error) {
	return x.DownloadReleaseAssetMock(ctx, owner, repo, assetID)
}

func (x *Mock) Authenticate(ctx context.Context, clientID string, clientSecret string, code string) (*model.GitHubToken, error) {
	return x.AuthenticateMock(ctx, clientID, clientSecret, code)
}

func (x *Mock) GetUser(ctx context.Context, token *model.GitHubToken) (*github.User, error) {
	return x.GetUserMock(ctx, token)
}
