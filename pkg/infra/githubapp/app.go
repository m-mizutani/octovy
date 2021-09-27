package githubapp

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v39/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/utils"
)

var logger = utils.Logger

type Interface interface {
	GetCodeZip(repo *model.GitHubRepo, commitID string, w io.WriteCloser) error
	CreateIssueComment(repo *model.GitHubRepo, prID int, body string) error
	CreateCheckRun(repo *model.GitHubRepo, commit string) (int64, error)
	UpdateCheckRun(repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error
}

type Client struct {
	appID     int64
	installID int64
	pem       []byte

	client *github.Client
}

type Factory func(appID, installID int64, pem []byte) Interface

func New(appID, installID int64, pem []byte) Interface {
	return &Client{
		appID:     appID,
		installID: installID,
		pem:       pem,
	}
}

func (x *Client) githubClient() (*github.Client, error) {
	if x.client != nil {
		return x.client, nil
	}

	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, x.appID, x.installID, x.pem)

	if err != nil {
		return nil, goerr.Wrap(err)
	}

	x.client = github.NewClient(&http.Client{Transport: itr})

	return x.client, nil
}

func (x *Client) GetCodeZip(repo *model.GitHubRepo, commitID string, w io.WriteCloser) error {
	client, err := x.githubClient()
	if err != nil {
		return err
	}

	opt := &github.RepositoryContentGetOptions{
		Ref: commitID,
	}
	ctx := context.Background()

	logger.Debug().
		Interface("appID", x.appID).
		Interface("repo", repo).
		Interface("installID", x.installID).
		Interface("privateKey.length", len(x.pem)).
		Msg("Sending GetArchiveLink request")

	// https://docs.github.com/en/rest/reference/repos#downloads
	url, r, err := client.Repositories.GetArchiveLink(ctx, repo.Owner, repo.RepoName, github.Zipball, opt, false)
	if err != nil {
		return goerr.Wrap(err)
	}

	logger.Debug().Interface("code", r.StatusCode).Send()

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

func (x *Client) CreateIssueComment(repo *model.GitHubRepo, prID int, body string) error {
	client, err := x.githubClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	comment := &github.IssueComment{Body: &body}

	ret, resp, err := client.Issues.CreateComment(ctx, repo.Owner, repo.RepoName, prID, comment)
	if err != nil {
		return goerr.Wrap(err, "Failed to create github comment").With("repo", repo).With("prID", prID).With("comment", comment)
	}
	if resp.StatusCode != http.StatusCreated {
		return goerr.Wrap(err, "Failed to ")
	}
	logger.With("comment", ret).Info("Commented to PR")

	return nil
}

func (x *Client) CreateCheckRun(repo *model.GitHubRepo, commit string) (int64, error) {
	client, err := x.githubClient()
	if err != nil {
		return 0, err
	}

	ctx := context.Background()
	opt := github.CreateCheckRunOptions{
		Name:    "Octovy: package vulnerability check",
		HeadSHA: commit,
		Status:  github.String("in_progress"),
	}

	run, resp, err := client.Checks.CreateCheckRun(ctx, repo.Owner, repo.RepoName, opt)
	if err != nil {
		return 0, goerr.Wrap(err, "Failed to create check run").With("repo", repo).With("commit", commit)
	}
	if resp.StatusCode != http.StatusCreated {
		return 0, goerr.Wrap(err, "Failed to ")
	}
	logger.With("run", run).Info("Created check run")

	return *run.ID, nil
}

func (x *Client) UpdateCheckRun(repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error {
	client, err := x.githubClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	_, resp, err := client.Checks.UpdateCheckRun(ctx, repo.Owner, repo.RepoName, checkID, *opt)
	if err != nil {
		return goerr.Wrap(err, "Failed to update check status to complete").With("repo", repo).With("id", checkID).With("opt", opt)
	}
	if resp.StatusCode != http.StatusOK {
		return goerr.Wrap(err, "Failed to update status to complete")
	}
	logger.With("repo", repo).With("id", checkID).Info("Created check run")

	return nil
}
