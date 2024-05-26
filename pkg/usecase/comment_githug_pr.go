package usecase

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/logic"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/model/trivy"
	"github.com/m-mizutani/octovy/pkg/domain/types"
)

func (x *UseCase) CommentGitHubPR(ctx context.Context, input *model.ScanGitHubRepoInput, report *trivy.Report) error {
	if err := input.Validate(); err != nil {
		return err
	}

	if nil == input.GitHubMetadata.PullRequest {
		return goerr.New("PullRequest is not set")
	}

	if x.clients.GitHubApp() == nil {
		return goerr.New("GitHubApp client is not set")
	}
	if x.clients.Storage() == nil {
		return goerr.New("Storage client is not configured")
	}

	var added, fixed trivy.Results
	target := model.GitHubMetadata{
		GitHubCommit: model.GitHubCommit{
			GitHubRepo: input.GitHubMetadata.GitHubRepo,
			CommitID:   input.GitHubMetadata.PullRequest.BaseCommitID,
		},
	}

	commitKey := toStorageCommitKey(target)
	r, err := x.clients.Storage().Get(ctx, commitKey)
	if err != nil {
		return err
	} else if r != nil {
		defer r.Close()

		var oldScan model.Scan
		if err := json.NewDecoder(r).Decode(&oldScan); err != nil {
			return goerr.Wrap(err, "Failed to decode old scan result")
		}
		fixed, added = logic.DiffResults(&oldScan.Report, report)
	}

	body, err := renderScanReport(report, added, fixed)
	if err != nil {
		return err
	}

	if err := x.hideGitHubOldComments(ctx, input); err != nil {
		return err
	}

	if err := x.clients.GitHubApp().CreateIssueComment(ctx, &input.GitHubMetadata.GitHubRepo, input.InstallID, input.PullRequest.Number, body); err != nil {
		return err
	}

	return nil
}

func (x *UseCase) hideGitHubOldComments(ctx context.Context, input *model.ScanGitHubRepoInput) error {
	if nil == input.GitHubMetadata.PullRequest {
		return goerr.New("PullRequest is not set")
	}

	if x.clients.GitHubApp() == nil {
		return goerr.New("GitHubApp client is not set")
	}

	comments, err := x.clients.GitHubApp().ListIssueComments(ctx, &input.GitHubMetadata.GitHubRepo, input.InstallID, input.PullRequest.Number)
	if err != nil {
		return err
	}

	for _, comment := range comments {
		if !comment.IsMinimized && strings.HasPrefix(comment.Body, types.GitHubCommentSignature) {
			if err := x.clients.GitHubApp().MinimizeComment(ctx, &input.GitHubMetadata.GitHubRepo, input.InstallID, comment.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

type scanReport struct {
	Signature string
	Metadata  scanReportMetadata
	Report    *trivy.Report
	Added     trivy.Results
	Fixed     trivy.Results
}

type scanReportMetadata struct {
	TotalVulnCount   int
	FixableVulnCount int
}

//go:embed templates/comment_body.md
var commentBodyTemplate string

func renderScanReport(report *trivy.Report, added, fixed trivy.Results) (string, error) {
	data := scanReport{
		Signature: types.GitHubCommentSignature,
		Report:    report,
		Added:     added,
		Fixed:     fixed,
	}

	for _, result := range report.Results {
		for _, vuln := range result.Vulnerabilities {
			data.Metadata.TotalVulnCount++
			if vuln.FixedVersion != "" {
				data.Metadata.FixableVulnCount++
			}
		}
	}

	tmpl, err := template.New("scanReport").Parse(commentBodyTemplate)
	if err != nil {
		return "", goerr.Wrap(err, "failed to parse comment body template")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", goerr.Wrap(err, "failed to render comment body template")
	}

	return buf.String(), nil
}
