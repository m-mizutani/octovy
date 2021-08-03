package db_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
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
		results := []*model.ScanReport{
			{
				ReportID: "aaaa",
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
								Type:            model.PkgRubyGems,
								Name:            "hoge",
								Version:         "1.2.3",
								Vulnerabilities: []string{},
							},
						},
					},
				},
				TrivyDBMeta: trivyMeta,
			},
			{
				ReportID: "bbbb",
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
								Type:            model.PkgRubyGems,
								Name:            "hoge",
								Version:         "bbbb",
								Vulnerabilities: []string{},
							},
						},
					},
				},
				TrivyDBMeta: trivyMeta,
			},
			{
				ReportID: "cccc",
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
								Type:            model.PkgRubyGems,
								Name:            "hoge",
								Version:         "1.2.5",
								Vulnerabilities: []string{},
							},
						},
					},
				},
				TrivyDBMeta: trivyMeta,
			},
		}

		require.NoError(t, client.InsertScanReport(results[0]))
		require.NoError(t, client.InsertScanReport(results[1]))
		require.NoError(t, client.InsertScanReport(results[2]))

		t.Run("Lookup report", func(t *testing.T) {
			r, err := client.LookupScanReport("cccc")
			require.NoError(t, err)
			assert.Equal(t, r, results[2])
		})

		t.Run("List latest scan results", func(t *testing.T) {
			r, err := client.FindScanLogsByBranch(&model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "five",
				},
				Branch: "dev",
			}, 2)
			require.NoError(t, err)
			require.Equal(t, 2, len(r))
			assert.Equal(t, "aaaa", r[0].Summary.ReportID)
			assert.Equal(t, "cccc", r[1].Summary.ReportID)
		})

		t.Run("List latest scan results (over)", func(t *testing.T) {
			r, err := client.FindScanLogsByBranch(&model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "five",
				},
				Branch: "dev",
			}, 5)
			require.NoError(t, err)
			require.Equal(t, 3, len(r))
			assert.Equal(t, "aaaa", r[0].Summary.ReportID)
			assert.Equal(t, "cccc", r[1].Summary.ReportID)
			assert.Equal(t, "bbbb", r[2].Summary.ReportID)
		})

		t.Run("No error by find not existing repo/branch", func(t *testing.T) {
			r1, err := client.FindScanLogsByBranch(&model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "five",
				},
				Branch: "end",
			}, 5)
			require.NoError(t, err)
			assert.Zero(t, len(r1))

			r2, err := client.FindScanLogsByBranch(&model.GitHubBranch{
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
			r, err := client.FindScanLogsByCommit(&model.GitHubCommit{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "five",
				},
				CommitID: "beef1111",
			}, 3)

			require.NoError(t, err)
			require.Equal(t, 2, len(r))
			assert.Contains(t, []string{r[0].Summary.ReportID, r[1].Summary.ReportID}, "aaaa")
			assert.Contains(t, []string{r[0].Summary.ReportID, r[1].Summary.ReportID}, "bbbb")
		})
	})
}
