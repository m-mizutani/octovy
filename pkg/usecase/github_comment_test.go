package usecase_test

import (
	"testing"

	"github.com/google/go-github/v39/github"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostGitHubComment(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var cntComment int
		_, mock := githubapp.NewMock()
		mock.CreateIssueCommentMock = func(repo *model.GitHubRepo, prID int, body string) error {
			cntComment++
			assert.Contains(t, body, "CVE-2000-0000")
			assert.Contains(t, body, "ansuz")
			assert.Contains(t, body, "https://octovy.example.com/scan/my-scan-report-id")
			return nil
		}

		input := &usecase.PostGitHubCommentInput{
			App:           mock,
			Target:        &model.ScanTarget{},
			PullReqNumber: github.Int(24),
			Scan: &ent.Scan{
				ID: "my-scan-report-id",
			},
			GitHubEvent: "opened",
			Report: &model.Report{
				Sources: map[string]*model.SourceChanges{
					"x": {
						Added: model.VulnChanges{
							{
								VulnRecord: model.VulnRecord{
									Pkg: &ent.PackageRecord{
										Source: "x",
										Name:   "box",
									},
									Vuln: &ent.Vulnerability{
										ID:    "CVE-2000-0000",
										Title: "ansuz",
									},
								},
							},
						},
					},
				},
			},
			FrontendURL: "https://octovy.example.com/",
		}
		require.NoError(t, usecase.PostGitHubComment(input))
		assert.Equal(t, 1, cntComment)
	})
}
