package model

import (
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

type ScanReport struct {
	Repo     Repository       `json:"repo"`
	CommitID string           `json:"commit_id"`
	Sources  []*PackageSource `json:"sources"`
}

func NewScanReport(scan *ent.Scan, statuses []*ent.VulnStatus, now int64) *ScanReport {
	repo := &ent.Repository{}
	if len(scan.Edges.Repository) > 0 {
		repo = scan.Edges.Repository[0]
	}

	db := NewVulnStatusDB(statuses, now)

	srcMap := map[string]*PackageSource{}
	for i := range scan.Edges.Packages {
		pkg := &Package{
			PackageRecord: *scan.Edges.Packages[i],
		}

		src, ok := srcMap[pkg.Source]
		if !ok {
			src = &PackageSource{
				Source: pkg.Source,
			}
			srcMap[pkg.Source] = src
		}

		for _, vuln := range scan.Edges.Packages[i].Edges.Vulnerabilities {
			v := &Vulnerability{
				Vulnerability: *vuln,
				Status:        db.Lookup(&pkg.PackageRecord, vuln.ID),
			}
			if vuln.Edges.CustomSeverity != nil {
				v.CustomSeverity = vuln.Edges.CustomSeverity.Label
			}

			pkg.Vulnerabilities = append(pkg.Vulnerabilities, v)
		}

		src.Packages = append(src.Packages, pkg)
	}

	labels := make([]string, len(repo.Edges.Labels))
	for i := range repo.Edges.Labels {
		labels[i] = repo.Edges.Labels[i].Name
	}

	inventory := &ScanReport{
		Repo: Repository{
			GitHubBranch: GitHubBranch{
				GitHubRepo: GitHubRepo{
					Owner: repo.Owner,
					Name:  repo.Name,
				},
				Branch: scan.Branch,
			},
			Labels: labels,
		},
		Sources: []*PackageSource{},
	}
	for _, v := range srcMap {
		inventory.Sources = append(inventory.Sources, v)
	}
	return inventory
}

type Repository struct {
	GitHubBranch
	Labels        []string `json:"labels"`
	DefaultBranch string   `json:"default_branch"`
}

type PackageSource struct {
	Source   string     `json:"source"`
	Packages []*Package `json:"packages"`
}

type Package struct {
	ent.PackageRecord
	Vulnerabilities []*Vulnerability `json:"vulnerabilities,omitempty"`

	// To remove "edges" field in JSON
	Edges *struct{} `json:"edges,omitempty"`
}

type Vulnerability struct {
	ent.Vulnerability
	Status         *ent.VulnStatus `json:"status,omitempty"`
	CustomSeverity string          `json:"custom_severity"`

	// To remove "edges" field in JSON
	Edges *struct{} `json:"edges,omitempty"`
}
