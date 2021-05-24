package github

import (
	"context"
	"io"
	"net/http"

	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
)

// This package is used to download trivy database, not used by GitHub App.

type Client struct {
	client *github.Client
}

func New() interfaces.GitHubClient {
	return &Client{
		client: github.NewClient(&http.Client{}),
	}
}

func (x *Client) ListReleases(owner string, repo string) ([]*github.RepositoryRelease, error) {
	ctx := context.Background()
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

func (x *Client) DownloadReleaseAsset(owner string, repo string, assetID int64) (io.ReadCloser, error) {
	ctx := context.Background()
	rc, url, err := x.client.Repositories.DownloadReleaseAsset(ctx, owner, repo, assetID, &http.Client{})
	if err != nil {
		return nil, goerr.Wrap(err).With("url", url)
	}

	return rc, nil
}
