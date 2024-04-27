package usecase

import (
	"context"

	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/model/trivy"
)

type Mock struct {
	MockInsertScanResult func(ctx context.Context, meta model.GitHubMetadata, report trivy.Report) error
	MockScanGitHubRepo   func(ctx context.Context, input *model.ScanGitHubRepoInput) error
}

func NewMock() *Mock {
	return &Mock{}
}

var _ interfaces.UseCase = &Mock{}

func (x *Mock) InsertScanResult(ctx context.Context, meta model.GitHubMetadata, report trivy.Report) error {
	return x.MockInsertScanResult(ctx, meta, report)
}

// ScanGitHubRepo implements interfaces.UseCase.
func (x *Mock) ScanGitHubRepo(ctx context.Context, input *model.ScanGitHubRepoInput) error {
	return x.MockScanGitHubRepo(ctx, input)
}
