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
			pkg.Vulnerabilities = append(pkg.Vulnerabilities, &Vulnerability{
				Vulnerability:  *vuln,
				Status:         db.Lookup(&pkg.PackageRecord, vuln.ID),
				CustomSeverity: vuln.Edges.CustomSeverity,
			})
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
					Owner:    repo.Owner,
					RepoName: repo.Name,
				},
				Branch: scan.Branch,
			},
			Labels: labels,
		},
		Sources: make([]*PackageSource, len(srcMap)),
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
	CustomSeverity *ent.Severity   `json:"custom_severity,omitempty"`

	// To remove "edges" field in JSON
	Edges *struct{} `json:"edges,omitempty"`
}
