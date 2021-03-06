package usecase_test

import (
	"context"
	"encoding/json"
	"testing"

	ftypes "github.com/aquasecurity/fanal/types"
	dtypes "github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aquasecurity/trivy/pkg/report"
	ttypes "github.com/aquasecurity/trivy/pkg/types"

	"github.com/google/go-github/v39/github"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/opa"
	"github.com/m-mizutani/octovy/pkg/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanProcedure(t *testing.T) {
	uc, mock := setupUsecase(t,
		optDBMock(),
		optTrivy(),
		optGitHubMock(),
		optGitHubAppMock(),
		optGitHubAppMockZip(),
	)

	var calledScan int
	mock.Trivy.ScanMock = func(dir string) (*model.TrivyReport, error) {
		calledScan++
		return &model.TrivyReport{
			Results: report.Results{
				{
					Target: "Gemfile.lock",
					Packages: []ftypes.Package{
						{
							Name:    "example",
							Version: "6.1.4",
						},
					},
					Vulnerabilities: []ttypes.DetectedVulnerability{
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

	assert.NoError(t, uc.Scan(model.NewContext(), &model.ScanRepositoryRequest{
		InstallID: 1,
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner: "blue",
					Name:  "five",
				},
				Branch: "main",
			},
			CommitID:    "1234567",
			UpdatedAt:   2000,
			RequestedAt: 2100,
		},
	}))

	assert.Equal(t, 1, calledScan)

	ctx := model.NewContext()
	scan, err := mock.DB.GetLatestScan(ctx, model.GitHubBranch{
		GitHubRepo: model.GitHubRepo{
			Owner: "blue",
			Name:  "five",
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

func TestScanProcedureWithRule(t *testing.T) {
	setup := func(t *testing.T, rule string, update func(repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error) *usecase.Usecase {
		uc, mock := setupUsecase(t,
			optDBMock(),
			optTrivy(),
			optGitHubMock(),
			optGitHubAppMock(),
			optGitHubAppMockZip(),
			optCheckRule(rule, update),
		)

		var calledScan int
		mock.Trivy.ScanMock = func(dir string) (*model.TrivyReport, error) {
			calledScan++
			return &model.TrivyReport{
				Results: report.Results{
					{
						Target: "Gemfile.lock",
						Packages: []ftypes.Package{
							{
								Name:    "example",
								Version: "6.1.4",
							},
						},
						Vulnerabilities: []ttypes.DetectedVulnerability{
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

		t.Cleanup(func() {
			assert.Equal(t, 1, calledScan)
		})

		return uc
	}

	scanReq := &model.ScanRepositoryRequest{
		InstallID: 1,
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner: "blue",
					Name:  "five",
				},
				Branch: "main",
			},
			CommitID:    "1234567",
			UpdatedAt:   2000,
			RequestedAt: 2100,
		},
	}

	testCases := []struct {
		title      string
		called     int
		rule       string
		conclusion string
	}{
		{
			title:  "always success",
			called: 1,
			rule: `package octovy.check
			conclusion = "success"`,
			conclusion: "success",
		},
		{
			title:  "always failure",
			called: 1,
			rule: `package octovy.check
			conclusion = "failure"`,
			conclusion: "failure",
		},
		{
			title:  "failure if vulnID has CVE-1000",
			called: 1,
			rule: `package octovy.check
			default conclusion = "success"
			conclusion = "failure" {
				vulnID := input.sources[_].packages[_].vuln_ids[_]
				vulnID == "CVE-1000"
			}
			`,
			conclusion: "failure",
		},
		{
			title:  "failure if vulnID has CVE-1001, then success",
			called: 1,
			rule: `package octovy.check
			default conclusion = "success"
			conclusion = "failure" {
				vulnID := input.sources[_].packages[_].vuln_ids[_]
				vulnID == "CVE-1001"
			}
			`,
			conclusion: "success",
		},
	}

	for _, c := range testCases {
		t.Run(c.title, func(t *testing.T) {
			var called int
			uc := setup(t, c.rule, func(repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error {
				called++
				assert.Equal(t, c.conclusion, *opt.Conclusion)
				return nil
			})

			assert.NoError(t, uc.Scan(model.NewContext(), scanReq))
			assert.Equal(t, 1, called)
		})
	}
}

func TestScanProcedureWithOPA(t *testing.T) {
	uc, mock := setupUsecase(t,
		optDBMock(),
		optTrivy(),
		optGitHubMock(),
		optGitHubAppMock(),
		optGitHubAppMockZip(),
		optOPAServer(),
	)

	var calledScan int
	mock.Trivy.ScanMock = func(dir string) (*model.TrivyReport, error) {
		calledScan++
		return &model.TrivyReport{
			Results: report.Results{
				{
					Target: "Gemfile.lock",
					Packages: []ftypes.Package{
						{
							Name:    "example",
							Version: "6.1.4",
						},
					},
					Vulnerabilities: []ttypes.DetectedVulnerability{
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

	var calledOPA int
	mock.OPA.MockData = func(ctx context.Context, pkg opa.RegoPkg, input, result interface{}) error {
		var report model.ScanReport
		raw, err := json.Marshal(input)
		require.NoError(t, err)
		json.Unmarshal(raw, &report)

		assert.Equal(t, "blue", report.Repo.Owner)
		assert.Equal(t, "five", report.Repo.Name)
		assert.Equal(t, opa.Check, pkg)
		calledOPA++
		return nil
	}

	var callCreateCheck, callUpdateCheck int
	mock.GtiHubApp.CreateCheckRunMock = func(repo *model.GitHubRepo, commit string) (int64, error) {
		callCreateCheck++
		return 0, nil
	}
	mock.GtiHubApp.UpdateCheckRunMock = func(repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error {
		callUpdateCheck++
		return nil
	}

	scanReq := &model.ScanRepositoryRequest{
		InstallID: 1,
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner: "blue",
					Name:  "five",
				},
				Branch: "main",
			},
			CommitID:    "1234567",
			UpdatedAt:   2000,
			RequestedAt: 2100,
		},
	}

	assert.NoError(t, uc.Scan(model.NewContext(), scanReq))
	assert.Equal(t, 1, calledScan)
	assert.Equal(t, 1, calledOPA)
	assert.Equal(t, 1, callCreateCheck)
	assert.Equal(t, 1, callUpdateCheck)
}
