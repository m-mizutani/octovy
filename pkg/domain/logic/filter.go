package logic

import (
	"time"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/model/trivy"
)

func FilterReport(oldReport *trivy.Report, cfg *model.Config, now time.Time) *trivy.Report {
	results := FilterResults(oldReport.Results, cfg, now)
	newReport := *oldReport
	newReport.Results = results
	return &newReport
}

func FilterResults(results trivy.Results, cfg *model.Config, now time.Time) trivy.Results {
	ignoreMap := make(map[string]map[string]struct{})
	for _, target := range cfg.IgnoreList {
		if _, ok := ignoreMap[target.Target]; !ok {
			ignoreMap[target.Target] = make(map[string]struct{})
		}

		for _, vuln := range target.Vulns {
			if vuln.ExpiresAt.Before(now) {
				continue
			}
			ignoreMap[target.Target][vuln.ID] = struct{}{}
		}
	}

	var filtered trivy.Results
	for _, result := range results {
		newResult := result
		ignoreVulns, ok := ignoreMap[result.Target]
		if !ok {
			filtered = append(filtered, newResult)
			continue
		}
		newResult.Vulnerabilities = nil

		for _, vuln := range result.Vulnerabilities {
			if _, ok := ignoreVulns[vuln.VulnerabilityID]; ok {
				continue
			}

			newResult.Vulnerabilities = append(newResult.Vulnerabilities, vuln)
		}

		if len(newResult.Vulnerabilities) > 0 {
			filtered = append(filtered, newResult)
		}
	}
	return filtered
}
