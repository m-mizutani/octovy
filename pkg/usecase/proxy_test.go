package usecase_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRepositories(t *testing.T) {
	uc, mock := setupUsecase(t)
	injectGitHubMock(t, mock)
	ctx := context.Background()
	branch := "main"
	var calledScan int
	mock.Trivy.ScanMock = func(dir string) (*model.TrivyReport, error) {
		calledScan++
		assert.FileExists(t, filepath.Join(dir, "Gemfile"))
		assert.FileExists(t, filepath.Join(dir, "Gemfile.lock"))
		return &model.TrivyReport{
			Results: model.TrivyResults{
				{
					Target: "Gemfile",
				},
				{
					Target: "Gemfile.lock",
				},
			},
		}, nil
	}

	uc.SendScanRequest(&model.ScanRepositoryRequest{
		InstallID: 1,
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "five",
				},
				Branch: branch,
			},
			CommitID:    "1234567",
			UpdatedAt:   2000,
			RequestedAt: 2100,
		},
	})
	usecase.CloseScanQueue(uc)

	require.NoError(t, uc.Init())

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
	require.NoError(t, usecase.RunScanThread(uc))

	t.Run("got result after scan", func(t *testing.T) {
		resp, err := uc.GetRepositories(ctx)
		require.NoError(t, err)
		require.Len(t, resp, 1)
		assert.Equal(t, "blue", resp[0].Owner)
		assert.Equal(t, "five", resp[0].Name)
		require.NotNil(t, resp[0].Edges.Latest, 1)
		assert.Equal(t, "1234567", resp[0].Edges.Latest.CommitID)
	})
	assert.Equal(t, calledScan, 1)
}

func TestGetVulnerability(t *testing.T) {
	uc, mock := setupUsecase(t)
	injectGitHubMock(t, mock)
	ctx := context.Background()
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
				Results: model.TrivyResults{
					{
						Target: "Gemfile.lock",
						Packages: []model.TrivyPackage{
							{
								Name:    "orange",
								Version: "0.0.1",
							},
						},
						Vulnerabilities: []model.DetectedVulnerability{
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
				Results: model.TrivyResults{
					{
						Target: "Gemfile.lock",
						Packages: []model.TrivyPackage{
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
				Results: model.TrivyResults{
					{
						Target: "Gemfile.lock",
						Packages: []model.TrivyPackage{
							{
								Name:    "orange",
								Version: "0.0.1",
							},
						},
						Vulnerabilities: []model.DetectedVulnerability{
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

	uc.SendScanRequest(&model.ScanRepositoryRequest{
		InstallID: 1,
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "five",
				},
				Branch: branch,
			},
			CommitID: "1234567",
		},
	})
	uc.SendScanRequest(&model.ScanRepositoryRequest{
		InstallID: 1,
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "timeless",
				},
				Branch: branch,
			},
			CommitID: "1234567",
		},
	})
	uc.SendScanRequest(&model.ScanRepositoryRequest{
		InstallID: 1,
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "words",
				},
				Branch: branch,
			},
			CommitID: "1234567",
		},
	})
	usecase.CloseScanQueue(uc)

	require.NoError(t, uc.Init())
	require.NoError(t, usecase.RunScanThread(uc))

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