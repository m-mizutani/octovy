package usecase

import "github.com/m-mizutani/octovy/pkg/infra/githubapp"

func postGitHubComment(app githubapp.Interface, scanID string, changes *pkgChanges, frontendURL string) error {
	// TODO: feedback
	logger.Info().Msg("DO Feedback")

	return nil
}
