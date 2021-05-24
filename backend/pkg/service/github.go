package service

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

func (x *Service) GetCodeZip(repo *model.GitHubRepo, commitID string, installID int64, w io.WriteCloser) error {
	defer w.Close()

	secrets, err := x.GetSecrets()
	if err != nil {
		return err
	}

	tr := http.DefaultTransport
	pem, err := secrets.GithubAppPEM()
	if err != nil {
		return err
	}
	appID, err := secrets.GetGitHubAppID()
	if err != nil {
		return err
	}

	itr, err := ghinstallation.New(tr, appID, installID, pem)

	if err != nil {
		return goerr.Wrap(err)
	}

	endpoint := strings.TrimLeft(x.config.GitHubEndpoint, "/")
	githubHTTP := x.Infra.NewHTTP(itr)

	var client *github.Client
	if endpoint == "" {
		client = github.NewClient(githubHTTP)
	} else {
		itr.BaseURL = endpoint
		httpClient := x.Infra.NewHTTP(itr)
		client, err = github.NewEnterpriseClient(endpoint, endpoint, httpClient)
		if err != nil {
			return goerr.Wrap(err).With("endpoint", endpoint)
		}
	}

	opt := &github.RepositoryContentGetOptions{
		Ref: commitID,
	}
	ctx := context.Background()

	logger.
		With("appID", appID).
		With("repo", repo).
		With("installID", installID).
		With("endpoint", endpoint).
		With("privateKey.length", len(secrets.GitHubAppPrivateKey)).
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

	httpClient := x.Infra.NewHTTP(nil)
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
