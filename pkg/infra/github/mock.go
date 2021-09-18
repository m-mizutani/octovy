package github

import (
	"context"
	"io"

	"github.com/google/go-github/v39/github"
)

func NewMock() *Mock {
	return &Mock{}
}

type Mock struct {
	ListReleasesMock         func(ctx context.Context, owner string, repo string) ([]*github.RepositoryRelease, error)
	DownloadReleaseAssetMock func(ctx context.Context, owner string, repo string, assetID int64) (io.ReadCloser, error)
}

func (x *Mock) ListReleases(ctx context.Context, owner string, repo string) ([]*github.RepositoryRelease, error) {
	return x.ListReleasesMock(ctx, owner, repo)
}

func (x *Mock) DownloadReleaseAsset(ctx context.Context, owner string, repo string, assetID int64) (io.ReadCloser, error) {
	return x.DownloadReleaseAssetMock(ctx, owner, repo, assetID)
}
