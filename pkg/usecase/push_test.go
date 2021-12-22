package usecase_test

import (
	"testing"

	ftypes "github.com/aquasecurity/fanal/types"
	ttypes "github.com/aquasecurity/trivy/pkg/types"

	"github.com/aquasecurity/trivy/pkg/report"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPushTrivyResult(t *testing.T) {
	uc, mock := setupUsecase(t, optDBMock())
	ctx := model.NewContext()
	req := &model.PushTrivyResultRequest{
		Target: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner: "blue",
					Name:  "five",
				},
				Branch: "master",
			},
			CommitID: "abcde12345abcde12345abcde12345abcde12345",
		},
		Report: model.TrivyReport{
			Results: report.Results{
				{
					Class:  report.ClassOSPkg,
					Target: "gcr.io/xxx",
					Type:   "debian",
					Packages: []ftypes.Package{
						{
							Name:    "libx",
							Version: "0.0.1",
						},
					},
					Vulnerabilities: []ttypes.DetectedVulnerability{
						{
							VulnerabilityID:  "CVE-0000",
							PkgName:          "libx",
							InstalledVersion: "0.0.1",
						},
					},
				},
			},
		},
	}
	require.NoError(t, uc.PushTrivyResult(ctx, req))

	scan, err := mock.DB.GetLatestScan(ctx, model.GitHubBranch{
		GitHubRepo: model.GitHubRepo{
			Owner: "blue",
			Name:  "five",
		},
		Branch: "master",
	})
	require.NoError(t, err)
	require.NotNil(t, scan)
	assert.Equal(t, "abcde12345abcde12345abcde12345abcde12345", scan.CommitID)
	require.Len(t, scan.Edges.Packages, 1)
	assert.Equal(t, "os-pkgs@debian", scan.Edges.Packages[0].Source)
	assert.Equal(t, "libx", scan.Edges.Packages[0].Name)
}
