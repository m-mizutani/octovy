package controller_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/controller"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
	"github.com/m-mizutani/octovy/backend/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockUsecase struct {
	usecase.Default
	sendScanRequest    func(r *model.ScanRepositoryRequest) error
	registerRepository func(repo *model.Repository) error
}

func (x *MockUsecase) SendScanRequest(r *model.ScanRepositoryRequest) error {
	return x.sendScanRequest(r)
}

func (x *MockUsecase) RegisterRepository(repo *model.Repository) error {
	return x.registerRepository(repo)
}

func TestLambdaAPIWebhook(t *testing.T) {
	ts := time.Now().UTC()
	t.Run("Push event", func(t *testing.T) {
		var req *model.ScanRepositoryRequest
		var repo *model.Repository
		ctrl := controller.New()
		ctrl.Usecase = &MockUsecase{
			sendScanRequest: func(r *model.ScanRepositoryRequest) error {
				req = r
				return nil
			},
			registerRepository: func(r *model.Repository) error {
				repo = r
				return nil
			},
		}

		pushEvent := github.PushEvent{
			Ref: github.String("refs/heads/master"),
			Repo: &github.PushEventRepository{
				Name: github.String("blue"),
				Owner: &github.User{
					Name: github.String("five"),
				},
				HTMLURL:       github.String("https://github-enterprise.example.com/blue/five"),
				DefaultBranch: github.String("default"),
			},

			Commits: []github.PushEventCommit{
				{
					ID:        github.String("abcdef123"),
					Timestamp: &github.Timestamp{Time: ts},
				},
				{
					ID:        github.String("beefcafe"),
					Timestamp: &github.Timestamp{Time: ts.Add(time.Minute)},
				},
				{
					ID:        github.String("bbbbbbbb"),
					Timestamp: &github.Timestamp{Time: ts.Add(time.Second)},
				},
			},
			Installation: &github.Installation{
				ID: github.Int64(1234),
			},
		}
		raw, err := json.Marshal(pushEvent)
		require.NoError(t, err)

		event := golambda.Event{
			Origin: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Path:       "/api/v1/webhook/github",
				Headers: map[string]string{
					"X-GitHub-Event": "push",
				},
				Body: string(raw),
			},
		}
		resp, err := ctrl.LambdaAPIHandler(event)
		require.NoError(t, err)
		assert.NotNil(t, resp)
		httpResp, ok := resp.(events.APIGatewayProxyResponse)
		require.True(t, ok)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		require.NotNil(t, req)
		assert.Equal(t, "five", req.Owner)
		assert.Equal(t, "blue", req.RepoName)
		assert.Equal(t, "master", req.Branch)
		assert.Equal(t, int64(1234), req.InstallID)
		assert.Equal(t, "beefcafe", req.CommitID)
		assert.Equal(t, ts.Add(time.Minute).Unix(), req.UpdatedAt)
		assert.False(t, req.IsPullRequest)
		assert.False(t, req.IsTargetBranch)

		require.NotNil(t, repo)
	})

	t.Run("Push event with default branch", func(t *testing.T) {
		var req *model.ScanRepositoryRequest
		ctrl := controller.New()
		ctrl.Usecase = &MockUsecase{
			sendScanRequest: func(r *model.ScanRepositoryRequest) error {
				req = r
				return nil
			},
			registerRepository: func(r *model.Repository) error { return nil },
		}

		pushEvent := github.PushEvent{
			Ref: github.String("refs/heads/master"),
			Repo: &github.PushEventRepository{
				Name: github.String("blue"),
				Owner: &github.User{
					Name: github.String("five"),
				},
				HTMLURL:       github.String("https://github-enterprise.example.com/blue/five"),
				DefaultBranch: github.String("master"),
			},

			Commits: []github.PushEventCommit{
				{
					ID:        github.String("abcdef123"),
					Timestamp: &github.Timestamp{Time: ts},
				},
			},
			Installation: &github.Installation{
				ID: github.Int64(1234),
			},
		}
		raw, err := json.Marshal(pushEvent)
		require.NoError(t, err)

		event := golambda.Event{
			Origin: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Path:       "/api/v1/webhook/github",
				Headers: map[string]string{
					"X-GitHub-Event": "push",
				},
				Body: string(raw),
			},
		}
		_, err = ctrl.LambdaAPIHandler(event)
		require.NoError(t, err)
		require.NotNil(t, req)
		assert.True(t, req.IsTargetBranch)
	})

	t.Run("Pull request event", func(t *testing.T) {
		var req *model.ScanRepositoryRequest
		var repo *model.Repository
		ctrl := controller.New()
		ctrl.Usecase = &MockUsecase{
			sendScanRequest: func(r *model.ScanRepositoryRequest) error {
				req = r
				return nil
			},
			registerRepository: func(r *model.Repository) error {
				repo = r
				return nil
			},
		}

		pullReqEvent := makePullRequestEvent(&ts)
		raw, err := json.Marshal(pullReqEvent)
		require.NoError(t, err)

		event := golambda.Event{
			Origin: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Path:       "/api/v1/webhook/github",
				Headers: map[string]string{
					"X-GitHub-Event": "pull_request",
				},
				Body: string(raw),
			},
		}

		resp, err := ctrl.LambdaAPIHandler(event)
		require.NoError(t, err)
		require.NotNil(t, resp)
		httpResp, ok := resp.(events.APIGatewayProxyResponse)
		require.True(t, ok)

		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		require.NotNil(t, req)
		assert.Equal(t, "five", req.Owner)
		assert.Equal(t, "blue", req.RepoName)
		assert.Equal(t, "ao:1", req.Branch)
		assert.Equal(t, int64(1234), req.InstallID)
		assert.Equal(t, "xyz", req.CommitID)
		assert.Equal(t, ts.Unix(), req.UpdatedAt)
		assert.True(t, req.IsPullRequest)
		assert.False(t, req.IsTargetBranch)

		require.NotNil(t, repo)
		assert.Equal(t, "five", repo.Owner)
		assert.Equal(t, "blue", repo.RepoName)
	})

	t.Run("Pull request event of synchronize", func(t *testing.T) {
		var req *model.ScanRepositoryRequest
		ctrl := controller.New()
		ctrl.Usecase = &MockUsecase{
			sendScanRequest: func(r *model.ScanRepositoryRequest) error {
				req = r
				return nil
			},
			registerRepository: func(r *model.Repository) error { return nil },
		}

		pullReqEvent := makePullRequestEvent(&ts)
		pullReqEvent.Action = github.String("synchronize")

		raw, err := json.Marshal(pullReqEvent)
		require.NoError(t, err)

		event := golambda.Event{
			Origin: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Path:       "/api/v1/webhook/github",
				Headers: map[string]string{
					"X-GitHub-Event": "pull_request",
				},
				Body: string(raw),
			},
		}

		_, err = ctrl.LambdaAPIHandler(event)
		require.NoError(t, err)
		require.NotNil(t, req)
	})

	t.Run("Ignore pull request event not opened or sync", func(t *testing.T) {
		var req *model.ScanRepositoryRequest
		ctrl := controller.New()
		ctrl.Usecase = &MockUsecase{
			sendScanRequest: func(r *model.ScanRepositoryRequest) error {
				req = r
				return nil
			},
			registerRepository: func(r *model.Repository) error { return nil },
		}

		pullReqEvent := makePullRequestEvent(&ts)
		pullReqEvent.Action = github.String("ready_for_review")
		raw, err := json.Marshal(pullReqEvent)
		require.NoError(t, err)

		event := golambda.Event{
			Origin: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Path:       "/api/v1/webhook/github",
				Headers: map[string]string{
					"X-GitHub-Event": "pull_request",
				},
				Body: string(raw),
			},
		}

		_, err = ctrl.LambdaAPIHandler(event)
		require.NoError(t, err)
		require.Nil(t, req)
	})
}

func makePullRequestEvent(ts *time.Time) *github.PullRequestEvent {
	return &github.PullRequestEvent{
		Action: github.String("opened"),
		Repo: &github.Repository{
			Name: github.String("blue"),
			Owner: &github.User{
				Login: github.String("five"),
			},
			HTMLURL:       github.String("https://github-enterprise.example.com/blue/five"),
			DefaultBranch: github.String("default"),
		},
		PullRequest: &github.PullRequest{
			Head: &github.PullRequestBranch{
				SHA:   github.String("xyz"),
				Label: github.String("ao:1"),
			},
			Base: &github.PullRequestBranch{
				Ref: github.String("master"),
			},
			CreatedAt: ts,
		},
		Installation: &github.Installation{
			ID: github.Int64(1234),
		},
	}
}
