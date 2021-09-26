package usecase

import (
	"fmt"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
)

type postGitHubCommentInput struct {
	App           githubapp.Interface
	GitHubEvent   string
	Target        *model.ScanTarget
	PullReqNumber *int
	Scan          *ent.Scan
	Report        *model.Report
	FrontendURL   string
}

type githubCommentBody struct {
	lines []string
}

func (x *githubCommentBody) Add(f string, v ...interface{}) {
	x.lines = append(x.lines, fmt.Sprintf(f, v...))
}
func (x *githubCommentBody) Break() {
	x.lines = append(x.lines, "")
}
func (x *githubCommentBody) Join() string { return strings.Join(x.lines, "\n") }

func postGitHubComment(input *postGitHubCommentInput) error {
	if input.PullReqNumber == nil {
		return goerr.Wrap(model.ErrInvalidSystemValue).With("input", input)
	}

	if input.Report.NothingToNotify(input.GitHubEvent) {
		logger.Debug().Msg("nothing to notify, returning")
		return nil
	}

	var b githubCommentBody
	b.Add("## Octovy scan result")
	b.Break()

	for src, changes := range input.Report.Sources {
		b.Add("### " + src)
		b.Break()

		for _, v := range changes.Added {
			b.Add("- 🚨 **New** %s (%s): %s", v.Vuln.ID, v.Pkg.Name, v.Vuln.Title)
		}
		for _, v := range changes.Deleted {
			b.Add("- ✅ **Fixed** %s (%s): %s", v.Vuln.ID, v.Pkg.Name, v.Vuln.Title)
		}
		if len(changes.Remained) > 0 {
			b.Add("- ⚠️ %d vulnerabilities are remained", len(changes.Remained))
		}
		b.Break()
	}

	b.Add("🗒️ See [report](%s/scan/%s) more detail", strings.Trim(input.FrontendURL, "/"), input.Scan.ID)

	if err := input.App.CreateIssueComment(&input.Target.GitHubRepo, *input.PullReqNumber, b.Join()); err != nil {
		return err
	}

	return nil
}
