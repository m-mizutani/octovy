package usecase

import (
	"github.com/google/go-github/v39/github"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
)

type checkRun struct {
	repo    *model.GitHubRepo
	checkID int64
	app     githubapp.Interface
}

func newCheckRun(app githubapp.Interface) *checkRun {
	return &checkRun{
		app: app,
	}
}

func (x *checkRun) create(repo *model.GitHubRepo, commitID string) error {
	checkID, err := x.app.CreateCheckRun(repo, commitID)
	if err != nil {
		return err
	}
	x.repo = repo
	x.checkID = checkID
	return nil
}

func (x *checkRun) complete(ctx *model.Context, scanID string, report *model.Report, frontendURL string, conclusion types.GitHubCheckResult) error {
	ctx.Log().
		With("scanID", scanID).
		With("url", frontendURL).
		Debug("updating check run")

	// TODO: ignore vulnerabilities having status
	body := report.ToMarkdown()
	opt := &github.UpdateCheckRunOptions{
		Name:       "Octovy: package vulnerability check",
		Status:     github.String("completed"),
		Conclusion: github.String(string(conclusion)),
		DetailsURL: github.String(frontendURL + "/scan/" + scanID),
		Output: &github.CheckRunOutput{
			Title:   github.String("Scanned new commit"),
			Summary: github.String("OK"),
			Text:    &body,
		},
	}

	if err := x.app.UpdateCheckRun(x.repo, x.checkID, opt); err != nil {
		return err
	}

	return nil
}
