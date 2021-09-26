package model

type SourceChanges struct {
	Added    VulnChanges
	Deleted  VulnChanges
	Remained VulnChanges
}

type Report struct {
	Sources map[string]*SourceChanges
}

func MakeReport(changes VulnChanges, db *VulnStatusDB) *Report {
	report := &Report{
		Sources: make(map[string]*SourceChanges),
	}
	for _, src := range changes.Sources() {
		target := changes.FilterBySource(src)
		qualified := target.Qualified(db)

		report.Sources[src] = &SourceChanges{
			Added:    qualified.FilterByType(VulnAdded),
			Deleted:  target.FilterByType(VulnDeleted),
			Remained: qualified.FilterByType(VulnRemained),
		}
	}

	return report
}

func (x *Report) NothingToNotify(githubEvent string) bool {
	switch githubEvent {
	case "opened":
		for _, src := range x.Sources {
			if len(src.Added) > 0 || len(src.Deleted) > 0 || len(src.Remained) > 0 {
				return false
			}
		}
		return true

	case "synchronize":
		for _, src := range x.Sources {
			if len(src.Added) > 0 || len(src.Deleted) > 0 {
				return false
			}
		}
		return true

	default:
		panic("unsupported github event: " + githubEvent)
	}
}
