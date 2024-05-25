package logic

import "github.com/m-mizutani/octovy/pkg/domain/model/trivy"

func DiffResults(oldReport, newReport *trivy.Report) (fixed, added trivy.Results) {
	resultMap := map[string]trivy.Result{}
	for _, result := range oldReport.Results {
		resultMap[result.Target] = result
	}

	for _, newResult := range newReport.Results {
		oldResult, ok := resultMap[newResult.Target]
		if !ok {
			added = append(added, newResult)
			continue
		}

		fixedVuln, addedVuln := DiffVulnerabilities(oldResult.Vulnerabilities, newResult.Vulnerabilities)
		if len(fixedVuln) > 0 {
			fixedResult := oldResult
			fixedResult.Vulnerabilities = fixedVuln
			fixed = append(fixed, fixedResult)
		}

		if len(addedVuln) > 0 {
			addedResult := newResult
			addedResult.Vulnerabilities = addedVuln
			added = append(added, addedResult)
		}

		delete(resultMap, newResult.Target)
	}

	for _, result := range resultMap {
		fixed = append(fixed, result)
	}

	return
}

func DiffVulnerabilities(oldVulns, newVulns []trivy.DetectedVulnerability) (fixed, added []trivy.DetectedVulnerability) {
	oldVulnMap := map[string]trivy.DetectedVulnerability{}
	for _, vuln := range oldVulns {

		oldVulnMap[vuln.ID()] = vuln
	}

	for _, newVuln := range newVulns {
		if _, ok := oldVulnMap[newVuln.ID()]; !ok {
			added = append(added, newVuln)
		}
		delete(oldVulnMap, newVuln.ID())
	}

	for _, vuln := range oldVulnMap {
		fixed = append(fixed, vuln)
	}

	return
}
