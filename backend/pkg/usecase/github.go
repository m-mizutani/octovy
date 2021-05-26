package usecase

import (
	"sort"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

func (x *Default) HandleGitHubInstallationEvent(event *github.InstallationEvent) error {
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
		// Do not scan private repository
		if repo.Private != nil && *repo.Private {
			logger.With("repo", repo).Info("Skip private repository")
			continue
		}

		parts := strings.Split(*repo.FullName, "/")
		if len(parts) != 2 {
			return goerr.Wrap(model.ErrInvalidWebhookData, "")
		}
		newRepo := &model.Repository{
			GitHubRepo: model.GitHubRepo{
				Owner:    parts[0],
				RepoName: parts[1],
			},
			URL:           *event.Installation.Account.HTMLURL + "/" + parts[1],
			DefaultBranch: "",
			InstallID:     *event.Installation.ID,
		}

		if err := x.RegisterRepository(newRepo); err != nil {
			return goerr.Wrap(err, "Failed RegisterRepository").With("repo", repo)
		}
	}

	return nil
}

func (x *Default) HandleGitHubPushEvent(event *github.PushEvent) error {
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
		logger.With("event", event).Warn("No commit push")
		return nil
	}
	// Do not scan private repository
	if event.Repo.Private != nil && *event.Repo.Private {
		logger.With("repo", event.Repo).Info("Skip private repository")
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
			CommitID:       *commit.ID,
			UpdatedAt:      commit.Timestamp.Unix(),
			IsPullRequest:  false,
			IsTargetBranch: (*event.Repo.DefaultBranch == refs[2]),
		},
		InstallID: *event.Installation.ID,
	}

	app, err := x.buildGitHubApp(req.InstallID)
	if err != nil {
		return err
	}
	checkID, err := app.CreateCheckRun(&req.GitHubRepo, req.CommitID)
	if err != nil {
		return err
	}
	req.Feedback = &model.FeedbackOptions{CheckID: &checkID}

	if err := x.SendScanRequest(&req); err != nil {
		return goerr.Wrap(err, "Failed SendScanRequest").With("req", req)
	}

	repo := &model.Repository{
		GitHubRepo:    req.GitHubRepo,
		URL:           *event.Repo.HTMLURL,
		DefaultBranch: *event.Repo.DefaultBranch,
		InstallID:     *event.Installation.ID,
	}
	if err := x.RegisterRepository(repo); err != nil {
		return goerr.Wrap(err, "Failed RegisterRepository").With("repo", repo)
	}

	logger.With("event", event).Info("Recv github push event")
	return nil

}

func (x *Default) HandleGitHubPullReqEvent(event *github.PullRequestEvent) error {
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
	if *event.Action != "opened" && *event.Action != "synchronize" {
		return nil
	}

	// Do not scan private repository
	if event.Repo.Private != nil && *event.Repo.Private {
		logger.With("repo", event.Repo).Info("Skip private repository")
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
			CommitID:       *event.PullRequest.Head.SHA,
			UpdatedAt:      event.PullRequest.CreatedAt.Unix(),
			IsPullRequest:  true,
			IsTargetBranch: false,
		},
		InstallID: *event.Installation.ID,
		Feedback: &model.FeedbackOptions{
			PullReqID:     event.PullRequest.Number,
			PullReqBranch: *event.PullRequest.Base.Ref,
		},
	}

	app, err := x.buildGitHubApp(req.InstallID)
	if err != nil {
		return err
	}
	checkID, err := app.CreateCheckRun(&req.GitHubRepo, req.CommitID)
	if err != nil {
		return err
	}
	req.Feedback.CheckID = &checkID

	if err := x.SendScanRequest(&req); err != nil {
		return goerr.Wrap(err, "Failed SendScanRequest").With("req", req)
	}

	repo := &model.Repository{
		GitHubRepo:    req.GitHubRepo,
		URL:           *event.Repo.HTMLURL,
		DefaultBranch: *event.Repo.DefaultBranch,
		InstallID:     *event.Installation.ID,
	}
	if err := x.RegisterRepository(repo); err != nil {
		return goerr.Wrap(err, "Failed RegisterRepository").With("repo", repo)
	}

	logger.With("event", event).Info("Recv github push event")
	return nil

}
