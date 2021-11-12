package model

import (
	"fmt"

	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

type VulnStatusDB struct {
	dict map[string]*ent.VulnStatus
}

func vulnStatusKey(src, pkgName, vulnID string) string {
	return fmt.Sprintf("%s|%s|%s", src, pkgName, vulnID)
}

func NewVulnStatusDB(statuses []*ent.VulnStatus, now int64) *VulnStatusDB {
	db := &VulnStatusDB{
		dict: make(map[string]*ent.VulnStatus),
	}
	for _, status := range statuses {
		key := vulnStatusKey(status.Source, status.PkgName, status.VulnID)
		if status.Status == types.StatusNone {
			continue
		}
		if status.Status == types.StatusSnoozed && status.ExpiresAt < now {
			continue
		}
		db.dict[key] = status
	}
	return db
}

func (x *VulnStatusDB) Lookup(pkg *ent.PackageRecord, vulnID string) *ent.VulnStatus {
	if status, ok := x.dict[vulnStatusKey(pkg.Source, pkg.Name, vulnID)]; ok {
		return status
	}
	return nil
}

func (x *VulnStatusDB) IsQualified(v *VulnRecord) bool {
	_, ok := x.dict[vulnStatusKey(v.Pkg.Source, v.Pkg.Name, v.Vuln.ID)]
	return !ok
}

type vulnChangeType int

const (
	VulnAdded vulnChangeType = iota
	VulnDeleted
	VulnRemained
)

type VulnRecord struct {
	Pkg  *ent.PackageRecord
	Vuln *ent.Vulnerability
}

type vulnChange struct {
	VulnRecord
	Type vulnChangeType
}
type VulnChanges []*vulnChange

func (x VulnChanges) Qualified(db *VulnStatusDB) VulnChanges {
	var resp VulnChanges
	for i := range x {
		if db.IsQualified(&x[i].VulnRecord) {
			resp = append(resp, x[i])
		}
	}
	return resp
}

func (x VulnChanges) FilterByType(t vulnChangeType) VulnChanges {
	var resp VulnChanges
	for i := range x {
		if x[i].Type == t {
			resp = append(resp, x[i])
		}
	}
	return resp
}

func (x VulnChanges) FilterBySource(src string) VulnChanges {
	var resp VulnChanges
	for i := range x {
		if x[i].Pkg.Source == src {
			resp = append(resp, x[i])
		}
	}
	return resp
}
func (x VulnChanges) Sources() []string {
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

type vulnRecordMap map[string]*VulnRecord

func (x vulnRecordMap) Put(pkg *ent.PackageRecord, vuln *ent.Vulnerability) {
	key := vulnStatusKey(pkg.Source, pkg.Name, vuln.ID)
	x[key] = &VulnRecord{
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
				VulnRecord: *oldVuln,
				Type:       VulnDeleted,
			})
		}
	}

	for newKey, newVuln := range newMap {
		if _, ok := oldMap[newKey]; !ok {
			changes = append(changes, &vulnChange{
				VulnRecord: *newVuln,
				Type:       VulnAdded,
			})
		} else {
			changes = append(changes, &vulnChange{
				VulnRecord: *newVuln,
				Type:       VulnRemained,
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

func DiffVulnRecords(oldPkgs, newPkgs []*ent.PackageRecord) VulnChanges {
	oldMap := pkgToVulnRecordMap(oldPkgs)
	newMap := pkgToVulnRecordMap(newPkgs)

	var changes []*vulnChange

	for oldKey, oldVuln := range oldMap {
		if _, ok := newMap[oldKey]; !ok {
			changes = append(changes, &vulnChange{
				VulnRecord: *oldVuln,
				Type:       VulnDeleted,
			})
		} else {
			changes = append(changes, &vulnChange{
				VulnRecord: *oldVuln,
				Type:       VulnRemained,
			})
		}
	}

	for newKey, newVuln := range newMap {
		if _, ok := oldMap[newKey]; !ok {
			changes = append(changes, &vulnChange{
				VulnRecord: *newVuln,
				Type:       VulnAdded,
			})
		}
	}

	return changes
}
