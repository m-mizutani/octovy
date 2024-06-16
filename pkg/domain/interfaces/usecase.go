package interfaces

import (
	"context"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/model/trivy"
)

type UseCase interface {
	InsertScanResult(ctx context.Context, meta model.GitHubMetadata, report trivy.Report, cfg model.Config) error
	ScanGitHubRepo(ctx context.Context, input *model.ScanGitHubRepoInput) error
}
