package usecase_test

import (
	"database/sql"
	"encoding/json"
	"os"
	"testing"
	"time"

	ptypes "github.com/aquasecurity/trivy-db/pkg/types"
	ftypes "github.com/aquasecurity/trivy/pkg/fanal/types"
	ttypes "github.com/aquasecurity/trivy/pkg/types"
	"github.com/google/uuid"
	"github.com/m-mizutani/gots/ptr"
	"github.com/m-mizutani/gots/rands"
	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/usecase"

	_ "github.com/lib/pq"
)

func TestCalcPackageID(t *testing.T) {
	hv1 := usecase.CalcPackageID("go", "pkgA", "v1.0.0")
	hv2 := usecase.CalcPackageID("go", "pkgA", "v1.0.1")
	hv3 := usecase.CalcPackageID("bundler", "pkgA", "v1.0.1")
	hv4 := usecase.CalcPackageID("go", "pkgB", "v1.0.0")

	gt.A(t, []string{hv1, hv2, hv3, hv4}).Distinct()
}

func newTestDB(t *testing.T) *sql.DB {
	testDSN, ok := os.LookupEnv("TEST_DB_DSN")
	if !ok {
		t.Errorf("TEST_DB_DSN is not set")
		t.FailNow()
	}
	dbClient := gt.R1(sql.Open("postgres", testDSN)).NoError(t)
	t.Cleanup(func() {
		gt.NoError(t, dbClient.Close())
	})

	return dbClient
}

func TestSaveScan(t *testing.T) {
	salt := rands.AlphaNum(10)
	report := ttypes.Report{
		ArtifactName: "github.com/m-mizutani/octovy",
		ArtifactType: ttypes.ClassLangPkg,
		Results: ttypes.Results{
			{
				Target: "Gemfile.lock",
				Class:  ttypes.ClassLangPkg,
				Type:   "bundler",
				Packages: []ftypes.Package{
					{
						Name:    "octokit_" + salt,
						Version: "4.18.0",
					},
				},
				Vulnerabilities: []ttypes.DetectedVulnerability{
					{
						VulnerabilityID:  "CVE-2020-1234-" + salt,
						PkgName:          "octokit_" + salt,
						InstalledVersion: "4.18.0",
						FixedVersion:     "4.18.1",
						Vulnerability: ptypes.Vulnerability{
							Title:       "CVE-2020-1234",
							Description: "test",
							Severity:    "HIGH",
							References:  []string{"https://example.com"},
							CVSS: ptypes.VendorCVSS{
								"nvd": {
									V2Vector: "AV:L/AC:M/Au:N/C:C/I:C/A:C",
									V3Vector: "CVSS:3.1/AV:L/AC:H/PR:H/UI:N/S:U/C:H/I:H/A:H",
									V2Score:  6.9,
									V3Score:  6.4,
								},
								"redhat": {
									V3Vector: "CVSS:3.1/AV:L/AC:H/PR:H/UI:N/S:U/C:H/I:H/A:H",
									V3Score:  6.4,
								},
							},
							CweIDs:           []string{"CWE-1234"},
							PublishedDate:    ptr.To(time.Now()),
							LastModifiedDate: ptr.To(time.Now()),
						},
					},
				},
			},
		},
	}

	dbClient := newTestDB(t)
	ctx := model.NewContext()
	meta := &usecase.GitHubRepoMetadata{
		GitHubCommit: usecase.GitHubCommit{
			GitHubRepo: usecase.GitHubRepo{
				Owner: "m-mizutani",
				Repo:  "octovy",
			},
			CommitID: "1234567890",
		},
	}
	gt.NoError(t, usecase.SaveScanReportGitHubRepo(ctx, dbClient, &report, meta))
}

func TestUpsertVulnerability(t *testing.T) {
	salt := rands.AlphaNum(10)

	now := time.Now()

	report := ttypes.Report{
		ArtifactName: "github.com/m-mizutani/octovy",
		ArtifactType: ttypes.ClassLangPkg,
		Results: ttypes.Results{
			{
				Target: "Gemfile.lock",
				Class:  ttypes.ClassLangPkg,
				Type:   "bundler",
				Packages: []ftypes.Package{
					{
						Name:    "octokit_" + salt,
						Version: "4.18.0",
					},
				},
				Vulnerabilities: []ttypes.DetectedVulnerability{
					{
						VulnerabilityID:  "CVE-2020-1234-" + salt,
						PkgName:          "octokit_" + salt,
						InstalledVersion: "4.18.0",
						Vulnerability: ptypes.Vulnerability{
							Title:            "CVE-2020-1234",
							Description:      "test",
							Severity:         "MIDDLE",
							References:       []string{"https://example.com"},
							PublishedDate:    ptr.To(now),
							LastModifiedDate: ptr.To(now),
						},
					},
				},
			},
		},
	}

	dbClient := newTestDB(t)
	q := db.New(dbClient)
	ctx := model.NewContext()
	meta := &usecase.GitHubRepoMetadata{
		GitHubCommit: usecase.GitHubCommit{
			GitHubRepo: usecase.GitHubRepo{
				Owner: "m-mizutani",
				Repo:  "octovy",
			},
			CommitID: uuid.NewString(),
		},
	}

	t.Run("save vulnerability", func(t *testing.T) {
		gt.NoError(t, usecase.SaveScanReportGitHubRepo(ctx, dbClient, &report, meta))
		vuln := gt.R1(q.GetVulnerability(ctx, "CVE-2020-1234-"+salt)).NoError(t)
		gt.V(t, vuln).NotNil()
		gt.V(t, vuln.Severity).Equal("MIDDLE")
	})

	t.Run("update severity, but not update last modified", func(t *testing.T) {
		report.Results[0].Vulnerabilities[0].Vulnerability.Severity = "HIGH"
		gt.NoError(t, usecase.SaveScanReportGitHubRepo(ctx, dbClient, &report, meta))
		vuln := gt.R1(q.GetVulnerability(ctx, "CVE-2020-1234-"+salt)).NoError(t)
		gt.V(t, vuln).NotNil()
		gt.V(t, vuln.Severity).Equal("MIDDLE")
	})

	t.Run("update severity and last modified", func(t *testing.T) {
		report.Results[0].Vulnerabilities[0].Vulnerability.Severity = "CRITICAL"
		report.Results[0].Vulnerabilities[0].Vulnerability.LastModifiedDate = ptr.To(now.Add(time.Second))
		gt.NoError(t, usecase.SaveScanReportGitHubRepo(ctx, dbClient, &report, meta))
		resp := gt.R1(q.GetVulnerability(ctx, "CVE-2020-1234-"+salt)).NoError(t)
		gt.V(t, resp).NotNil()
		gt.V(t, resp.Severity).Equal("CRITICAL")

		var vuln ttypes.DetectedVulnerability
		gt.NoError(t, json.Unmarshal(resp.Data.RawMessage, &vuln))
		gt.V(t, vuln.LastModifiedDate.Unix()).Equal(now.Add(time.Second).Unix())
		gt.V(t, vuln.Severity).Equal("CRITICAL")
	})
}
