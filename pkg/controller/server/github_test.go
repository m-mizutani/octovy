package server_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/controller/server"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/usecase"
)

//go:embed testdata/github/pull_request.opened.json
var testGitHubPullRequestOpened []byte

//go:embed testdata/github/pull_request.synchronize.json
var testGitHubPullRequestSynchronize []byte

//go:embed testdata/github/push.json
var testGitHubPush []byte

//go:embed testdata/github/push.default.json
var testGitHubPushDefault []byte

func TestGitHubPullRequestSync(t *testing.T) {
	const secret = "dummy"

	testCases := map[string]struct {
		event string
		body  []byte
		input usecase.ScanGitHubRepoInput
	}{
		"pull_request.opened": {
			event: "pull_request",
			body:  testGitHubPullRequestOpened,
			input: usecase.ScanGitHubRepoInput{
				GitHubRepoMetadata: usecase.GitHubRepoMetadata{
					GitHubCommit: model.GitHubCommit{
						GitHubRepo: model.GitHubRepo{
							Owner: "m-mizutani",
							Repo:  "masq",
						},
						CommitID: "aa0378cad00d375c1897c1b5b5a4dd125984b511",
					},
					PullRequestID:   13,
					Branch:          "update/packages/20230918",
					BaseCommitID:    "8acdc26c9f12b9cc88e5f0b23f082f648d9e5645",
					IsDefaultBranch: false,
				},
				InstallID: 41633205,
			},
		},
		"pull_request.synchronize": {
			event: "pull_request",
			body:  testGitHubPullRequestSynchronize,
			input: usecase.ScanGitHubRepoInput{
				GitHubRepoMetadata: usecase.GitHubRepoMetadata{
					GitHubCommit: model.GitHubCommit{
						GitHubRepo: model.GitHubRepo{
							Owner: "m-mizutani",
							Repo:  "octovy",
						},
						CommitID: "69454c171c2f0f2dbc9ccb0c9ef9b72fd769f046",
					},
					PullRequestID:   89,
					Branch:          "release/v0.2.0",
					BaseCommitID:    "bca5ddd2023d5c906a0420492deb2ede8d99eb79",
					IsDefaultBranch: false,
				},
				InstallID: 41633205,
			},
		},

		"push": {
			event: "push",
			body:  testGitHubPush,
			input: usecase.ScanGitHubRepoInput{
				GitHubRepoMetadata: usecase.GitHubRepoMetadata{
					GitHubCommit: model.GitHubCommit{
						GitHubRepo: model.GitHubRepo{
							Owner: "m-mizutani",
							Repo:  "masq",
						},
						CommitID: "aa0378cad00d375c1897c1b5b5a4dd125984b511",
					},
					PullRequestID:   0,
					Branch:          "update/packages/20230918",
					BaseCommitID:    "0000000000000000000000000000000000000000",
					IsDefaultBranch: false,
				},
				InstallID: 41633205,
			},
		},
		"push: to default": {
			event: "push",
			body:  testGitHubPushDefault,
			input: usecase.ScanGitHubRepoInput{
				GitHubRepoMetadata: usecase.GitHubRepoMetadata{
					GitHubCommit: model.GitHubCommit{
						GitHubRepo: model.GitHubRepo{
							Owner: "m-mizutani",
							Repo:  "ops",
						},
						CommitID: "f58ae7668c3dfc193a1d2c0372cc52847613cde4",
					},
					PullRequestID:   0,
					Branch:          "master",
					BaseCommitID:    "987e1005c2e3c79631b620c4a76afd4b8111b7b1",
					IsDefaultBranch: true,
				},
				InstallID: 41633205,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var called int
			mock := &usecase.Mock{
				MockScanGitHubRepo: func(ctx *model.Context, input *usecase.ScanGitHubRepoInput) error {
					called++
					gt.V(t, input).Equal(&tc.input)
					return nil
				},
			}

			serv := server.New(mock, secret)
			req := newGitHubWebhookRequest(t, tc.event, tc.body, secret)
			w := httptest.NewRecorder()
			serv.Mux().ServeHTTP(w, req)
			gt.V(t, w.Code).Equal(http.StatusOK)
			gt.V(t, called).Equal(1)
		})
	}
}

func newGitHubWebhookRequest(t *testing.T, event string, body []byte, secret types.GitHubAppSecret) *http.Request {
	req := gt.R1(http.NewRequest(http.MethodPost, "/webhook/github", bytes.NewReader(body))).NoError(t)

	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)

	req.Header.Set("X-GitHub-Event", event)
	req.Header.Set("X-Hub-Signature-256", "sha256="+hex.EncodeToString(h.Sum(nil)))
	req.Header.Set("X-GitHub-Delivery", uuid.NewString())
	req.Header.Set("Content-Type", "application/json")

	return req
}
