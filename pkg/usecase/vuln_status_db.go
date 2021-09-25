package usecase

import (
	"fmt"

	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

type vulnStatusDB struct {
	dict map[string]*ent.VulnStatus
}

func vulnStatusKey(src, pkgName, vulnID string) string {
	return fmt.Sprintf("%s|%s|%s", src, pkgName, vulnID)
}

func newVulnStatusDB(statuses []*ent.VulnStatus, now int64) *vulnStatusDB {
	db := &vulnStatusDB{
		dict: make(map[string]*ent.VulnStatus),
	}
	for _, status := range statuses {
		key := vulnStatusKey(status.Source, status.PkgName, status.VulnID)
		if now < status.ExpiresAt && status.Status != types.StatusNone {
			db.dict[key] = status
		}
	}
	return db
}

func (x *vulnStatusDB) Filter(pkg *ent.PackageRecord) []*ent.Vulnerability {
	var resp []*ent.Vulnerability
	for _, vuln := range pkg.Edges.Vulnerabilities {
		if _, ok := x.dict[vulnStatusKey(pkg.Source, pkg.Name, vuln.ID)]; !ok {
			resp = append(resp, vuln)
		}
	}
	return resp
}

func (x *vulnStatusDB) IsIgnored(v *vulnRecord) bool {
	_, ok := x.dict[vulnStatusKey(v.Pkg.Source, v.Pkg.Name, v.Vuln.ID)]
	return ok
}
