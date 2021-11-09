package model

import (
	"fmt"
	"strings"
)

type SourceChanges struct {
	Added    VulnChanges
	Deleted  VulnChanges
	Remained VulnChanges
}

type Report struct {
	scanID      string
	frontendURL string
	sources     map[string]*SourceChanges
}

func MakeReport(scanID string, changes VulnChanges, db *VulnStatusDB, url string) *Report {
	report := &Report{
		scanID:      scanID,
		sources:     make(map[string]*SourceChanges),
		frontendURL: strings.Trim(url, "/"),
	}
	for _, src := range changes.Sources() {
		target := changes.FilterBySource(src)
		qualified := target.Qualified(db)

		report.sources[src] = &SourceChanges{
			Added:    qualified.FilterByType(VulnAdded),
			Deleted:  target.FilterByType(VulnDeleted),
			Remained: qualified.FilterByType(VulnRemained),
		}
	}

	return report
}

func (x *Report) ToMarkdown() string {
	var b githubCommentBody
	b.Add("## Octovy scan result")
	b.Break()

	for src, changes := range x.sources {
		b.Add("### " + src)
		b.Break()

		for _, v := range changes.Added {
			b.Add("- 🚨 **New** %s (%s): %s", v.Vuln.ID, v.Pkg.Name, v.Vuln.Title)
		}
		for _, v := range changes.Deleted {
			b.Add("- ✅ **Fixed** %s (%s): %s", v.Vuln.ID, v.Pkg.Name, v.Vuln.Title)
		}
		if len(changes.Remained) > 0 {
			b.Add("- ⚠️ %d vulnerabilities are remained", len(changes.Remained))
		}
		b.Break()
	}

	b.Add("🗒️ See [report](%s/scan/%s) more detail", x.frontendURL, x.scanID)

	return b.Join()
}

type githubCommentBody struct {
	lines []string
}

func (x *githubCommentBody) Add(f string, v ...interface{}) {
	x.lines = append(x.lines, fmt.Sprintf(f, v...))
}
func (x *githubCommentBody) Break() {
	x.lines = append(x.lines, "")
}
func (x *githubCommentBody) Join() string { return strings.Join(x.lines, "\n") }

func (x *Report) NothingToNotify(githubEvent string) bool {
	switch githubEvent {
	case "opened":
		for _, src := range x.sources {
			if len(src.Added) > 0 || len(src.Deleted) > 0 || len(src.Remained) > 0 {
				return false
			}
		}
		return true

	case "synchronize":
		for _, src := range x.sources {
			if len(src.Added) > 0 || len(src.Deleted) > 0 {
				return false
			}
		}
		return true

	default:
		panic("unsupported github event: " + githubEvent)
	}
}
