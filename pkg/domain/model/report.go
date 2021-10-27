package model

type SourceChanges struct {
	Added    VulnChanges
	Deleted  VulnChanges
	Remained VulnChanges
}

type Advisory struct {
	Sources map[string]*SourceChanges
}

func MakeAdvisory(changes VulnChanges, db *VulnStatusDB) *Advisory {
	advisory := &Advisory{
		Sources: make(map[string]*SourceChanges),
	}
	for _, src := range changes.Sources() {
		target := changes.FilterBySource(src)
		qualified := target.Qualified(db)

		advisory.Sources[src] = &SourceChanges{
			Added:    qualified.FilterByType(VulnAdded),
			Deleted:  target.FilterByType(VulnDeleted),
			Remained: qualified.FilterByType(VulnRemained),
		}
	}

	return advisory
}

func (x *Advisory) NothingToNotify(githubEvent string) bool {
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
