package usecase

import (
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

func (x *githubCommentBody) Add(s ...string) { x.lines = append(x.lines, s...) }
func (x *githubCommentBody) Join() string    { return strings.Join(x.lines, "\n") }

func postGitHubComment(input *postGitHubCommentInput) error {
	if input.PullReqNumber == nil {
		return goerr.Wrap(model.ErrInvalidSystemValue).With("input", input)
	}

	q := input.Changes.Qualified(input.DB)
	qAdded := q.Filter(vulnAdded)
	qRemained := q.Filter(vulnRemained)
	rDeleted := input.Changes.Filter(vulnDeleted)

	if len(qAdded) == 0 && len(qRemained) == 0 && len(rDeleted) == 0 {
		// Nothing to comment
		return nil
	}

	var b githubCommentBody
	b.Add("## Octovy scan result", "")

	if err := input.App.CreateIssueComment(&input.Target.GitHubRepo, *input.PullReqNumber, b.Join()); err != nil {
		return err
	}

	return nil
}
