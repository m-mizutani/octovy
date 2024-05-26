package usecase

import (
	"context"

	"github.com/m-mizutani/octovy/pkg/domain/model"
)

var RenderScanReport = renderScanReport

func (x *UseCase) HideGitHubOldComments(ctx context.Context, input *model.ScanGitHubRepoInput) error {
	return x.hideGitHubOldComments(ctx, input)
}
