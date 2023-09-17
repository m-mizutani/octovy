package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/go-github/v53/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/usecase"
	"github.com/m-mizutani/octovy/pkg/utils"
)

func handleGitHubEvent(uc *usecase.UseCase, r *http.Request, key types.GitHubAppSecret) error {
	payload, err := github.ValidatePayload(r, []byte(key))
	if err != nil {
		return goerr.Wrap(err, "validating payload")
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return goerr.Wrap(err, "parsing webhook")
	}

	ctx := r.Context()
	var octovyCtx *model.Context
	if v, ok := ctx.(*model.Context); !ok {
		octovyCtx = model.NewContext(model.WithBase(ctx))
	} else {
		octovyCtx = v
	}

	scanInput := githubEventToScanInput(event)
	if scanInput == nil {
		return nil
	}

	utils.Logger().With(slog.Any("input", scanInput)).Info("Scan GitHub repository")
	if err := uc.ScanGitHubRepo(octovyCtx, scanInput); err != nil {
		return err
	}

	return nil
}

func githubEventToScanInput(event interface{}) *usecase.ScanGitHubRepoInput {
	switch ev := event.(type) {
	case *github.PushEvent:
		var branch string
		if ref := strings.Split(ev.GetRef(), "/"); len(ref) == 3 && ref[1] == "heads" {
			branch = ref[2]
		}

		return &usecase.ScanGitHubRepoInput{
			GitHubRepoMetadata: usecase.GitHubRepoMetadata{
				Owner:         ev.GetRepo().GetOwner().GetLogin(),
				Repo:          ev.GetRepo().GetName(),
				CommitID:      ev.GetHeadCommit().GetID(),
				Branch:        branch,
				BaseCommitID:  ev.GetBefore(),
				PullRequestID: 0,
			},
			InstallID: types.GitHubAppInstallID(ev.GetInstallation().GetID()),
		}

	case *github.PullRequestEvent:
		if ev.GetAction() != "opened" && ev.GetAction() != "synchronize" {
			utils.Logger().Debug("ignore PR event", slog.String("action", ev.GetAction()))
			return nil
		}

		var branch string
		if ref := strings.Split(ev.GetPullRequest().GetHead().GetRef(), "/"); len(ref) == 3 && ref[1] == "heads" {
			branch = ref[2]
		}

		return &usecase.ScanGitHubRepoInput{
			GitHubRepoMetadata: usecase.GitHubRepoMetadata{
				Owner:         ev.GetRepo().GetOwner().GetLogin(),
				Repo:          ev.GetRepo().GetName(),
				CommitID:      ev.GetPullRequest().GetHead().GetSHA(),
				Branch:        branch,
				BaseCommitID:  ev.GetPullRequest().GetBase().GetSHA(),
				PullRequestID: ev.GetPullRequest().GetNumber(),
			},
			InstallID: types.GitHubAppInstallID(ev.GetInstallation().GetID()),
		}

	case *github.CheckRunEvent, *github.CheckSuiteEvent:
		return nil // ignore

	default:
		utils.Logger().Warn("unsupported event", slog.Any("event", fmt.Sprintf("%T", event)))
		return nil
	}
}
