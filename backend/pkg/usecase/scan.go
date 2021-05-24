package usecase

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"time"

	"github.com/aquasecurity/go-dep-parser/pkg/bundler"
	"github.com/aquasecurity/go-dep-parser/pkg/gomod"
	"github.com/aquasecurity/go-dep-parser/pkg/npm"
	"github.com/aquasecurity/go-dep-parser/pkg/pipenv"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
	"github.com/aquasecurity/go-dep-parser/pkg/yarn"
	"github.com/google/uuid"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
	"github.com/m-mizutani/octovy/backend/pkg/service"
)

type parser struct {
	Parse   func(r io.Reader) ([]types.Library, error)
	PkgType model.PkgType
}

var parserMap = map[string]parser{
	"Gemfile.lock":      {Parse: bundler.Parse, PkgType: model.PkgRubyGems},
	"go.sum":            {Parse: gomod.Parse, PkgType: model.PkgGoModule},
	"Pipfile.lock":      {Parse: pipenv.Parse, PkgType: model.PkgPyPI},
	"yarn.lock":         {Parse: yarn.Parse, PkgType: model.PkgNPM},
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

func putScanResult(scannedAt time.Time, target *model.ScanTarget, pkgs []*model.PackageRecord) error {

	return nil
}

func putPackageRecords(svc *service.Service, branch *model.GitHubBranch, newPkgs []*model.PackageRecord) error {

	oldPkgs, err := svc.DB().FindPackageRecordsByBranch(branch)
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

	return nil
}

func (x *Default) ScanRepository(req *model.ScanRepositoryRequest) error {
	tmp, err := x.svc.Infra.Utils.TempFile("", "*.zip")
	if err != nil {
		return goerr.Wrap(err)
	}
	if err := x.svc.GetCodeZip(&req.GitHubRepo, req.CommitID, req.InstallID, tmp); err != nil {
		return err
	}

	zipFile, err := x.svc.Infra.Utils.OpenZip(tmp.Name())
	if err != nil {
		return goerr.Wrap(err).With("file", tmp.Name())
	}
	defer func() {
		if err := zipFile.Close(); err != nil {
			logger.With("zip", zipFile).With("error", err).Error("Failed to close zip file")
		}
		if err := x.svc.Infra.Utils.Remove(tmp.Name()); err != nil {
			logger.With("filename", tmp.Name()).Error("Failed to remove zip file")
		}
	}()

	dt, err := x.svc.Detector()
	if err != nil {
		return err
	}
	trivyDBMeta, err := dt.TrivyDBMeta()
	if err != nil {
		return err
	}

	var newPkgs []*model.PackageRecord
	detectedVulnMap := map[string]*model.Vulnerability{}
	sourcePkgMap := map[string][]*model.Package{}

	scannedAt := x.svc.Infra.Utils.TimeNow()
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
				vuln.FirstSeenAt = scannedAt.Unix()
				if vuln.Detail.LastModifiedDate != nil {
					vuln.LastModifiedAt = vuln.Detail.LastModifiedDate.Unix()
				}

				detectedVulnMap[vuln.VulnID] = vuln
			}

			pkg := &model.PackageRecord{
				Detected:  req.ScanTarget,
				ScannedAt: scannedAt.Unix(),
				Source:    stepDownDirectory(f.Name),
				Package: model.Package{
					Type:            psr.PkgType,
					Name:            pkgs[i].Name,
					Version:         pkgs[i].Version,
					Vulnerabilities: vulnIDs,
				},
			}
			parsed[i] = pkg

			sourcePkgMap[pkg.Source] = append(sourcePkgMap[pkg.Source], &pkg.Package)
		}

		newPkgs = append(newPkgs, parsed...)
	}

	if len(newPkgs) > 0 && req.IsTargetBranch {
		if err := putPackageRecords(x.svc, &req.GitHubBranch, newPkgs); err != nil {
			return err
		}
	}

	// ScanReport
	var sources []*model.PackageSource
	for src, pkgs := range sourcePkgMap {
		sources = append(sources, &model.PackageSource{
			Source:   src,
			Packages: pkgs,
		})
	}
	report := &model.ScanReport{
		ReportID:    uuid.New().String(),
		Target:      req.ScanTarget,
		Sources:     sources,
		ScannedAt:   scannedAt.Unix(),
		TrivyDBMeta: *trivyDBMeta,
	}
	if err := x.svc.DB().InsertScanReport(report); err != nil {
		return err
	}

	scanLog := report.ToLog()
	if !req.IsPullRequest {
		branch := &model.Branch{
			GitHubBranch:  req.GitHubBranch,
			LastScannedAt: report.ScannedAt,
			ReportSummary: scanLog.Summary,
		}
		if err := x.svc.DB().UpdateBranch(branch); err != nil {
			return err
		}
		if err := x.svc.DB().UpdateBranchIfDefault(&req.GitHubRepo, branch); err != nil {
			return err
		}
	}

	for _, vuln := range detectedVulnMap {
		if err := x.svc.DB().InsertVulnerability(vuln); err != nil {
			return goerr.Wrap(err).With("vuln", vuln)
		}
	}

	logger.With("log", scanLog).Info("Done repository scan")

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

func (x *Default) LookupScanReport(reportID string) (*model.ScanReport, error) {
	return x.svc.DB().LookupScanReport(reportID)
}

func (x *Default) FindPackageRecordsByBranch(branch *model.GitHubBranch) ([]*model.PackageRecord, error) {
	return x.svc.DB().FindPackageRecordsByBranch(branch)
}

func (x *Default) FindPackageRecordsByName(pkgType model.PkgType, pkgName string) ([]*model.PackageRecord, error) {
	return x.svc.DB().FindPackageRecordsByName(pkgType, pkgName)
}
