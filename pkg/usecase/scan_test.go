package usecase_test

import (
	"testing"

	dtypes "github.com/aquasecurity/trivy-db/pkg/types"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanProcedure(t *testing.T) {
	uc, mock := setupUsecase(t)
	injectGitHubMock(t, mock, false)
	var calledScan int
	mock.Trivy.ScanMock = func(dir string) (*model.TrivyReport, error) {
		calledScan++
		return &model.TrivyReport{
			Results: model.TrivyResults{
				{
					Target: "Gemfile.lock",
					Packages: []model.TrivyPackage{
						{
							Name:    "example",
							Version: "6.1.4",
						},
					},
					Vulnerabilities: []model.DetectedVulnerability{
						{
							VulnerabilityID:  "CVE-1000",
							PkgName:          "example",
							InstalledVersion: "6.1.4",
							FixedVersion:     "6.1.5",
							Vulnerability: dtypes.Vulnerability{
								Title:       "test vuln",
								Description: "it's test",
								Severity:    "low",
							},
						},
					},
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
				Branch: "main",
			},
			CommitID:    "1234567",
			UpdatedAt:   2000,
			RequestedAt: 2100,
		},
	})
	usecase.CloseScanQueue(uc)

	require.NoError(t, uc.Init())
	require.NoError(t, usecase.RunScanThread(uc))

	assert.Equal(t, 1, calledScan)

	ctx := model.NewContext()
	scan, err := mock.DB.GetLatestScan(ctx, model.GitHubBranch{
		GitHubRepo: model.GitHubRepo{
			Owner:    "blue",
			RepoName: "five",
		},
		Branch: "main",
	})
	require.NoError(t, err)
	assert.Equal(t, "1234567", scan.CommitID)
	require.Len(t, scan.Edges.Packages, 1)
	assert.Equal(t, "example", scan.Edges.Packages[0].Name)
	require.Len(t, scan.Edges.Packages[0].Edges.Vulnerabilities, 1)
	assert.Equal(t, "test vuln", scan.Edges.Packages[0].Edges.Vulnerabilities[0].Title)
}
