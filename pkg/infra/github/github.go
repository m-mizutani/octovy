package github

import (
	"context"
	"io"
	"net/http"

	"github.com/google/go-github/v39/github"
	"github.com/m-mizutani/goerr"
)

// This package is used to download trivy database, not used by GitHub App.
type Interface interface {
	ListReleases(ctx context.Context, owner string, repo string) ([]*github.RepositoryRelease, error)
	DownloadReleaseAsset(ctx context.Context, owner string, repo string, assetID int64) (io.ReadCloser, error)
}

type Client struct {
	client *github.Client
}

func New() Interface {
	return &Client{
		client: github.NewClient(&http.Client{}),
	}
}

func (x *Client) ListReleases(ctx context.Context, owner string, repo string) ([]*github.RepositoryRelease, error) {
	opt := &github.ListOptions{}

	releases, resp, err := x.client.Repositories.ListReleases(ctx, owner, repo, opt)
	if err != nil {
		return nil, goerr.Wrap(err, "ListRelease error")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, goerr.New("Returned not 200").With("code", resp.StatusCode)
	}

	return releases, nil
}

func (x *Client) DownloadReleaseAsset(ctx context.Context, owner string, repo string, assetID int64) (io.ReadCloser, error) {
	rc, url, err := x.client.Repositories.DownloadReleaseAsset(ctx, owner, repo, assetID, &http.Client{})
	if err != nil {
		return nil, goerr.Wrap(err).With("url", url)
	}

	return rc, nil
}
