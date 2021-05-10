package db_test

import (
	"testing"

	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBranch(t *testing.T) {
	repo := model.GitHubRepo{
		Owner:    "five",
		RepoName: "timeless",
	}

	t.Run("Insert and lookup", func(t *testing.T) {
		client := newTestTable(t)

		b1 := &model.Branch{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: repo,
				Branch:     "blue",
			},
			LastScannedAt: 1234,
			ReportSummary: model.ScanReportSummary{
				ReportID: "aaaa",
				PkgTypes: []model.PkgType{model.PkgBundler},
			},
		}
		b2 := &model.Branch{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: repo,
				Branch:     "orange",
			},
			LastScannedAt: 2345,
			ReportSummary: model.ScanReportSummary{
				ReportID: "bbbb",
				PkgTypes: []model.PkgType{model.PkgBundler},
			},
		}
		require.NoError(t, client.UpdateBranch(b1))
		require.NoError(t, client.UpdateBranch(b2))

		t.Run("Lookup by onwer/name/branch", func(t *testing.T) {
			r1, err := client.LookupBranch(&model.GitHubBranch{
				GitHubRepo: repo,
				Branch:     "blue",
			})
			require.NoError(t, err)
			require.NotNil(t, r1)
			assert.Equal(t, r1, b1)

			r2, err := client.LookupBranch(&model.GitHubBranch{
				GitHubRepo: repo,
				Branch:     "orange",
			})
			require.NoError(t, err)
			require.NotNil(t, r2)
			assert.Equal(t, r2, b2)
		})

		t.Run("Lookup latest", func(t *testing.T) {
			r1, err := client.FindLatestScannedBranch(&repo, 1)
			require.NoError(t, err)
			require.Equal(t, 1, len(r1))
			assert.Equal(t, b2, r1[0])

			r2, err := client.FindLatestScannedBranch(&repo, 2)
			require.NoError(t, err)
			require.Equal(t, 2, len(r2))
			assert.Equal(t, b2, r2[0])
			assert.Equal(t, b1, r2[1])
		})
	})

	t.Run("Update", func(t *testing.T) {
		base := &model.Branch{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: repo,
				Branch:     "blue",
			},
			LastScannedAt: 1234,
			ReportSummary: model.ScanReportSummary{
				ReportID: "aaaa",
				PkgTypes: []model.PkgType{model.PkgBundler},
			},
		}

		t.Run("Update with greater timestamp", func(t *testing.T) {
			client := newTestTable(t)
			require.NoError(t, client.UpdateBranch(base))

			updated := &model.Branch{
				GitHubBranch: model.GitHubBranch{
					GitHubRepo: repo,
					Branch:     "blue",
				},
				LastScannedAt: 2000,
				ReportSummary: model.ScanReportSummary{
					ReportID: "bbbb",
					PkgTypes: []model.PkgType{model.PkgBundler},
				},
			}

			require.NoError(t, client.UpdateBranch(updated))
			r1, err := client.LookupBranch(&model.GitHubBranch{
				GitHubRepo: repo,
				Branch:     "blue",
			})
			require.NoError(t, err)
			assert.Equal(t, updated, r1)
		})

		t.Run("Do not update with lesser timestamp", func(t *testing.T) {
			client := newTestTable(t)
			require.NoError(t, client.UpdateBranch(base))

			updated := &model.Branch{
				GitHubBranch: model.GitHubBranch{
					GitHubRepo: repo,
					Branch:     "blue",
				},
				LastScannedAt: 1000,
				ReportSummary: model.ScanReportSummary{
					ReportID: "bbbb",
					PkgTypes: []model.PkgType{model.PkgBundler},
				},
			}

			require.NoError(t, client.UpdateBranch(updated))
			r1, err := client.LookupBranch(&model.GitHubBranch{
				GitHubRepo: repo,
				Branch:     "blue",
			})
			require.NoError(t, err)
			assert.Equal(t, base, r1)
		})
	})
}
