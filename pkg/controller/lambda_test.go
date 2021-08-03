package controller_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/pkg/controller"
	"github.com/m-mizutani/octovy/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockUsecase struct {
	usecase.Default
	push    *github.PushEvent
	pr      *github.PullRequestEvent
	install *github.InstallationEvent
}

func (x *MockUsecase) HandleGitHubPushEvent(event *github.PushEvent) error {
	x.push = event
	return nil
}

func (x *MockUsecase) HandleGitHubPullReqEvent(event *github.PullRequestEvent) error {
	x.pr = event
	return nil
}

func (x *MockUsecase) HandleGitHubInstallationEvent(event *github.InstallationEvent) error {
	x.install = event
	return nil
}

func TestLambdaAPIWebhook(t *testing.T) {
	t.Run("Push event", func(t *testing.T) {
		mock := &MockUsecase{}
		ctrl := controller.New()
		ctrl.Usecase = mock

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

		require.NotNil(t, mock.push)
		require.Nil(t, mock.pr)
		require.Nil(t, mock.install)
		assert.Equal(t, "blue", *mock.push.Repo.Name)
	})

	t.Run("Pull request event", func(t *testing.T) {
		mock := &MockUsecase{}
		ctrl := controller.New()
		ctrl.Usecase = mock

		pullReqEvent := &github.PullRequestEvent{
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
				Number: github.Int(875),
			},
			Installation: &github.Installation{
				ID: github.Int64(1234),
			},
		}
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
		require.Nil(t, mock.push)
		require.NotNil(t, mock.pr)
		require.Nil(t, mock.install)
		assert.Equal(t, "blue", *mock.pr.Repo.Name)
	})

	t.Run("Installation event", func(t *testing.T) {
		mock := &MockUsecase{}
		ctrl := controller.New()
		ctrl.Usecase = mock

		installationEvent := &github.InstallationEvent{
			Repositories: []*github.Repository{
				{
					Name: github.String("orange"),
				},
			},
		}
		raw, err := json.Marshal(installationEvent)
		require.NoError(t, err)

		event := golambda.Event{
			Origin: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Path:       "/api/v1/webhook/github",
				Headers: map[string]string{
					"X-GitHub-Event": "installation",
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
		require.Nil(t, mock.push)
		require.Nil(t, mock.pr)
		require.NotNil(t, mock.install)
		assert.Equal(t, "orange", *mock.install.Repositories[0].Name)
	})
}
