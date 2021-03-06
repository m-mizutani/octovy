package usecase

import (
	"sort"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/utils"
)

func (x *Usecase) HandleGitHubPushEvent(ctx *model.Context, event *github.PushEvent) error {
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
		ctx.Log().With("event", event).Warn("No commit push")
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
					Owner: *event.Repo.Owner.Name,
					Name:  *event.Repo.Name,
				},
				Branch: refs[2],
			},
			CommitID:  *commit.ID,
			UpdatedAt: commit.Timestamp.Unix(),
			URL:       *event.Repo.HTMLURL,
		},
		InstallID: *event.Installation.ID,
	}

	if err := x.Scan(ctx, &req); err != nil {
		return goerr.Wrap(err, "Failed Scan").With("req", req)
	}

	repo := &ent.Repository{
		Owner:         req.Owner,
		Name:          req.Name,
		URL:           *event.Repo.HTMLURL,
		InstallID:     *event.Installation.ID,
		DefaultBranch: event.Repo.DefaultBranch,
		AvatarURL:     event.Repo.Owner.AvatarURL,
	}
	if _, err := x.RegisterRepository(ctx, repo); err != nil {
		return goerr.Wrap(err, "Failed RegisterRepository").With("repo", repo)
	}

	ctx.Log().With("event", event).Debug("Recv github push event")
	return nil

}

func (x *Usecase) HandleGitHubPullReqEvent(ctx *model.Context, event *github.PullRequestEvent) error {
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
					Owner: *event.Repo.Owner.Login,
					Name:  *event.Repo.Name,
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

	if err := x.Scan(ctx, &req); err != nil {
		return goerr.Wrap(err, "Failed Scan").With("req", req)
	}

	repo := &ent.Repository{
		Owner:         req.Owner,
		Name:          req.Name,
		URL:           *event.Repo.HTMLURL,
		InstallID:     *event.Installation.ID,
		DefaultBranch: event.Repo.DefaultBranch,
		AvatarURL:     event.Repo.Owner.AvatarURL,
	}

	if _, err := x.RegisterRepository(ctx, repo); err != nil {
		return goerr.Wrap(err, "Failed RegisterRepository").With("repo", repo)
	}

	ctx.Log().With("event", event).Debug("Recv github PR event")
	return nil
}

func (x *Usecase) HandleGitHubInstallationEvent(ctx *model.Context, event *github.InstallationEvent) error {
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

	ctx.Log().With("event", event).Debug("Recv github installation event")
	return nil
}

func (x *Usecase) VerifyGitHubSecret(sigSHA256 string, body []byte) error {
	if x.config.GitHubWebhookSecret == "" {
		if sigSHA256 == "" {
			return nil // No secret and no signature
		}
		utils.Logger.With("signature", sigSHA256).Warn("Got X-Hub-Signature-256, but no secret is configured. Octovy ignore the signature and continue processing")
		return nil
	}

	if err := github.ValidateSignature(sigSHA256, body, []byte(x.config.GitHubWebhookSecret)); err != nil {
		return model.ErrInvalidWebhookData.Wrap(err)
	}

	return nil
}
