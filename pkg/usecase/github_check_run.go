package usecase

import (
	"github.com/google/go-github/v39/github"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
)

type checkRun struct {
	repo      *model.GitHubRepo
	checkID   int64
	app       githubapp.Interface
	completed bool
}

func newCheckRun(app githubapp.Interface) *checkRun {
	return &checkRun{
		app: app,
	}
}

func (x *checkRun) create(ctx *model.Context, repo *model.GitHubRepo, commitID string) error {
	checkID, err := x.app.CreateCheckRun(repo, commitID)
	if err != nil {
		return err
	}
	x.repo = repo
	x.checkID = checkID
	ctx.Log().With("checkID", checkID).Debug("created github check")
	return nil
}

func str(s string) *string { return &s }

func (x *checkRun) fallback(ctx *model.Context) {
	if x.completed {
		return
	}

	ctx.Log().Debug("Complete to close check run")

	opt := &github.UpdateCheckRunOptions{
		Name:       "Octovy: package vulnerability check",
		Status:     github.String("completed"),
		Conclusion: github.String("success"),
		Output: &github.CheckRunOutput{
			Title:   github.String("❗ Octovy got error, but success to avoid CI failure"),
			Summary: github.String("Failed scan procedure"),
			Text:    github.String("Please contact to administrator of Octovy"),
		},
	}

	if err := x.app.UpdateCheckRun(x.repo, x.checkID, opt); err != nil {
		ctx.Log().With("opt", opt).Err(err).Error("failed to submit fallback")
	}
}

func (x *checkRun) complete(ctx *model.Context, scanID string, report *model.Report, frontendURL string, result *model.GitHubCheckResult) error {
	ctx.Log().
		With("scanID", scanID).
		With("url", frontendURL).
		Debug("updating check run")

	title := "Package scan report"

	// TODO: use result.Messages
	body := report.ToMarkdown()
	opt := &github.UpdateCheckRunOptions{
		Name:       "Octovy",
		Status:     str("completed"),
		Conclusion: str(string(result.Conclusion)),
		DetailsURL: str(frontendURL + "/scan/" + scanID),
		Output: &github.CheckRunOutput{
			Title:   &title,
			Summary: str(report.Summary()),
			Text:    &body,
		},
	}

	if err := x.app.UpdateCheckRun(x.repo, x.checkID, opt); err != nil {
		return err
	}
	ctx.Log().With("opt", opt).Debug("update github check")
	x.completed = true
	return nil
}
