package usecase

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/utils"
)

type postGitHubCommentInput struct {
	App           githubapp.Interface
	GitHubEvent   string
	Target        *model.ScanTarget
	PullReqNumber *int
	Report        *model.Report
}

func postGitHubComment(input *postGitHubCommentInput) error {
	if input.PullReqNumber == nil {
		return goerr.Wrap(model.ErrInvalidSystemValue).With("input", input)
	}

	if input.Report.NothingToNotify(input.GitHubEvent) {
		utils.Logger.Debug("nothing to notify, returning")
		return nil
	}

	body := input.Report.ToMarkdown()
	if err := input.App.CreateIssueComment(&input.Target.GitHubRepo, *input.PullReqNumber, body); err != nil {
		return err
	}

	return nil
}
