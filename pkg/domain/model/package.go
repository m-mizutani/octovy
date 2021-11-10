package model

import (
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

type PackageInventory struct {
	Sources []*PackageSource
}

func NewPackageInventory(pkgs []*ent.PackageRecord, statuses []*ent.VulnStatus, now int64) *PackageInventory {
	db := NewVulnStatusDB(statuses, now)

	srcMap := map[string]*PackageSource{}
	for i := range pkgs {
		pkg := &Package{
			PackageRecord: *pkgs[i],
		}

		src, ok := srcMap[pkg.Source]
		if !ok {
			src = &PackageSource{
				Source: pkg.Source,
			}
			srcMap[pkg.Source] = src
		}

		for _, vuln := range pkgs[i].Edges.Vulnerabilities {
			pkg.Vulnerabilities = append(pkg.Vulnerabilities, &Vulnerability{
				Vulnerability:  *vuln,
				Status:         db.Lookup(&pkg.PackageRecord, vuln.ID),
				CustomSeverity: vuln.Edges.CustomSeverity,
			})
		}

		src.Packages = append(src.Packages, pkg)
	}

	inventory := &PackageInventory{
		Sources: make([]*PackageSource, len(srcMap)),
	}
	for _, v := range srcMap {
		inventory.Sources = append(inventory.Sources, v)
	}
	return inventory
}

func (x *PackageInventory) CheckResult(rules []*ent.CheckRule) types.GitHubCheckResult {
	return types.CheckSuccess
}

type PackageSource struct {
	Source   string     `json:"source"`
	Packages []*Package `json:"packages"`

	// To remove "edges" field in JSON
	Edges *struct{} `json:"edges.omitempty"`
}

type Package struct {
	ent.PackageRecord
	Vulnerabilities []*Vulnerability `json:"vulnerabilities"`

	// To remove "edges" field in JSON
	Edges *struct{} `json:"edges.omitempty"`
}

type Vulnerability struct {
	ent.Vulnerability
	Status         *ent.VulnStatus `json:"status"`
	CustomSeverity *ent.Severity   `json:"custom_severity"`

	// To remove "edges" field in JSON
	Edges *struct{} `json:"edges.omitempty"`
}

func (x *Vulnerability) Feedback(repo *ent.Repository, rules []*ent.CheckRule) types.GitHubCheckResult {
	for _, rule := range rules {
		return rule.Result
	}

	return types.CheckSuccess
}
