package model_test

import (
	"testing"
	"time"

	ftypes "github.com/aquasecurity/fanal/types"
	dtypes "github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aquasecurity/trivy/pkg/report"
	ttypes "github.com/aquasecurity/trivy/pkg/types"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrivyReportToEnt(t *testing.T) {
	t.Run("empty to empty", func(t *testing.T) {
		pkg, vuln := model.TrivyReportToEnt(&model.TrivyReport{}, time.Now())
		assert.Len(t, pkg, 0)
		assert.Len(t, vuln, 0)
	})

	t.Run("fully check", func(t *testing.T) {
		ts1 := time.Now()
		ts2 := ts1.Add(-time.Hour * 30)

		pkg, vuln := model.TrivyReportToEnt(&model.TrivyReport{
			Results: report.Results{
				{
					Target: "Gemfile.lock",
					Class:  "lang-pkgs",
					Type:   "bundler",
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
								CweIDs:      []string{"CWE-000"},
								CVSS: dtypes.VendorCVSS{
									"x": dtypes.CVSS{
										V2Vector: "test2",
										V3Vector: "test3",
									},
								},
								References: []string{
									"https://example.com",
								},
								LastModifiedDate: &ts2,
							},
						},
					},
				},
			},
		}, ts1)

		require.Len(t, pkg, 1)
		assert.Equal(t, &ent.PackageRecord{
			Type:    "bundler",
			Name:    "example",
			Source:  "Gemfile.lock",
			Version: "6.1.4",
			VulnIds: []string{"CVE-1000"},
		}, pkg[0])

		require.Len(t, vuln, 1)
		assert.Equal(t, &ent.Vulnerability{
			ID:             "CVE-1000",
			FirstSeenAt:    ts1.Unix(),
			LastModifiedAt: ts2.Unix(),
			Title:          "test vuln",
			Description:    "it's test",
			Severity:       "low",
			CweID:          []string{"CWE-000"},
			Cvss: []string{
				"x,V2Vector,test2",
				"x,V3Vector,test3",
			},
			References: []string{
				"https://example.com",
			},
		}, vuln[0])
	})

	t.Run("no matched package (invalid data, but ignore)", func(t *testing.T) {
		t.Run("version is not matched", func(t *testing.T) {
			pkg, vuln := model.TrivyReportToEnt(&model.TrivyReport{
				Results: report.Results{
					{
						Target: "Gemfile.lock",
						Type:   "bundler",
						Packages: []ftypes.Package{
							{
								Name:    "example",
								Version: "1.1.4",
							},
						},
						Vulnerabilities: []ttypes.DetectedVulnerability{
							{
								VulnerabilityID:  "CVE-1000",
								PkgName:          "example",
								InstalledVersion: "6.1.4",
							},
						},
					},
				},
			}, time.Now())

			require.Len(t, pkg, 1)
			require.Len(t, vuln, 1)
			assert.Len(t, pkg[0].VulnIds, 0)
		})

		t.Run("name is not matched", func(t *testing.T) {
			pkg, vuln := model.TrivyReportToEnt(&model.TrivyReport{
				Results: report.Results{
					{
						Target: "Gemfile.lock",
						Type:   "bundler",
						Packages: []ftypes.Package{
							{
								Name:    "blue",
								Version: "6.1.4",
							},
						},
						Vulnerabilities: []ttypes.DetectedVulnerability{
							{
								VulnerabilityID:  "CVE-1000",
								PkgName:          "example",
								InstalledVersion: "6.1.4",
							},
						},
					},
				},
			}, time.Now())

			require.Len(t, pkg, 1)
			require.Len(t, vuln, 1)
			assert.Len(t, pkg[0].VulnIds, 0)
		})
	})

	t.Run("vulnerability not duplicated", func(t *testing.T) {
		pkg, vuln := model.TrivyReportToEnt(&model.TrivyReport{
			Results: report.Results{
				{
					Target: "Gemfile.lock",
					Type:   "bundler",
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
						},
					},
				},
				{
					Target: "tmp/Gemfile.lock",
					Type:   "bundler",
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
						},
					},
				},
			},
		}, time.Now())

		require.Len(t, pkg, 2)
		require.Len(t, vuln, 1)
	})

	t.Run("import OS package", func(t *testing.T) {
		ts1 := time.Now()
		ts2 := ts1.Add(-time.Hour * 30)

		pkg, vuln := model.TrivyReportToEnt(&model.TrivyReport{
			Results: report.Results{
				{
					Target: "gcr.io/example:xxxxxxxx",
					Class:  "os-pkgs",
					Type:   "debian",
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
								CweIDs:      []string{"CWE-000"},
								CVSS: dtypes.VendorCVSS{
									"x": dtypes.CVSS{
										V2Vector: "test2",
										V3Vector: "test3",
									},
								},
								References: []string{
									"https://example.com",
								},
								LastModifiedDate: &ts2,
							},
						},
					},
				},
			},
		}, ts1)

		require.Len(t, pkg, 1)
		require.Len(t, vuln, 1)
		assert.Len(t, pkg[0].VulnIds, 1)
		assert.Equal(t, pkg[0].Source, "os-pkgs@debian")
		assert.Contains(t, pkg[0].VulnIds, "CVE-1000")
	})
}
