package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/go-github/v53/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/utils"
)

func handleGitHubAppEvent(uc interfaces.UseCase, r *http.Request, key types.GitHubAppSecret) error {
	ctx := r.Context()
	payload, err := github.ValidatePayload(r, []byte(key))
	if err != nil {
		return goerr.Wrap(err, "validating payload")
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return goerr.Wrap(err, "parsing webhook")
	}

	utils.CtxLogger(ctx).With(slog.Any("event", event)).Info("Received GitHub App event")

	scanInput := githubEventToScanInput(event)
	if scanInput == nil {
		return nil
	}

	utils.Logger().With(slog.Any("input", scanInput)).Info("Scan GitHub repository")

	if err := uc.ScanGitHubRepo(r.Context(), scanInput); err != nil {
		return goerr.Wrap(err, "failed to scan GitHub repository")
	}

	return nil
}

func refToBranch(v string) string {
	if ref := strings.SplitN(v, "/", 3); len(ref) == 3 && ref[0] == "refs" && ref[1] == "heads" {
		return ref[2]
	}
	return v
}

func githubEventToScanInput(event interface{}) *model.ScanGitHubRepoInput {
	switch ev := event.(type) {
	case *github.PushEvent:
		return &model.ScanGitHubRepoInput{
			GitHubMetadata: model.GitHubMetadata{
				GitHubCommit: model.GitHubCommit{
					GitHubRepo: model.GitHubRepo{
						RepoID:   ev.GetRepo().GetID(),
						Owner:    ev.GetRepo().GetOwner().GetLogin(),
						RepoName: ev.GetRepo().GetName(),
					},
					CommitID: ev.GetHeadCommit().GetID(),
					Branch:   refToBranch(ev.GetRef()),
					Ref:      ev.GetRef(),
					Committer: model.GitHubUser{
						Login: ev.GetHeadCommit().GetCommitter().GetLogin(),
						Email: ev.GetHeadCommit().GetCommitter().GetEmail(),
					},
				},
				DefaultBranch: ev.GetRepo().GetDefaultBranch(),
			},
			InstallID: types.GitHubAppInstallID(ev.GetInstallation().GetID()),
		}

	case *github.PullRequestEvent:
		if ev.GetAction() != "opened" && ev.GetAction() != "synchronize" {
			utils.Logger().Debug("ignore PR event", slog.String("action", ev.GetAction()))
			return nil
		}

		pr := ev.GetPullRequest()

		input := &model.ScanGitHubRepoInput{
			GitHubMetadata: model.GitHubMetadata{
				GitHubCommit: model.GitHubCommit{
					GitHubRepo: model.GitHubRepo{
						RepoID:   ev.GetRepo().GetID(),
						Owner:    ev.GetRepo().GetOwner().GetLogin(),
						RepoName: ev.GetRepo().GetName(),
					},
					CommitID: pr.GetHead().GetSHA(),
					Ref:      pr.GetHead().GetRef(),
					Branch:   pr.GetHead().GetRef(),
					Committer: model.GitHubUser{
						ID:    pr.GetHead().GetUser().GetID(),
						Login: pr.GetHead().GetUser().GetLogin(),
						Email: pr.GetHead().GetUser().GetEmail(),
					},
				},
				DefaultBranch: ev.GetRepo().GetDefaultBranch(),
				PullRequest: &model.GitHubPullRequest{
					ID:           pr.GetID(),
					Number:       pr.GetNumber(),
					BaseBranch:   pr.GetBase().GetRef(),
					BaseCommitID: pr.GetBase().GetSHA(),
					User: model.GitHubUser{
						ID:    pr.GetBase().GetUser().GetID(),
						Login: pr.GetBase().GetUser().GetLogin(),
						Email: pr.GetBase().GetUser().GetEmail(),
					},
				},
			},
			InstallID: types.GitHubAppInstallID(ev.GetInstallation().GetID()),
		}

		return input

	case *github.InstallationEvent, *github.InstallationRepositoriesEvent:
		return nil // ignore

	default:
		utils.Logger().Warn("unsupported event", slog.Any("event", fmt.Sprintf("%T", event)))
		return nil
	}
}

func handleGitHubActionEvent(_ interfaces.UseCase, _ *http.Request) error {
	return nil
}
