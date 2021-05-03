package usecase

import (
	"fmt"
	"io"
	"path/filepath"

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
	if err := svc.GetCodeZip(&req.GitHubRepo, req.Ref, req.InstallID, tmp); err != nil {
		return err
	}

	zipFile, err := svc.FS.OpenZip(tmp.Name())
	if err != nil {
		return goerr.Wrap(err).With("file", tmp.Name())
	}

	var newPkgs []*model.Package
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

		parsed := make([]*model.Package, len(pkgs))
		for i := range pkgs {
			parsed[i] = &model.Package{
				ScanTarget: req.ScanTarget,

				Source:  stepDownDirectory(f.Name),
				PkgType: psr.PkgType,
				PkgName: pkgs[i].Name,
				Version: pkgs[i].Version,
			}
		}

		newPkgs = append(newPkgs, parsed...)
	}

	if len(newPkgs) == 0 {
		return nil
	}

	oldPkgs, err := svc.DB().FindPackagesByBranch(&req.GitHubBranch)
	if err != nil {
		return goerr.Wrap(err)
	}

	addPkgs, delPkgs := diffPackageList(oldPkgs, newPkgs)
	for _, pkg := range addPkgs {
		if err := svc.DB().InsertPackage(pkg); err != nil {
			return goerr.Wrap(err).With("pkg", pkg)
		}
	}
	for _, pkg := range delPkgs {
		if err := svc.DB().DeletePackage(pkg); err != nil {
			return goerr.Wrap(err).With("pkg", pkg)
		}
	}

	return nil
}

func mapPackages(pkgs []*model.Package) map[string]*model.Package {
	resp := make(map[string]*model.Package)
	for _, pkg := range pkgs {
		key := fmt.Sprintf("%s/%s/%s", pkg.Source, pkg.PkgName, pkg.Version)
		resp[key] = pkg
	}
	return resp
}

func diffPackageList(oldPkgs, newPkgs []*model.Package) (addPkgs, delPkgs []*model.Package) {
	oldMap := mapPackages(oldPkgs)
	newMap := mapPackages(newPkgs)

	for oldKey, oldPkg := range oldMap {
		if _, ok := newMap[oldKey]; !ok {
			delPkgs = append(delPkgs, oldPkg)
		}
	}

	for newKey, newPkg := range newMap {
		if _, ok := oldMap[newKey]; !ok {
			addPkgs = append(addPkgs, newPkg)
		}
	}

	return
}
