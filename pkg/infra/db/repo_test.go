package db_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepo(t *testing.T) {
	repo := &model.Repository{
		GitHubRepo: model.GitHubRepo{
			Owner:    "tower",
			RepoName: "div",
		},
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
			PkgTypes:     []model.PkgType{model.PkgRubyGems},
			PkgCount:     3,
			VulnCount:    2,
			VulnPkgCount: 1,
		},
	}
	t.Run("Update default branch status", func(t *testing.T) {
		client := newTestTable(t)

		inserted, err := client.InsertRepo(repo)
		require.NoError(t, err)
		assert.True(t, inserted)
		require.NoError(t, client.UpdateBranchIfDefault(&repo.GitHubRepo, branch))

		r1, err := client.FindRepoByFullName(repo.Owner, repo.RepoName)
		require.NoError(t, err)
		assert.Equal(t, r1.Branch, *branch)
	})

	t.Run("Update default branch status if LastScannedAt is greater", func(t *testing.T) {
		client := newTestTable(t)

		inserted, err := client.InsertRepo(repo)
		require.NoError(t, err)
		assert.True(t, inserted)
		require.NoError(t, client.UpdateBranchIfDefault(&repo.GitHubRepo, branch))

		b2 := &model.Branch{
			GitHubBranch: model.GitHubBranch{
				Branch: "blue",
			},
			LastScannedAt: 2345,
			ReportSummary: model.ScanReportSummary{
				PkgTypes:     []model.PkgType{model.PkgRubyGems},
				PkgCount:     4,
				VulnCount:    5,
				VulnPkgCount: 6,
			},
		}

		require.NoError(t, client.UpdateBranchIfDefault(&repo.GitHubRepo, b2))

		r1, err := client.FindRepoByFullName(repo.Owner, repo.RepoName)
		require.NoError(t, err)
		assert.Equal(t, r1.Branch, *b2)
	})

	t.Run("Do not update default branch status if LastScannedAt is lesser", func(t *testing.T) {
		client := newTestTable(t)

		inserted, err := client.InsertRepo(repo)
		require.NoError(t, err)
		assert.True(t, inserted)
		require.NoError(t, client.UpdateBranchIfDefault(&repo.GitHubRepo, branch))

		b2 := &model.Branch{
			GitHubBranch: model.GitHubBranch{
				Branch: "blue",
			},
			LastScannedAt: 1000,
			ReportSummary: model.ScanReportSummary{
				PkgTypes:     []model.PkgType{model.PkgRubyGems},
				PkgCount:     4,
				VulnCount:    5,
				VulnPkgCount: 6,
			},
		}

		// No error, but not updated
		require.NoError(t, client.UpdateBranchIfDefault(&repo.GitHubRepo, b2))

		r1, err := client.FindRepoByFullName(repo.Owner, repo.RepoName)
		require.NoError(t, err)
		assert.Equal(t, r1.Branch, *branch)
	})

	t.Run("Do not update if not default branch", func(t *testing.T) {
		client := newTestTable(t)
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

func TestOwner(t *testing.T) {
	client := newTestTable(t)

	_, err := client.InsertRepo(&model.Repository{
		GitHubRepo: model.GitHubRepo{
			Owner:    "blue",
			RepoName: "five",
		},
	})
	require.NoError(t, err)

	_, err = client.InsertRepo(&model.Repository{
		GitHubRepo: model.GitHubRepo{
			Owner:    "orange",
			RepoName: "doll",
		},
	})
	require.NoError(t, err)

	_, err = client.InsertRepo(&model.Repository{
		GitHubRepo: model.GitHubRepo{
			Owner:    "blue",
			RepoName: "timeless",
		},
	})
	require.NoError(t, err)

	owners, err := client.FindOwners()
	require.NoError(t, err)
	assert.Equal(t, 2, len(owners))
	assert.Contains(t, []string{owners[0].Name, owners[1].Name}, "blue")
	assert.Contains(t, []string{owners[0].Name, owners[1].Name}, "orange")

	repos, err := client.FindRepoByOwner("blue")
	require.NoError(t, err)
	assert.Equal(t, 2, len(repos))
	assert.Contains(t, []string{repos[0].RepoName, repos[1].RepoName}, "five")
	assert.Contains(t, []string{repos[0].RepoName, repos[1].RepoName}, "timeless")
}
