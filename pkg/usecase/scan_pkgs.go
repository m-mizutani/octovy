package usecase

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

func crawlPackages(req *model.ScanRepositoryRequest, clients *scanClients) ([]*ent.PackageRecord, error) {
	tmp, err := ioutil.TempFile("", "*.zip")
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	defer func() {
		if err := os.Remove(tmp.Name()); err != nil {
			logger.Error().Interface("filename", tmp.Name()).Msg("Failed to remove zip file")
		}
	}()

	if err := clients.GitHubApp.GetCodeZip(&req.GitHubRepo, req.CommitID, tmp); err != nil {
		return nil, err
	}

	zipFile, err := zip.OpenReader(tmp.Name())
	if err != nil {
		return nil, goerr.Wrap(err).With("file", tmp.Name())
	}
	defer func() {
		if err := zipFile.Close(); err != nil {
			logger.Error().Interface("zip", zipFile).Err(err).Msg("Failed to close zip file")
		}
	}()

	var newPkgs []*ent.PackageRecord

	for _, f := range zipFile.File {
		psr, ok := parserMap[filepath.Base(f.Name)]
		if !ok {
			continue
		}

		fd, err := f.Open()
		if err != nil {
			return nil, goerr.Wrap(err)
		}
		defer fd.Close()

		pkgs, err := psr.Parse(fd)
		if err != nil {
			return nil, goerr.Wrap(err)
		}

		parsed := make([]*ent.PackageRecord, len(pkgs))
		for i := range pkgs {
			pkg := &ent.PackageRecord{
				Source:  stepDownDirectory(f.Name),
				Type:    psr.PkgType,
				Name:    pkgs[i].Name,
				Version: pkgs[i].Version,
			}
			parsed[i] = pkg
		}

		newPkgs = append(newPkgs, parsed...)
	}

	return newPkgs, nil
}

type pkgChanges struct {
	Added    []*ent.PackageRecord
	Modified []*ent.PackageRecord
	Deleted  []*ent.PackageRecord
}

func diffPackages(oldPkgs, newPkgs []*ent.PackageRecord) *pkgChanges {
	var changes pkgChanges

	oldMap := mapPackages(oldPkgs)
	newMap := mapPackages(newPkgs)

	for oldKey, oldPkg := range oldMap {
		if newPkg, ok := newMap[oldKey]; !ok {
			changes.Deleted = append(changes.Deleted, oldPkg)
		} else {
			if !matchVulnerabilities(oldPkg, newPkg) {
				changes.Modified = append(changes.Modified, newPkg)
			}
		}
	}

	for newKey, newPkg := range newMap {
		if _, ok := oldMap[newKey]; !ok {
			changes.Added = append(changes.Added, newPkg)
		}
	}

	return &changes
}

func mapPackages(pkgs []*ent.PackageRecord) map[string]*ent.PackageRecord {
	resp := make(map[string]*ent.PackageRecord)
	for _, pkg := range pkgs {
		key := fmt.Sprintf("%s|%s|%s", pkg.Source, pkg.Name, pkg.Version)
		resp[key] = pkg
	}
	return resp
}
