package usecase_test

import (
	"context"
	"os"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/mock"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/model/trivy"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/usecase"
)

func TestRenderScanReport(t *testing.T) {
	report := trivy.Report{
		Results: []trivy.Result{
			{
				Target: "target1",
				Vulnerabilities: []trivy.DetectedVulnerability{
					{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1", Vulnerability: trivy.Vulnerability{Title: "Vuln title1", Severity: "HIGH"}},
					{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2", Vulnerability: trivy.Vulnerability{Title: "Vuln title2", Severity: "CRITICAL"}},
				},
			},
			{
				Target: "target2",
				Vulnerabilities: []trivy.DetectedVulnerability{
					{VulnerabilityID: "CVE-0000-0003", PkgName: "pkg4", Vulnerability: trivy.Vulnerability{Title: "Vuln title3", Severity: "CRITICAL"}},
				},
			},
			{
				Target: "target3",
				Secrets: []trivy.SecretFinding{
					{
						RuleID:    "slack-web-hook",
						Category:  "Slack",
						Severity:  "HIGH",
						Title:     "Slack Web Hook",
						StartLine: 14,
						EndLine:   15,
					},
				},
			},
		},
	}
	added := trivy.Results{
		{
			Target: "target1",
			Vulnerabilities: []trivy.DetectedVulnerability{
				{
					VulnerabilityID: "CVE-0000-0002",
					PkgName:         "pkg2",
					Vulnerability: trivy.Vulnerability{
						Title:       "Vuln title2",
						Description: "Vuln description2",
						Severity:    "CRITICAL",
						References: []string{
							"https://example.com",
							"https://example.com/2",
						},
					},
				},
			},
		},
	}
	fixed := trivy.Results{
		{
			Target: "target2",
			Vulnerabilities: []trivy.DetectedVulnerability{
				{
					VulnerabilityID: "CVE-0000-0003",
					PkgName:         "pkg3",
					Vulnerability: trivy.Vulnerability{
						Title: "Vuln title3",
					},
				},
			},
		},
	}

	body, err := usecase.RenderScanReport(&report, added, fixed)
	gt.NoError(t, err)
	gt.NoError(t, os.WriteFile("templates/test_comment_body.md", []byte(body), 0644))
}

func TestIgnoreIfNoResults(t *testing.T) {
	report := trivy.Report{
		SchemaVersion: 1,
		ArtifactName:  "test",
	}

	csMock := mock.StorageMock{}
	ghMock := mock.GitHubMock{}
	uc := usecase.New(infra.New(
		infra.WithGitHubApp(&ghMock),
		infra.WithStorage(&csMock),
	))
	input := &model.ScanGitHubRepoInput{
		GitHubMetadata: model.GitHubMetadata{
			GitHubCommit: model.GitHubCommit{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "magic",
					RepoID:   12345,
				},
				Committer: model.GitHubUser{Login: "octovy-bot"},
				CommitID:  "9b7cea90596429d5b1243caecc15b1f79598cb85",
				Branch:    "main",
				Ref:       "refs/pull/123/merge",
			},
			PullRequest: &model.GitHubPullRequest{Number: 123},
		},
		InstallID: 12345,
	}

	ctx := context.Background()
	gt.NoError(t, uc.CommentGitHubPR(ctx, input, &report, &model.Config{}))
}

func TestHideGitHubOldComments(t *testing.T) {
	mockGH := &mock.GitHubMock{}

	uc := usecase.New(infra.New(
		infra.WithGitHubApp(mockGH),
	))

	type testCase struct {
		comments   []*model.GitHubIssueComment
		subjectIDs []string
	}

	runTest := func(tc testCase) func(t *testing.T) {
		return func(t *testing.T) {
			mockGH.ListIssueCommentsFunc = func(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int) ([]*model.GitHubIssueComment, error) {
				return tc.comments, nil
			}

			var minimized []string
			mockGH.MinimizeCommentFunc = func(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, subjectID string) error {
				minimized = append(minimized, subjectID)
				return nil
			}

			input := &model.ScanGitHubRepoInput{
				GitHubMetadata: model.GitHubMetadata{
					GitHubCommit: model.GitHubCommit{
						GitHubRepo: model.GitHubRepo{
							Owner:    "blue",
							RepoName: "magic",
						},
					},
					PullRequest: &model.GitHubPullRequest{Number: 123},
				},
				InstallID: 12345,
			}

			ctx := context.Background()
			gt.NoError(t, uc.HideGitHubOldComments(ctx, input))
			gt.V(t, minimized).Equal(tc.subjectIDs)
		}
	}

	t.Run("no comments", runTest(testCase{}))

	t.Run("no minimized comments without signature", runTest(testCase{
		comments: []*model.GitHubIssueComment{
			{ID: "abc", Body: "comment1", IsMinimized: false},
			{ID: "edf", Body: "comment2", IsMinimized: true},
		},
		subjectIDs: nil,
	}))

	t.Run("minimize comments with signature", runTest(testCase{
		comments: []*model.GitHubIssueComment{
			{ID: "abc", Body: types.GitHubCommentSignature + "\ncomment1", IsMinimized: false},
			{ID: "edf", Body: types.GitHubCommentSignature + "\ncomment2", IsMinimized: true},
		},
		subjectIDs: []string{"abc"},
	}))
}
