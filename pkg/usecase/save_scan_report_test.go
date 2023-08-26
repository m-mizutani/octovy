package usecase_test

import (
	"database/sql"
	"os"
	"testing"
	"time"

	ptypes "github.com/aquasecurity/trivy-db/pkg/types"
	ftypes "github.com/aquasecurity/trivy/pkg/fanal/types"
	ttypes "github.com/aquasecurity/trivy/pkg/types"
	"github.com/m-mizutani/gots/ptr"
	"github.com/m-mizutani/gots/rands"
	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/model"
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
		t.Skip("TEST_DB_DSN is not set")
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

	gt.NoError(t, usecase.SaveScanReport(ctx, dbClient, &report))
}
