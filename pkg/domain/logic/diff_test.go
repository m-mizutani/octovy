package logic_test

import (
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/logic"
	"github.com/m-mizutani/octovy/pkg/domain/model/trivy"
)

func TestDiffResults(t *testing.T) {
	type testCase struct {
		oldReport, newReport trivy.Report
		fixed, added         trivy.Results
	}

	test := func(c testCase) func(t *testing.T) {
		return func(t *testing.T) {
			fixed, added := logic.DiffResults(&c.oldReport, &c.newReport)
			gt.Equal(t, fixed, c.fixed)
			gt.Equal(t, added, c.added)
		}
	}

	t.Run("No diff", test(testCase{
		oldReport: trivy.Report{
			Results: []trivy.Result{
				{
					Target: "target1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
					},
				},
			},
		},
		newReport: trivy.Report{
			Results: []trivy.Result{
				{
					Target: "target1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
					},
				},
			},
		},
		fixed: nil,
		added: nil,
	}))

	t.Run("Add new vulnerability", test(testCase{
		oldReport: trivy.Report{
			Results: []trivy.Result{
				{
					Target: "target1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
					},
				},
			},
		},
		newReport: trivy.Report{
			Results: []trivy.Result{
				{
					Target: "target1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
						{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2"},
					},
				},
			},
		},
		fixed: nil,
		added: []trivy.Result{
			{
				Target: "target1",
				Vulnerabilities: []trivy.DetectedVulnerability{
					{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2"},
				},
			},
		},
	}))

	t.Run("Fix vulnerability", test(testCase{
		oldReport: trivy.Report{
			Results: []trivy.Result{
				{
					Target: "target1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
						{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2"},
					},
				},
			},
		},
		newReport: trivy.Report{
			Results: []trivy.Result{
				{
					Target: "target1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
					},
				},
			},
		},
		fixed: []trivy.Result{
			{
				Target: "target1",
				Vulnerabilities: []trivy.DetectedVulnerability{
					{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2"},
				},
			},
		},
		added: nil,
	}))

	t.Run("Add and fix vulnerability", test(testCase{
		oldReport: trivy.Report{
			Results: []trivy.Result{
				{
					Target: "target1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
						{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2"},
					},
				},
			},
		},
		newReport: trivy.Report{
			Results: []trivy.Result{
				{
					Target: "target1",
					Vulnerabilities: []trivy.DetectedVulnerability{

						{VulnerabilityID: "CVE-0000-0003", PkgName: "pkg3"},

						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
					},
				},
			},
		},
		fixed: []trivy.Result{
			{
				Target: "target1",
				Vulnerabilities: []trivy.DetectedVulnerability{
					{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2"},
				},
			},
		},
		added: []trivy.Result{
			{
				Target: "target1",
				Vulnerabilities: []trivy.DetectedVulnerability{
					{VulnerabilityID: "CVE-0000-0003", PkgName: "pkg3"},
				},
			},
		},
	}))

	t.Run("No diff with multiple results", test(testCase{
		oldReport: trivy.Report{
			Results: []trivy.Result{
				{
					Target: "target1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
					},
				},
				{
					Target: "target2",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2"},
					},
				},
			},
		},
		newReport: trivy.Report{
			Results: []trivy.Result{
				{
					Target: "target1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
					},
				},
				{
					Target: "target2",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2"},
					},
				},
			},
		},
		fixed: nil,
		added: nil,
	}))

	t.Run("Add new vulnerability with multiple results", test(testCase{
		oldReport: trivy.Report{
			Results: []trivy.Result{
				{
					Target: "target1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
					},
				},
				{
					Target: "target2",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2"},
					},
				},
			},
		},
		newReport: trivy.Report{
			Results: []trivy.Result{
				{
					Target: "target1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
						{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2"},
					},
				},
				{
					Target: "target2",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
						{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2"},
					},
				},
			},
		},
		fixed: nil,
		added: []trivy.Result{
			{
				Target: "target1",
				Vulnerabilities: []trivy.DetectedVulnerability{
					{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2"},
				},
			},
			{
				Target: "target2",
				Vulnerabilities: []trivy.DetectedVulnerability{
					{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
				},
			},
		},
	}))

	t.Run("Fix vulnerability with multiple results", test(testCase{
		oldReport: trivy.Report{
			Results: []trivy.Result{
				{
					Target: "target1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
						{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2"},
					},
				},
				{
					Target: "target2",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0003", PkgName: "pkg3"},
					},
				},
			},
		},
		newReport: trivy.Report{
			Results: []trivy.Result{
				{
					Target: "target1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0001", PkgName: "pkg1"},
					},
				},
				{
					Target: "target2",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "CVE-0000-0003", PkgName: "pkg3"},
					},
				},
			},
		},
		fixed: []trivy.Result{
			{
				Target: "target1",
				Vulnerabilities: []trivy.DetectedVulnerability{
					{VulnerabilityID: "CVE-0000-0002", PkgName: "pkg2"},
				},
			},
		},
		added: nil,
	}))
}
