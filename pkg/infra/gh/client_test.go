package gh_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/gh"
	"github.com/m-mizutani/octovy/pkg/utils"
)

func TestGitHubComment(t *testing.T) {
	ghApp, installID := buildGitHubApp(t)

	ctx := context.Background()
	repo := &model.GitHubRepo{
		Owner:    "m-mizutani",
		RepoName: "octovy-test-code",
	}
	ghApp.CreateIssueComment(ctx, repo, types.GitHubAppInstallID(installID), 1, "Hello, world")

	comments := gt.R1(ghApp.ListIssueComments(ctx, repo, types.GitHubAppInstallID(installID), 1)).NoError(t)

	utils.DumpJson(t, "comments.json", comments)
}

func buildGitHubApp(t *testing.T) (*gh.Client, types.GitHubAppInstallID) {
	var (
		strAppID          = utils.LoadEnv(t, "TEST_GITHUB_APP_ID")
		strInstallationID = utils.LoadEnv(t, "TEST_GITHUB_INSTALLATION_ID")
		privateKey        = utils.LoadEnv(t, "TEST_GITHUB_PRIVATE_KEY")
	)

	appID := gt.R1(strconv.ParseInt(strAppID, 10, 64)).NoError(t)
	installID := gt.R1(strconv.ParseInt(strInstallationID, 10, 64)).NoError(t)

	ghApp := gt.R1(gh.New(types.GitHubAppID(appID), types.GitHubAppPrivateKey(privateKey))).NoError(t)

	return ghApp, types.GitHubAppInstallID(installID)
}

func TestListComments(t *testing.T) {
	ghApp, installID := buildGitHubApp(t)

	ctx := context.Background()
	repo := &model.GitHubRepo{
		Owner:    "m-mizutani",
		RepoName: "octovy-test-code",
	}

	comments, err := ghApp.ListIssueComments(ctx, repo, types.GitHubAppInstallID(installID), 2)
	gt.NoError(t, err)
	gt.A(t, comments).Longer(1).At(0, func(t testing.TB, v *model.GitHubIssueComment) {
		gt.Equal(t, v.Body, "testing")
	})
}

func TestHideComment(t *testing.T) {
	ghApp, installID := buildGitHubApp(t)

	ctx := context.Background()
	repo := &model.GitHubRepo{
		Owner:    "m-mizutani",
		RepoName: "octovy-test-code",
	}
	testIssueID := 2

	slag := "comment-test:" + uuid.NewString()
	gt.NoError(t, ghApp.CreateIssueComment(ctx, repo, installID, 2, slag))

	comments, err := ghApp.ListIssueComments(ctx, repo, types.GitHubAppInstallID(installID), testIssueID)
	gt.NoError(t, err)

	var subjectID string
	for _, c := range comments {
		if c.Body == slag {
			gt.False(t, c.IsMinimized)
			subjectID = c.ID
			break
		}
	}

	gt.NotEqual(t, subjectID, "")

	gt.NoError(t, ghApp.MinimizeComment(ctx, repo, installID, subjectID))

	comments, err = ghApp.ListIssueComments(ctx, repo, types.GitHubAppInstallID(installID), testIssueID)
	gt.NoError(t, err)
	gt.A(t, comments).Longer(1).Any(func(v *model.GitHubIssueComment) bool {
		return v.ID == subjectID && v.IsMinimized
	})
}
