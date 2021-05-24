package githubapp

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

var logger = golambda.Logger

type GitHubApp struct {
	appID     int64
	installID int64
	pem       []byte
	endpoint  string

	client *github.Client
}

func New(appID, installID int64, pem []byte, endpoint string) interfaces.GitHubApp {
	return &GitHubApp{
		appID:     appID,
		installID: installID,
		pem:       pem,
		endpoint:  endpoint,
	}
}

func (x *GitHubApp) githubClient() (*github.Client, error) {
	if x.client != nil {
		return x.client, nil
	}

	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, x.appID, x.installID, x.pem)

	if err != nil {
		return nil, goerr.Wrap(err)
	}

	endpoint := strings.TrimLeft(x.endpoint, "/")

	if endpoint == "" {
		x.client = github.NewClient(&http.Client{Transport: itr})
	} else {
		itr.BaseURL = endpoint
		httpClient := &http.Client{Transport: itr}
		x.client, err = github.NewEnterpriseClient(endpoint, endpoint, httpClient)
		if err != nil {
			return nil, goerr.Wrap(err).With("endpoint", endpoint)
		}
	}

	return x.client, nil
}

func (x *GitHubApp) GetCodeZip(repo *model.GitHubRepo, commitID string, w io.WriteCloser) error {
	client, err := x.githubClient()
	if err != nil {
		return err
	}

	opt := &github.RepositoryContentGetOptions{
		Ref: commitID,
	}
	ctx := context.Background()

	logger.
		With("appID", x.appID).
		With("repo", repo).
		With("installID", x.installID).
		With("endpoint", x.endpoint).
		With("privateKey.length", len(x.pem)).
		Debug("Sending GetArchiveLink request")

	// https://docs.github.com/en/rest/reference/repos#downloads
	url, r, err := client.Repositories.GetArchiveLink(ctx, repo.Owner, repo.RepoName, github.Zipball, opt, false)
	if err != nil {
		return goerr.Wrap(err)
	}

	logger.With("code", r.StatusCode).Debug("")

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return goerr.Wrap(err)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return goerr.Wrap(err)
	}
	if resp.StatusCode != http.StatusOK {
		var msg string
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			msg = err.Error()
		} else {
			msg = string(body)
		}
		return goerr.New("GitHub download request is failed").With("status", resp.StatusCode).With("body", msg)
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}
