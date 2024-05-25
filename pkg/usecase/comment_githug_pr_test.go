package usecase_test

import (
	"os"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/model/trivy"
	"github.com/m-mizutani/octovy/pkg/usecase"
)

func TestRenderScanReport(t *testing.T) {
	report := trivy.Report{
		Results: []trivy.Result{
			{
				Target: "target1",
				Vulnerabilities: []trivy.DetectedVulnerability{
					{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1", Vulnerability: trivy.Vulnerability{Title: "Vuln title1", Severity: "HIGH"}},
					{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2", Vulnerability: trivy.Vulnerability{Title: "Vuln title2", Severity: "CRITICAL"}},
				},
			},
			{
				Target: "target2",
				Vulnerabilities: []trivy.DetectedVulnerability{
					{VulnerabilityID: "CVE-0000-0003", PkgName: "pkg4", Vulnerability: trivy.Vulnerability{Title: "Vuln title3", Severity: "CRITICAL"}},
				},
			},
		},
	}
	added := trivy.Results{
		{
			Target: "target1",
			Vulnerabilities: []trivy.DetectedVulnerability{
				{
					VulnerabilityID: "CVE-0000-0002",
					PkgName:         "pkg2",
					Vulnerability: trivy.Vulnerability{
						Title: "Vuln title2",
					},
				},
			},
		},
	}
	fixed := trivy.Results{
		{
			Target: "target2",
			Vulnerabilities: []trivy.DetectedVulnerability{
				{
					VulnerabilityID: "CVE-0000-0003",
					PkgName:         "pkg3",
					Vulnerability: trivy.Vulnerability{
						Title: "Vuln title3",
					},
				},
			},
		},
	}

	body, err := usecase.RenderScanReport(&report, added, fixed)
	gt.NoError(t, err)
	gt.NoError(t, os.WriteFile("templates/test_comment_body.md", []byte(body), 0644))
}
