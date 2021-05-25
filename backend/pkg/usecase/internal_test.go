package usecase

import (
	"testing"

	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestStepDownDirectory(t *testing.T) {
	assert.Equal(t, "blue/Gemfile.lock", stepDownDirectory("root/blue/Gemfile.lock"))
	assert.Equal(t, "blue/Gemfile.lock", stepDownDirectory("./root/blue/Gemfile.lock"))
	assert.Equal(t, "blue/green/Gemfile.lock", stepDownDirectory("/root/blue/green/Gemfile.lock"))
}

func TestDiffReport(t *testing.T) {
	newReport := &model.ScanReport{
		Sources: []*model.PackageSource{
			{
				Source: "abc",
				Packages: []*model.Package{
					{
						Name:            "blue",
						Version:         "1.1",
						Vulnerabilities: []string{"CVE-2999-0001"},
					},
					{
						Name:            "orange",
						Version:         "2.1",
						Vulnerabilities: []string{"CVE-2999-0002"},
					},
					{
						Name:            "red",
						Version:         "3.1",
						Vulnerabilities: []string{"CVE-2999-0003"},
					},
				},
			},
		},
	}
	oldReport := &model.ScanReport{
		Sources: []*model.PackageSource{
			{
				Source: "abc",
				Packages: []*model.Package{
					{
						Name:            "orange",
						Version:         "2.1",
						Vulnerabilities: []string{"CVE-2999-0002", "CVE-2999-0011"},
					},
					{
						Name:            "red",
						Version:         "3.1",
						Vulnerabilities: []string{"CVE-2999-0003"},
					},
					{
						Name:            "timeless",
						Version:         "5.0",
						Vulnerabilities: []string{"CVE-2999-0005"},
					},
				},
			},
			{
				Source: "xyz",
				Packages: []*model.Package{
					{
						Name:            "red",
						Version:         "3.1",
						Vulnerabilities: []string{"CVE-2999-0003"},
					},
				},
			},
		},
	}

	newVuln, fixedVuln, remainedVuln := diffReport(newReport, oldReport)
	assert.Equal(t, 3, len(newVuln))
	assert.Equal(t, 2, len(remainedVuln))
	assert.Equal(t, 1, len(fixedVuln))

	assert.Contains(t, newVuln, &vulnRecord{
		Source:     "abc",
		VulnID:     "CVE-2999-0011",
		PkgName:    "orange",
		PkgVersion: "2.1",
	})
	assert.Contains(t, newVuln, &vulnRecord{
		Source:     "abc",
		VulnID:     "CVE-2999-0005",
		PkgName:    "timeless",
		PkgVersion: "5.0",
	})
	assert.Contains(t, newVuln, &vulnRecord{
		Source:     "xyz",
		VulnID:     "CVE-2999-0003",
		PkgName:    "red",
		PkgVersion: "3.1",
	})
}
