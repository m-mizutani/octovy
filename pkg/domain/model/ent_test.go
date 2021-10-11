package model_test

import (
	"testing"
	"time"

	"github.com/aquasecurity/trivy-db/pkg/types"

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
			Results: model.TrivyResults{
				{
					Target: "Gemfile.lock",
					Type:   "bundler",
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
							Vulnerability: types.Vulnerability{
								Title:       "test vuln",
								Description: "it's test",
								Severity:    "low",
								CweIDs:      []string{"CWE-000"},
								CVSS: types.VendorCVSS{
									"x": types.CVSS{
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
				Results: model.TrivyResults{
					{
						Target: "Gemfile.lock",
						Type:   "bundler",
						Packages: []model.TrivyPackage{
							{
								Name:    "example",
								Version: "1.1.4",
							},
						},
						Vulnerabilities: []model.DetectedVulnerability{
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
				Results: model.TrivyResults{
					{
						Target: "Gemfile.lock",
						Type:   "bundler",
						Packages: []model.TrivyPackage{
							{
								Name:    "blue",
								Version: "6.1.4",
							},
						},
						Vulnerabilities: []model.DetectedVulnerability{
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
			Results: model.TrivyResults{
				{
					Target: "Gemfile.lock",
					Type:   "bundler",
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
						},
					},
				},
				{
					Target: "tmp/Gemfile.lock",
					Type:   "bundler",
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
						},
					},
				},
			},
		}, time.Now())

		require.Len(t, pkg, 2)
		require.Len(t, vuln, 1)
	})
}
