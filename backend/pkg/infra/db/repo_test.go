package db_test

import (
	"testing"

	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepo(t *testing.T) {
	t.Run("Update default branch status", func(t *testing.T) {
		client := newTestTable(t)
		repo := &model.Repository{
			GitHubRepo:    model.GitHubRepo{},
			URL:           "https://xxx",
			DefaultBranch: "blue",
			InstallID:     9,
		}
		branch := &model.Branch{
			GitHubBranch: model.GitHubBranch{
				Branch: "blue",
			},
			LastScannedAt: 1234,
			ReportSummary: model.ScanReportSummary{
				PkgTypes:     []model.PkgType{model.PkgBundler},
				PkgCount:     3,
				VulnCount:    2,
				VulnPkgCount: 1,
			},
		}
		inserted, err := client.InsertRepo(repo)
		require.NoError(t, err)
		assert.True(t, inserted)
		require.NoError(t, client.UpdateBranchIfDefault(&repo.GitHubRepo, branch))

		r1, err := client.FindRepoByFullName(repo.Owner, repo.RepoName)
		require.NoError(t, err)
		assert.Equal(t, r1.Branch, *branch)
	})

	t.Run("Do not update if not default branch", func(t *testing.T) {
		client := newTestTable(t)
		repo := &model.Repository{
			GitHubRepo:    model.GitHubRepo{},
			URL:           "https://xxx",
			DefaultBranch: "blue",
			InstallID:     9,
		}
		branch := &model.Branch{
			GitHubBranch: model.GitHubBranch{
				Branch: "blue",
			},
			LastScannedAt: 1234,
			ReportSummary: model.ScanReportSummary{
				PkgTypes:     []model.PkgType{model.PkgBundler},
				PkgCount:     3,
				VulnCount:    2,
				VulnPkgCount: 1,
			},
		}
		inserted, err := client.InsertRepo(repo)
		require.NoError(t, err)

		// Change default branch name
		require.NoError(t, client.SetRepoDefaultBranchName(&repo.GitHubRepo, "main"))

		assert.True(t, inserted)
		require.NoError(t, client.UpdateBranchIfDefault(&repo.GitHubRepo, branch))

		r1, err := client.FindRepoByFullName(repo.Owner, repo.RepoName)
		require.NoError(t, err)
		assert.NotEqual(t, r1.Branch, *branch)
	})
}
