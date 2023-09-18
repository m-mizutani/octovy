package usecase

import "github.com/m-mizutani/octovy/pkg/domain/model"

type Mock struct {
	MockScanGitHubRepo func(ctx *model.Context, input *ScanGitHubRepoInput) error
}

func NewMock() UseCase {
	return &Mock{}
}

func (x *Mock) ScanGitHubRepo(ctx *model.Context, input *ScanGitHubRepoInput) error {
	return x.MockScanGitHubRepo(ctx, input)
}
