package usecase

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"

	"github.com/aquasecurity/go-dep-parser/pkg/bundler"
	"github.com/aquasecurity/go-dep-parser/pkg/gomod"
	"github.com/aquasecurity/go-dep-parser/pkg/npm"
	"github.com/aquasecurity/go-dep-parser/pkg/pipenv"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
	"github.com/aquasecurity/go-dep-parser/pkg/yarn"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/m-mizutani/octovy/backend/pkg/service"
)

type parser struct {
	Parse   func(r io.Reader) ([]types.Library, error)
	PkgType model.PkgType
}

var parserMap = map[string]parser{
	"Gemfile.lock":      {Parse: bundler.Parse, PkgType: model.PkgBundler},
	"go.sum":            {Parse: gomod.Parse, PkgType: model.PkgGoModule},
	"Pipfile.lock":      {Parse: pipenv.Parse, PkgType: model.PkgPipenv},
	"yarn.lock":         {Parse: yarn.Parse, PkgType: model.PkgYarn},
	"package-lock.json": {Parse: npm.Parse, PkgType: model.PkgNPM},
}

func stepDownDirectory(fpath string) string {
	if len(fpath) > 0 && fpath[0] == filepath.Separator {
		fpath = fpath[1:]
	}

	p := fpath
	var arr []string
	for {
		d, f := filepath.Split(p)
		if d == "" {
			break
		}
		arr = append([]string{f}, arr...)
		p = filepath.Clean(d)
	}

	return filepath.Join(arr...)
}

func (x *Default) ScanRepository(svc *service.Service, req *model.ScanRepositoryRequest) error {
	tmp, err := svc.FS.TempFile("", "*.zip")
	if err != nil {
		return goerr.Wrap(err)
	}
	if err := svc.GetCodeZip(&req.GitHubRepo, req.CommitID, req.InstallID, tmp); err != nil {
		return err
	}

	zipFile, err := svc.FS.OpenZip(tmp.Name())
	if err != nil {
		return goerr.Wrap(err).With("file", tmp.Name())
	}

	dt, err := svc.Detector()
	if err != nil {
		return err
	}

	var newPkgs []*model.PackageRecord
	detectedVulnMap := map[string]*model.Vulnerability{}

	for _, f := range zipFile.File {
		psr, ok := parserMap[filepath.Base(f.Name)]
		if !ok {
			continue
		}

		fd, err := f.Open()
		if err != nil {
			return goerr.Wrap(err)
		}
		defer fd.Close()

		pkgs, err := psr.Parse(fd)
		if err != nil {
			return goerr.Wrap(err)
		}

		parsed := make([]*model.PackageRecord, len(pkgs))
		for i := range pkgs {
			vulns, err := dt.Detect(psr.PkgType, pkgs[i].Name, pkgs[i].Version)
			if err != nil {
				return err
			}

			var vulnIDs []string
			for _, vuln := range vulns {
				vulnIDs = append(vulnIDs, vuln.VulnID)
				detectedVulnMap[vuln.VulnID] = vuln
			}

			parsed[i] = &model.PackageRecord{
				Detected: req.ScanTarget,

				Source: stepDownDirectory(f.Name),
				Package: model.Package{
					Type:            psr.PkgType,
					Name:            pkgs[i].Name,
					Version:         pkgs[i].Version,
					Vulnerabilities: vulnIDs,
				},
			}
		}

		newPkgs = append(newPkgs, parsed...)
	}

	if len(newPkgs) == 0 {
		return nil
	}

	oldPkgs, err := svc.DB().FindPackageRecordsByBranch(&req.GitHubBranch)
	if err != nil {
		return goerr.Wrap(err)
	}

	addPkgs, modPkgs, delPkgs := diffPackageList(oldPkgs, newPkgs)
	for _, pkg := range addPkgs {
		if inserted, err := svc.DB().InsertPackageRecord(pkg); err != nil {
			return goerr.Wrap(err).With("pkg", pkg)
		} else if !inserted {
			modPkgs = append(modPkgs, pkg)
		}
	}

	for _, pkg := range modPkgs {
		if err := svc.DB().UpdatePackageRecord(pkg); err != nil {
			return goerr.Wrap(err).With("pkg", pkg)
		}
	}

	for _, pkg := range delPkgs {
		if err := svc.DB().RemovePackageRecord(pkg); err != nil {
			return goerr.Wrap(err).With("pkg", pkg)
		}
	}

	for _, vuln := range detectedVulnMap {
		if err := svc.DB().InsertVulnerability(vuln); err != nil {
			return goerr.Wrap(err).With("vuln", vuln)
		}
	}

	return nil
}

func mapPackages(pkgs []*model.PackageRecord) map[string]*model.PackageRecord {
	resp := make(map[string]*model.PackageRecord)
	for _, pkg := range pkgs {
		key := fmt.Sprintf("%s|%s|%s", pkg.Source, pkg.Name, pkg.Version)
		resp[key] = pkg
	}
	return resp
}

func matchVulnerabilities(a, b *model.PackageRecord) bool {
	copyVulnList := func(p *model.PackageRecord) []string {
		v := make([]string, len(p.Vulnerabilities))
		for i := range p.Vulnerabilities {
			v[i] = p.Vulnerabilities[i]
		}
		sort.Slice(v, func(i int, j int) bool {
			return v[i] < v[j]
		})
		return v
	}

	v1 := copyVulnList(a)
	v2 := copyVulnList(b)
	if len(v1) != len(v2) {
		return false
	}
	for i := range v1 {
		if v1[i] != v2[i] {
			return false
		}
	}
	return true
}

func diffPackageList(oldPkgs, newPkgs []*model.PackageRecord) (addPkgs, modPkgs, delPkgs []*model.PackageRecord) {
	oldMap := mapPackages(oldPkgs)
	newMap := mapPackages(newPkgs)

	for oldKey, oldPkg := range oldMap {
		if newPkg, ok := newMap[oldKey]; !ok {
			delPkgs = append(delPkgs, oldPkg)
		} else {
			if !matchVulnerabilities(oldPkg, newPkg) {
				modPkgs = append(modPkgs, newPkg)
			}
		}
	}

	for newKey, newPkg := range newMap {
		if _, ok := oldMap[newKey]; !ok {
			addPkgs = append(addPkgs, newPkg)
		}
	}

	return
}
