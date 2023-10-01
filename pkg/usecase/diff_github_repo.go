package usecase

import (
	"database/sql"

	ttypes "github.com/aquasecurity/trivy/pkg/types"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func DiffGitHubRepoResults(ctx *model.Context, dbClient *sql.DB, report *ttypes.Report, meta *GitHubRepoMetadata) error {
	
	return nil
}
