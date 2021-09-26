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
	Target        *model.ScanTarget
	PullReqNumber *int
	Scan          *ent.Scan
	Changes       vulnChanges
	FrontendURL   string
	DB            *vulnStatusDB
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

	q := input.Changes.Qualified(input.DB)
	qAdded := q.FilterByType(vulnAdded)
	qRemained := q.FilterByType(vulnRemained)
	rDeleted := input.Changes.FilterByType(vulnDeleted)

	if len(qAdded) == 0 && len(qRemained) == 0 && len(rDeleted) == 0 {
		logger.Debug().Interface("input", input).Msg("nothing to notify, returning")
		return nil
	}

	var b githubCommentBody
	b.Add("## Octovy scan result")
	b.Break()

	for _, src := range input.Changes.Sources() {
		b.Add("### "+src, "")
		b.Break()

		for _, v := range qAdded.FilterBySource(src) {
			b.Add("- 🚨 **New** %s (%s): %s", v.Vuln.ID, v.Pkg.Name, v.Vuln.Title)
		}
		for _, v := range rDeleted.FilterBySource(src) {
			b.Add("- ✅ **Fixed** %s (%s): %s", v.Vuln.ID, v.Pkg.Name, v.Vuln.Title)
		}
		if remained := qRemained.FilterBySource(src); len(remained) > 0 {
			b.Add("- ⚠️ %d vulnerabilities remained", len(remained))
		}
		b.Break()
	}

	b.Add("🗒️ See [report](%s/scan/%s) more detail", input.FrontendURL, input.Scan.ID)

	if err := input.App.CreateIssueComment(&input.Target.GitHubRepo, *input.PullReqNumber, b.Join()); err != nil {
		return err
	}

	return nil
}
