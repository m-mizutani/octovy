package githubapp

import (
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v53/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
)

type Client interface {
	GetArchiveURL(ctx *model.Context, input *GetArchiveURLInput) (*url.URL, error)
	// CreateIssueComment(repo *model.GitHubRepo, prID int, body string) error
	// CreateCheckRun(repo *model.GitHubRepo, commit string) (int64, error)
	// UpdateCheckRun(repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error
}

type GetArchiveURLInput struct {
	Owner     string
	Repo      string
	CommitID  string
	InstallID types.GitHubAppInstallID
}

type clientImpl struct {
	appID types.GitHubAppID
	pem   types.GitHubAppPrivateKey
}

func New(appID types.GitHubAppID, pem types.GitHubAppPrivateKey) (Client, error) {
	if appID == 0 {
		return nil, goerr.Wrap(types.ErrInvalidOption, "appID is empty")
	}
	if pem == "" {
		return nil, goerr.Wrap(types.ErrInvalidOption, "pem is empty")
	}

	return &clientImpl{
		appID: appID,
		pem:   pem,
	}, nil
}

func (x *clientImpl) buildGithubClient(installID types.GitHubAppInstallID) (*github.Client, error) {
	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, int64(x.appID), int64(installID), []byte(x.pem))

	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return github.NewClient(&http.Client{Transport: itr}), nil
}

func (x *clientImpl) GetArchiveURL(ctx *model.Context, input *GetArchiveURLInput) (*url.URL, error) {
	client, err := x.buildGithubClient(input.InstallID)
	if err != nil {
		return nil, err
	}

	opt := &github.RepositoryContentGetOptions{
		Ref: input.CommitID,
	}

	ctx.Logger().Debug("Sending GetArchiveLink request",
		slog.Any("appID", x.appID),
		slog.Any("repo", input.Repo),
		slog.Any("installID", input.InstallID),
		slog.Any("privateKey.length", len(x.pem)),
	)

	// https://docs.github.com/en/rest/reference/repos#downloads
	// https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28#get-archive-link
	url, r, err := client.Repositories.GetArchiveLink(ctx, input.Owner, input.Repo, github.Zipball, opt, false)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	if r.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(r.Body)
		return nil, goerr.Wrap(err, "Failed to get archive link").With("status", r.StatusCode).With("body", string(body))
	}

	ctx.Logger().Debug("GetArchiveLink response", slog.Any("url", url), slog.Any("r", r))

	return url, nil
}

/*
func (x *clientImpl) CreateIssueComment(repo *model.GitHubRepo, prID int, body string) error {
	client, err := x.githubClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	comment := &github.IssueComment{Body: &body}

	ret, resp, err := client.Issues.CreateComment(ctx, repo.Owner, repo.Name, prID, comment)
	if err != nil {
		return goerr.Wrap(err, "Failed to create github comment").With("repo", repo).With("prID", prID).With("comment", comment)
	}
	if resp.StatusCode != http.StatusCreated {
		return goerr.Wrap(err, "Failed to ")
	}
	utils.Logger.With("comment", ret).Debug("Commented to PR")

	return nil
}

func (x *clientImpl) CreateCheckRun(repo *model.GitHubRepo, commit string) (int64, error) {
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

	run, resp, err := client.Checks.CreateCheckRun(ctx, repo.Owner, repo.Name, opt)
	if err != nil {
		return 0, goerr.Wrap(err, "Failed to create check run").With("repo", repo).With("commit", commit)
	}
	if resp.StatusCode != http.StatusCreated {
		return 0, goerr.Wrap(err, "Failed to ")
	}
	utils.Logger.With("run", run).Debug("Created check run")

	return *run.ID, nil
}

func (x *clientImpl) UpdateCheckRun(repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error {
	client, err := x.githubClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	_, resp, err := client.Checks.UpdateCheckRun(ctx, repo.Owner, repo.Name, checkID, *opt)
	if err != nil {
		return goerr.Wrap(err, "Failed to update check status to complete").With("repo", repo).With("id", checkID).With("opt", opt)
	}
	if resp.StatusCode != http.StatusOK {
		return goerr.Wrap(err, "Failed to update status to complete")
	}
	utils.Logger.With("repo", repo).With("id", checkID).Debug("Created check run")

	return nil
}
*/
