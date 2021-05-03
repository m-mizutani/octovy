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
	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/m-mizutani/octovy/backend/pkg/service"
	"github.com/m-mizutani/octovy/backend/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockUsecase struct {
	usecase.Default
	sendScanRequest    func(svc *service.Service, r *model.ScanRepositoryRequest) error
	registerRepository func(svc *service.Service, repo *model.Repository) error
}

func (x *MockUsecase) SendScanRequest(svc *service.Service, r *model.ScanRepositoryRequest) error {
	return x.sendScanRequest(svc, r)
}

func (x *MockUsecase) RegisterRepository(svc *service.Service, repo *model.Repository) error {
	return x.registerRepository(svc, repo)
}

func TestLambdaAPI(t *testing.T) {
	ts := time.Now().UTC()
	t.Run("GHE Webhook", func(t *testing.T) {
		var req *model.ScanRepositoryRequest
		var repo *model.Repository
		ctrl := controller.New()
		ctrl.Usecase = &MockUsecase{
			sendScanRequest: func(svc *service.Service, r *model.ScanRepositoryRequest) error {
				req = r
				return nil
			},
			registerRepository: func(svc *service.Service, r *model.Repository) error {
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
		assert.Equal(t, "beefcafe", req.Ref)
		assert.Equal(t, ts.Add(time.Minute).Unix(), req.UpdatedAt)

		require.NotNil(t, repo)
	})
}
