package logic_test

import (
	"testing"
	"time"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/logic"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/model/trivy"
)

func TestFilterResults(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		results  trivy.Results
		cfg      *model.Config
		expected trivy.Results
	}{
		{
			name: "No ignore targets",
			results: trivy.Results{
				{
					Target: "file1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "vuln1"},
						{VulnerabilityID: "vuln2"},
					},
				},
			},
			cfg: &model.Config{
				IgnoreTargets: []model.IgnoreTarget{},
			},
			expected: trivy.Results{
				{
					Target: "file1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "vuln1"},
						{VulnerabilityID: "vuln2"},
					},
				},
			},
		},
		{
			name: "Ignore expired vulnerabilities",
			results: trivy.Results{
				{
					Target: "file1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "vuln1"},
						{VulnerabilityID: "vuln2"},
					},
				},
			},
			cfg: &model.Config{
				IgnoreTargets: []model.IgnoreTarget{
					{
						File: "file1",
						Vulns: []model.IgnoreVuln{
							{ID: "vuln1", ExpiresAt: now.Add(-time.Hour)},
							{ID: "vuln2", ExpiresAt: now.Add(time.Hour)},
						},
					},
				},
			},
			expected: trivy.Results{
				{
					Target: "file1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "vuln1"},
					},
				},
			},
		},
		{
			name: "Ignore non-expired vulnerabilities",
			results: trivy.Results{
				{
					Target: "file1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "vuln1"},
						{VulnerabilityID: "vuln2"},
					},
				},
			},
			cfg: &model.Config{
				IgnoreTargets: []model.IgnoreTarget{
					{
						File: "file1",
						Vulns: []model.IgnoreVuln{
							{ID: "vuln1", ExpiresAt: now.Add(time.Hour)},
						},
					},
				},
			},
			expected: trivy.Results{
				{
					Target: "file1",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "vuln2"},
					},
				},
			},
		},
		{
			name: "No vulnerabilities to ignore",
			results: trivy.Results{
				{
					Target: "file2",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "vuln3"},
					},
				},
			},
			cfg: &model.Config{
				IgnoreTargets: []model.IgnoreTarget{
					{
						File: "file1",
						Vulns: []model.IgnoreVuln{
							{ID: "vuln1", ExpiresAt: now.Add(time.Hour)},
						},
					},
				},
			},
			expected: trivy.Results{
				{
					Target: "file2",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "vuln3"},
					},
				},
			},
		},
		{
			name: "Ignore vulnerabilities for different file",
			results: trivy.Results{
				{
					Target: "file2",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "vuln3"},
					},
				},
			},
			cfg: &model.Config{
				IgnoreTargets: []model.IgnoreTarget{
					{
						File: "file1",
						Vulns: []model.IgnoreVuln{
							{ID: "vuln1", ExpiresAt: now.Add(time.Hour)},
						},
					},
				},
			},
			expected: trivy.Results{
				{
					Target: "file2",
					Vulnerabilities: []trivy.DetectedVulnerability{
						{VulnerabilityID: "vuln3"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := logic.FilterResults(tt.results, tt.cfg, now)
			gt.Equal(t, tt.expected, actual)
		})
	}
}
