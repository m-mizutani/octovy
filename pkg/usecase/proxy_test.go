package usecase_test

import (
	"path/filepath"
	"testing"

	ftypes "github.com/aquasecurity/fanal/types"
	"github.com/aquasecurity/trivy/pkg/report"
	ttypes "github.com/aquasecurity/trivy/pkg/types"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRepositories(t *testing.T) {
	uc, mock := setupUsecase(t,
		optDBMock(),
		optTrivy(),
		optGitHubMock(),
		optGitHubAppMock(),
		optGitHubAppMockZip(),
	)

	ctx := model.NewContext()
	branch := "main"
	var calledScan int
	mock.Trivy.ScanMock = func(dir string) (*model.TrivyReport, error) {
		calledScan++
		assert.FileExists(t, filepath.Join(dir, "Gemfile"))
		assert.FileExists(t, filepath.Join(dir, "Gemfile.lock"))
		return &model.TrivyReport{
			Results: report.Results{
				{
					Target: "Gemfile",
				},
				{
					Target: "Gemfile.lock",
				},
			},
		}, nil
	}

	t.Run("no result before scan", func(t *testing.T) {
		r0, err := uc.GetRepositories(ctx)
		require.NoError(t, err)
		require.Len(t, r0, 0)
	})

	_, err := mock.DB.CreateRepo(ctx, &ent.Repository{
		Owner:         "blue",
		Name:          "five",
		DefaultBranch: &branch,
	})
	require.NoError(t, err)

	assert.NoError(t, uc.Scan(ctx, &model.ScanRepositoryRequest{
		InstallID: 1,
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner: "blue",
					Name:  "five",
				},
				Branch: branch,
			},
			CommitID:    "1234567",
			UpdatedAt:   2000,
			RequestedAt: 2100,
		},
	}))

	{
		resp, err := uc.GetRepositories(ctx)
		require.NoError(t, err)
		require.Len(t, resp, 1)
		assert.Equal(t, "blue", resp[0].Owner)
		assert.Equal(t, "five", resp[0].Name)
		require.NotNil(t, resp[0].Edges.Latest)
		assert.Equal(t, "1234567", resp[0].Edges.Latest.CommitID)
	}
	assert.Equal(t, calledScan, 1)
}

func TestGetVulnerability(t *testing.T) {
	uc, mock := setupUsecase(t,
		optDBMock(),
		optTrivy(),
		optGitHubMock(),
		optGitHubAppMock(),
		optGitHubAppMockZip(),
	)

	ctx := model.NewContext()
	branch := "main"
	var calledScan int

	var err error
	_, err = mock.DB.CreateRepo(ctx, &ent.Repository{
		Owner:         "blue",
		Name:          "five",
		DefaultBranch: &branch,
	})
	require.NoError(t, err)
	_, err = mock.DB.CreateRepo(ctx, &ent.Repository{
		Owner:         "blue",
		Name:          "timeless",
		DefaultBranch: &branch,
	})
	require.NoError(t, err)
	_, err = mock.DB.CreateRepo(ctx, &ent.Repository{
		Owner:         "blue",
		Name:          "words",
		DefaultBranch: &branch,
	})
	require.NoError(t, err)

	mock.Trivy.ScanMock = func(dir string) (*model.TrivyReport, error) {
		calledScan++

		switch calledScan {
		case 1: // has targeted vuln -> blue/five
			return &model.TrivyReport{
				Results: report.Results{
					{
						Target: "Gemfile.lock",
						Packages: []ftypes.Package{
							{
								Name:    "orange",
								Version: "0.0.1",
							},
						},
						Vulnerabilities: []ttypes.DetectedVulnerability{
							{
								VulnerabilityID:  "CVE-0001",
								PkgName:          "orange",
								InstalledVersion: "0.0.1",
							},
						},
					},
				},
			}, nil

		case 2: // not matched -> blue/timeless
			return &model.TrivyReport{
				Results: report.Results{
					{
						Target: "Gemfile.lock",
						Packages: []ftypes.Package{
							{
								Name:    "orange",
								Version: "0.0.1",
							},
						},
					},
				},
			}, nil

		case 3: // matched vuln, but not target -> blue/words
			return &model.TrivyReport{
				Results: report.Results{
					{
						Target: "Gemfile.lock",
						Packages: []ftypes.Package{
							{
								Name:    "orange",
								Version: "0.0.1",
							},
						},
						Vulnerabilities: []ttypes.DetectedVulnerability{
							{
								VulnerabilityID:  "CVE-0002",
								PkgName:          "orange",
								InstalledVersion: "0.0.1",
							},
						},
					},
				},
			}, nil
		}

		return nil, nil
	}

	require.NoError(t, uc.Scan(ctx, &model.ScanRepositoryRequest{
		InstallID: 1,
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner: "blue",
					Name:  "five",
				},
				Branch: branch,
			},
			CommitID: "1234567",
		},
	}))
	require.NoError(t, uc.Scan(ctx, &model.ScanRepositoryRequest{
		InstallID: 1,
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner: "blue",
					Name:  "timeless",
				},
				Branch: branch,
			},
			CommitID: "1234567",
		},
	}))
	require.NoError(t, uc.Scan(ctx, &model.ScanRepositoryRequest{
		InstallID: 1,
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner: "blue",
					Name:  "words",
				},
				Branch: branch,
			},
			CommitID: "1234567",
		},
	}))

	t.Run("3 repository with latset scan found", func(t *testing.T) {
		resp, err := uc.GetRepositories(ctx)
		require.NoError(t, err)
		require.Len(t, resp, 3)
		for _, repo := range resp {
			require.NotNil(t, repo.Edges.Latest)
		}
	})
	assert.Equal(t, calledScan, 3)

	t.Run("only 1 repository with CVE-0001 found", func(t *testing.T) {
		vuln, err := uc.GetVulnerability(ctx, "CVE-0001")
		require.NoError(t, err)
		require.NotNil(t, vuln)
		require.NotNil(t, vuln.Vulnerability)
		require.Len(t, vuln.Affected, 1)
		assert.Equal(t, "five", vuln.Affected[0].Name)
	})
}
