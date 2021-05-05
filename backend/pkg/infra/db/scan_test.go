package db_test

import (
	"testing"

	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanResult(t *testing.T) {
	t.Run("Insert and find vulnerabilities", func(t *testing.T) {
		client := newTestTable(t)
		trivyMeta := model.TrivyDBMeta{
			Version:   1,
			Type:      1,
			UpdatedAt: 2345,
		}
		results := []*model.ScanResult{
			{
				Target: model.ScanTarget{
					GitHubBranch: model.GitHubBranch{
						GitHubRepo: model.GitHubRepo{
							Owner:    "blue",
							RepoName: "five",
						},
						Branch: "dev",
					},
					CommitID:  "beef1111",
					UpdatedAt: 1230,
				},
				ScannedAt: 3000,
				Sources: []*model.PackageSource{
					{
						Source: "Gemfile.lock",
						Packages: []*model.Package{
							{
								Type:    model.PkgBundler,
								Name:    "hoge",
								Version: "1.2.3",
							},
						},
					},
				},
				TrivyDBMeta: trivyMeta,
			},
			{
				Target: model.ScanTarget{
					GitHubBranch: model.GitHubBranch{
						GitHubRepo: model.GitHubRepo{
							Owner:    "blue",
							RepoName: "five",
						},
						Branch: "dev",
					},
					CommitID:  "beef1111",
					UpdatedAt: 1230,
				},
				ScannedAt: 1000,
				Sources: []*model.PackageSource{
					{
						Source: "Gemfile.lock",
						Packages: []*model.Package{
							{
								Type:    model.PkgBundler,
								Name:    "hoge",
								Version: "1.2.4",
							},
						},
					},
				},
				TrivyDBMeta: trivyMeta,
			},
			{
				Target: model.ScanTarget{
					GitHubBranch: model.GitHubBranch{
						GitHubRepo: model.GitHubRepo{
							Owner:    "blue",
							RepoName: "five",
						},
						Branch: "dev",
					},
					CommitID:  "beef2222",
					UpdatedAt: 1240,
				},
				ScannedAt: 2000,
				Sources: []*model.PackageSource{
					{
						Source: "Gemfile.lock",
						Packages: []*model.Package{
							{
								Type:    model.PkgBundler,
								Name:    "hoge",
								Version: "1.2.5",
							},
						},
					},
				},
				TrivyDBMeta: trivyMeta,
			},
		}

		require.NoError(t, client.InsertScanResult(results[0]))
		require.NoError(t, client.InsertScanResult(results[1]))
		require.NoError(t, client.InsertScanResult(results[2]))

		t.Run("List latest scan results", func(t *testing.T) {
			r, err := client.FindLatestScanResults(&model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "five",
				},
				Branch: "dev",
			}, 2)
			require.NoError(t, err)
			require.Equal(t, 2, len(r))
			assert.Equal(t, "1.2.3", r[0].Sources[0].Packages[0].Version)
			assert.Equal(t, "1.2.5", r[1].Sources[0].Packages[0].Version)
		})

		t.Run("List latest scan results (over)", func(t *testing.T) {
			r, err := client.FindLatestScanResults(&model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "five",
				},
				Branch: "dev",
			}, 5)
			require.NoError(t, err)
			require.Equal(t, 3, len(r))
			assert.Equal(t, "1.2.3", r[0].Sources[0].Packages[0].Version)
			assert.Equal(t, "1.2.5", r[1].Sources[0].Packages[0].Version)
			assert.Equal(t, "1.2.4", r[2].Sources[0].Packages[0].Version)
		})

		t.Run("No error by find not existing repo/branch", func(t *testing.T) {
			r1, err := client.FindLatestScanResults(&model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "five",
				},
				Branch: "end",
			}, 5)
			require.NoError(t, err)
			assert.Zero(t, len(r1))

			r2, err := client.FindLatestScanResults(&model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "six",
				},
				Branch: "dev",
			}, 5)
			require.NoError(t, err)
			assert.Zero(t, len(r2))
		})

		t.Run("Find latest result of commitID", func(t *testing.T) {
			r, err := client.FindScanResult(&model.GitHubCommit{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "five",
				},
				CommitID: "beef1111",
			})

			require.NoError(t, err)
			require.NotNil(t, r)
			assert.Equal(t, "1.2.3", r.Sources[0].Packages[0].Version)
		})
	})
}
