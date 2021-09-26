package usecase

import (
	"context"
	"sort"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

func (x *usecase) HandleGitHubPushEvent(ctx context.Context, event *github.PushEvent) error {
	if event == nil ||
		event.Repo == nil ||
		event.Repo.HTMLURL == nil ||
		event.Repo.DefaultBranch == nil ||
		event.Repo.Name == nil ||
		event.Repo.Owner == nil ||
		event.Repo.Owner.Name == nil ||
		event.Installation == nil ||
		event.Installation.ID == nil {
		return goerr.Wrap(model.ErrInvalidWebhookData, "Not enough fields").With("event", event)
	}

	if len(event.Commits) == 0 {
		logger.Warn().Interface("event", event).Msg("No commit push")
		return nil
	}

	sort.Slice(event.Commits, func(i, j int) bool {
		return event.Commits[i].Timestamp.After(event.Commits[j].Timestamp.Time)
	})
	commit := event.Commits[0]
	refs := strings.Split(*event.Ref, "/")

	req := model.ScanRepositoryRequest{
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    *event.Repo.Owner.Name,
					RepoName: *event.Repo.Name,
				},
				Branch: refs[2],
			},
			CommitID:  *commit.ID,
			UpdatedAt: commit.Timestamp.Unix(),
			URL:       *event.Repo.HTMLURL,
		},
		InstallID: *event.Installation.ID,
	}

	if err := x.SendScanRequest(&req); err != nil {
		return goerr.Wrap(err, "Failed SendScanRequest").With("req", req)
	}

	repo := &ent.Repository{
		Owner:         req.Owner,
		Name:          req.RepoName,
		URL:           *event.Repo.HTMLURL,
		InstallID:     *event.Installation.ID,
		DefaultBranch: event.Repo.DefaultBranch,
		AvatarURL:     event.Repo.Owner.AvatarURL,
	}
	if _, err := x.RegisterRepository(ctx, repo); err != nil {
		return goerr.Wrap(err, "Failed RegisterRepository").With("repo", repo)
	}

	logger.Debug().Interface("event", event).Msg("Recv github push event")
	return nil

}

func (x *usecase) HandleGitHubPullReqEvent(ctx context.Context, event *github.PullRequestEvent) error {
	if event == nil ||
		event.Action == nil ||
		event.Repo == nil ||
		event.Repo.HTMLURL == nil ||
		event.Repo.DefaultBranch == nil ||
		event.Repo.Name == nil ||
		event.Repo.Owner == nil ||
		event.Repo.Owner.Login == nil ||
		event.PullRequest == nil ||
		event.PullRequest.Head == nil ||
		event.PullRequest.Head.SHA == nil ||
		event.PullRequest.Head.Label == nil ||
		event.PullRequest.Base == nil ||
		event.PullRequest.Base.Ref == nil ||
		event.PullRequest.CreatedAt == nil ||
		event.PullRequest.Number == nil ||
		event.Installation == nil ||
		event.Installation.ID == nil {
		return goerr.Wrap(model.ErrInvalidWebhookData, "Not enough fields").With("event", event)
	}

	// Check only PR opened and synchronize
	var targetBranch string
	switch *event.Action {
	case "opened":
		targetBranch = *event.PullRequest.Base.Ref
	case "synchronize":
		targetBranch = *event.PullRequest.Head.Label
	default:
		return nil
	}

	req := model.ScanRepositoryRequest{
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    *event.Repo.Owner.Login,
					RepoName: *event.Repo.Name,
				},
				Branch: *event.PullRequest.Head.Label,
			},
			CommitID:     *event.PullRequest.Head.SHA,
			UpdatedAt:    event.PullRequest.CreatedAt.Unix(),
			URL:          *event.Repo.HTMLURL,
			TargetBranch: targetBranch,
		},
		InstallID:     *event.Installation.ID,
		PullReqNumber: event.PullRequest.Number,
		PullReqAction: *event.Action,
	}

	if err := x.SendScanRequest(&req); err != nil {
		return goerr.Wrap(err, "Failed SendScanRequest").With("req", req)
	}

	repo := &ent.Repository{
		Owner:         req.Owner,
		Name:          req.RepoName,
		URL:           *event.Repo.HTMLURL,
		InstallID:     *event.Installation.ID,
		DefaultBranch: event.Repo.DefaultBranch,
		AvatarURL:     event.Repo.Owner.AvatarURL,
	}

	if _, err := x.RegisterRepository(ctx, repo); err != nil {
		return goerr.Wrap(err, "Failed RegisterRepository").With("repo", repo)
	}

	logger.Debug().Interface("event", event).Msg("Recv github PR event")
	return nil
}

func (x *usecase) HandleGitHubInstallationEvent(ctx context.Context, event *github.InstallationEvent) error {
	if event == nil ||
		event.Installation == nil ||
		event.Installation.ID == nil ||
		event.Installation.Account == nil ||
		event.Installation.Account.HTMLURL == nil {
		return goerr.Wrap(model.ErrInvalidWebhookData, "Not enough event fields").With("event", event)
	}

	for _, repo := range event.Repositories {
		if repo == nil || repo.FullName == nil {
			return goerr.Wrap(model.ErrInvalidWebhookData, "Not enough repository fields").With("repo", repo)
		}

		parts := strings.Split(*repo.FullName, "/")
		if len(parts) != 2 {
			return goerr.Wrap(model.ErrInvalidWebhookData, "")
		}
		repo := &ent.Repository{
			Owner:     parts[0],
			Name:      parts[1],
			URL:       *event.Installation.Account.HTMLURL,
			InstallID: *event.Installation.ID,
		}

		if _, err := x.RegisterRepository(ctx, repo); err != nil {
			return goerr.Wrap(err, "Failed RegisterRepository").With("repo", repo)
		}
	}

	logger.Debug().Interface("event", event).Msg("Recv github installation event")
	return nil
}
