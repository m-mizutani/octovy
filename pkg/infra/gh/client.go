package gh

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v53/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/utils"
)

type Client struct {
	appID types.GitHubAppID
	pem   types.GitHubAppPrivateKey
}

var _ interfaces.GitHub = (*Client)(nil)

func New(appID types.GitHubAppID, pem types.GitHubAppPrivateKey) (*Client, error) {
	if appID == 0 {
		return nil, goerr.Wrap(types.ErrInvalidOption, "appID is empty")
	}
	if pem == "" {
		return nil, goerr.Wrap(types.ErrInvalidOption, "pem is empty")
	}

	return &Client{
		appID: appID,
		pem:   pem,
	}, nil
}

func (x *Client) buildGithubClient(installID types.GitHubAppInstallID) (*github.Client, error) {
	httpClient, err := x.buildGithubHTTPClient(installID)
	if err != nil {
		return nil, err
	}
	return github.NewClient(httpClient), nil
}

func (x *Client) buildGithubHTTPClient(installID types.GitHubAppInstallID) (*http.Client, error) {
	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, int64(x.appID), int64(installID), []byte(x.pem))

	if err != nil {
		return nil, goerr.Wrap(err, "Failed to create github client")
	}

	client := &http.Client{Transport: itr}
	return client, nil
}

func (x *Client) GetArchiveURL(ctx context.Context, input *interfaces.GetArchiveURLInput) (*url.URL, error) {
	utils.CtxLogger(ctx).Info("Sending GetArchiveLink request",
		slog.Any("appID", x.appID),
		slog.Any("privateKey", x.pem),
		slog.Any("input", input),
	)

	client, err := x.buildGithubClient(input.InstallID)
	if err != nil {
		return nil, err
	}

	opt := &github.RepositoryContentGetOptions{
		Ref: input.CommitID,
	}

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

	utils.CtxLogger(ctx).Debug("GetArchiveLink response", slog.Any("url", url), slog.Any("r", r))

	return url, nil
}

func (x *Client) CreateIssue(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, req *github.IssueRequest) (*github.Issue, error) {
	client, err := x.buildGithubClient(id)
	if err != nil {
		return nil, err
	}

	issue, resp, err := client.Issues.Create(ctx, repo.Owner, repo.RepoName, req)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to create github comment").With("repo", repo).With("req", req)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, goerr.Wrap(err, "failed to create issue").With("repo", repo).With("req", req).With("resp", resp)
	}

	return issue, nil
}

func (x *Client) CreateIssueComment(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int, body string) error {
	client, err := x.buildGithubClient(id)
	if err != nil {
		return err
	}

	comment := &github.IssueComment{Body: &body}

	ret, resp, err := client.Issues.CreateComment(ctx, repo.Owner, repo.RepoName, prID, comment)
	if err != nil {
		return goerr.Wrap(err, "Failed to create github comment").With("repo", repo).With("prID", prID).With("comment", comment)
	}
	if resp.StatusCode != http.StatusCreated {
		return goerr.Wrap(err, "Failed to ")
	}
	utils.Logger().Debug("Commented to PR", "comment", ret)

	return nil
}

//go:embed queries/list_comments.graphql
var queryListIssueComments string

func (x *Client) ListIssueComments(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int) ([]*model.GitHubIssueComment, error) {
	type response struct {
		Repository struct {
			PullRequest struct {
				Comments struct {
					Edges []struct {
						Cursor string `json:"cursor"`
						Node   struct {
							ID     string `json:"id"`
							Author struct {
								Login string `json:"login"`
							} `json:"author"`
							Body        string `json:"body"`
							IsMinimized bool   `json:"isMinimized"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"comments"`
				Title string `json:"title"`
			} `json:"pullRequest"`
		} `json:"repository"`
	}

	var comments []*model.GitHubIssueComment

	var cursor *string
	for {
		vars := map[string]any{
			"owner":       repo.Owner,
			"name":        repo.RepoName,
			"issueNumber": prID,
		}
		if cursor != nil {
			vars["cursor"] = *cursor
		}
		resp, err := x.queryGraphQL(ctx, id, &gqlRequest{
			Query:     queryListIssueComments,
			Variables: vars,
		})

		if err != nil {
			return nil, err
		}

		var data response
		if err := json.Unmarshal(resp.Data, &data); err != nil {
			return nil, goerr.Wrap(err, "Failed to unmarshal response")
		}

		if len(data.Repository.PullRequest.Comments.Edges) == 0 {
			break
		}

		for _, edge := range data.Repository.PullRequest.Comments.Edges {
			comments = append(comments, &model.GitHubIssueComment{
				ID:          edge.Node.ID,
				Login:       edge.Node.Author.Login,
				Body:        edge.Node.Body,
				IsMinimized: edge.Node.IsMinimized,
			})
			cursor = &edge.Cursor
		}
	}

	return comments, nil
}

//go:embed queries/minimize_comment.graphql
var queryMinimizeComment string

func (x *Client) MinimizeComment(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, subjectID string) error {
	req := &gqlRequest{
		Query: queryMinimizeComment,
		Variables: map[string]any{
			"id": subjectID,
		},
	}

	_, err := x.queryGraphQL(ctx, id, req)
	if err != nil {
		return err
	}

	return nil
}

type gqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}
type gqlResponse struct {
	Data  json.RawMessage `json:"data"`
	Error json.RawMessage `json:"errors"`
}

func (x *Client) queryGraphQL(ctx context.Context, id types.GitHubAppInstallID, req *gqlRequest) (*gqlResponse, error) {
	client, err := x.buildGithubHTTPClient(id)
	if err != nil {
		return nil, err
	}

	rawReq, err := json.Marshal(req)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to marshal graphQL request").With("req", req)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiGraphQLEndpoint, bytes.NewReader(rawReq))
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to create graphQL request").With("req", req)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to send graphQL request").With("req", req)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, goerr.Wrap(err, "Failed to get graphQL response").With("req", httpReq).With("resp", resp).With("body", string(body))
	}

	var gqlResp gqlResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to read response body").With("resp", resp)
	}
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return nil, goerr.Wrap(err, "Failed to decode response").With("resp", resp)
	}

	return &gqlResp, nil
}

const (
	apiGraphQLEndpoint = "https://api.github.com/graphql"
)

/*
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

func (x *Client) UpdateCheckRun(repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error {
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
