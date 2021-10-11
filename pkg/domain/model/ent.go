package model

import (
	"strings"
	"time"

	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/utils"
)

var logger = utils.Logger

func TrivyReportToEnt(report *TrivyReport, now time.Time) (pkgList []*ent.PackageRecord, vulnList []*ent.Vulnerability) {
	ptrKey := func(name, ver string) string {
		return name + "|" + ver
	}
	vulnMap := map[string]*ent.Vulnerability{}

	for _, result := range report.Results {
		pkgPtr := map[string]*ent.PackageRecord{}
		for _, pkg := range result.Packages {
			p := &ent.PackageRecord{
				Type:    result.Type,
				Source:  result.Target,
				Name:    pkg.Name,
				Version: pkg.Version,
			}
			pkgList = append(pkgList, p)
			pkgPtr[ptrKey(pkg.Name, pkg.Version)] = p
		}

		for _, vuln := range result.Vulnerabilities {
			var cvss []string
			for vendor, v := range vuln.CVSS {
				if v.V2Vector != "" {
					cvss = append(cvss, strings.Join([]string{vendor, "V2Vector", v.V2Vector}, ","))
				}
				if v.V3Vector != "" {
					cvss = append(cvss, strings.Join([]string{vendor, "V3Vector", v.V3Vector}, ","))
				}
			}

			v := &ent.Vulnerability{
				ID:          vuln.VulnerabilityID,
				FirstSeenAt: now.Unix(),
				Title:       vuln.Title,
				Description: vuln.Description,
				CweID:       vuln.CweIDs,
				Severity:    vuln.Severity,
				Cvss:        cvss,
				References:  vuln.References,
			}
			if vuln.LastModifiedDate != nil {
				v.LastModifiedAt = vuln.LastModifiedDate.Unix()
			}

			vulnMap[v.ID] = v

			if pkg, ok := pkgPtr[ptrKey(vuln.PkgName, vuln.InstalledVersion)]; ok {
				pkg.VulnIds = append(pkg.VulnIds, vuln.VulnerabilityID)
			} else {
				logger.Warn().Interface("vuln", vuln).Msg("package is not inserted")
			}
		}
	}

	for _, v := range vulnMap {
		vulnList = append(vulnList, v)
	}

	return
}
