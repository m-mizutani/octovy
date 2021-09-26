package usecase

import "github.com/m-mizutani/octovy/pkg/infra/ent"

type vulnChangeType int

const (
	vulnAdded vulnChangeType = iota
	vulnDeleted
	vulnRemained
)

type vulnRecord struct {
	Pkg  *ent.PackageRecord
	Vuln *ent.Vulnerability
}

type vulnChange struct {
	vulnRecord
	Type vulnChangeType
}
type vulnChanges []*vulnChange

func (x vulnChanges) Qualified(db *vulnStatusDB) vulnChanges {
	var resp vulnChanges
	for i := range x {
		if db.IsIgnored(&x[i].vulnRecord) {
			continue
		}
		resp = append(resp, x[i])
	}
	return resp
}

func (x vulnChanges) FilterByType(t vulnChangeType) vulnChanges {
	var resp vulnChanges
	for i := range x {
		if x[i].Type == t {
			resp = append(resp, x[i])
		}
	}
	return resp
}

func (x vulnChanges) FilterBySource(src string) vulnChanges {
	var resp vulnChanges
	for i := range x {
		if x[i].Pkg.Source == src {
			resp = append(resp, x[i])
		}
	}
	return resp
}
func (x vulnChanges) Sources() []string {
	src := map[string]struct{}{}
	for i := range x {
		src[x[i].Pkg.Source] = struct{}{}
	}

	var srcList []string
	for s := range src {
		srcList = append(srcList, s)
	}
	return srcList
}

type vulnRecordMap map[string]*vulnRecord

func (x vulnRecordMap) Put(pkg *ent.PackageRecord, vuln *ent.Vulnerability) {
	key := vulnStatusKey(pkg.Source, pkg.Name, vuln.ID)
	x[key] = &vulnRecord{
		Pkg:  pkg,
		Vuln: vuln,
	}
}
func (x vulnRecordMap) Diff(y vulnRecordMap) []*vulnChange {
	oldMap, newMap := x, y
	var changes []*vulnChange

	for oldKey, oldVuln := range oldMap {
		if _, ok := newMap[oldKey]; !ok {
			changes = append(changes, &vulnChange{
				vulnRecord: *oldVuln,
				Type:       vulnDeleted,
			})
		} else {
			changes = append(changes, &vulnChange{
				vulnRecord: *oldVuln,
				Type:       vulnRemained,
			})
		}
	}

	for newKey, newVuln := range newMap {
		if _, ok := oldMap[newKey]; !ok {
			changes = append(changes, &vulnChange{
				vulnRecord: *newVuln,
				Type:       vulnAdded,
			})
		}
	}

	return changes
}

func pkgToVulnRecordMap(pkgs []*ent.PackageRecord) vulnRecordMap {
	m := vulnRecordMap{}
	for _, pkg := range pkgs {
		for _, vuln := range pkg.Edges.Vulnerabilities {
			m.Put(pkg, vuln)
		}
	}
	return m
}
