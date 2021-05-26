package githubapp

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

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

func (x *GitHubApp) CreateIssueComment(repo *model.GitHubRepo, prID int, body string) error {
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

func (x *GitHubApp) CreateCheckRun(repo *model.GitHubRepo, commit string) (int64, error) {
	client, err := x.githubClient()
	if err != nil {
		return 0, err
	}

	ctx := context.Background()
	opt := github.CreateCheckRunOptions{
		Name:    "octovy: package vulnerability check",
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

func (x *GitHubApp) UpdateCheckStatus(repo *model.GitHubRepo, checkID int64, status string) error {
	client, err := x.githubClient()
	if err != nil {
		return err
	}
	opt := github.UpdateCheckRunOptions{
		Status: &status,
	}

	ctx := context.Background()

	_, resp, err := client.Checks.UpdateCheckRun(ctx, repo.Owner, repo.RepoName, checkID, opt)
	if err != nil {
		return goerr.Wrap(err, "Failed to update check status").With("repo", repo).With("id", checkID).With("status", status)
	}
	if resp.StatusCode != http.StatusOK {
		return goerr.Wrap(err, "Failed to update status")
	}
	logger.With("repo", repo).With("id", checkID).Info("Created check run")

	return nil
}

func (x *GitHubApp) PutCheckResult(repo *model.GitHubRepo, checkID int64, conclusion string, completedAt time.Time, url string) error {
	client, err := x.githubClient()
	if err != nil {
		return err
	}
	opt := github.UpdateCheckRunOptions{
		Status:      github.String("completed"),
		CompletedAt: &github.Timestamp{Time: completedAt},
		Conclusion:  &conclusion,
		DetailsURL:  &url,
	}

	ctx := context.Background()

	_, resp, err := client.Checks.UpdateCheckRun(ctx, repo.Owner, repo.RepoName, checkID, opt)
	if err != nil {
		return goerr.Wrap(err, "Failed to update check status to complete").With("repo", repo).With("id", checkID).With("opt", opt)
	}
	if resp.StatusCode != http.StatusOK {
		return goerr.Wrap(err, "Failed to update status to complete")
	}
	logger.With("repo", repo).With("id", checkID).Info("Created check run")

	return nil
}
