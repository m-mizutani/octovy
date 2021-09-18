package usecase

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"

	"github.com/aquasecurity/go-dep-parser/pkg/golang/mod"
	"github.com/aquasecurity/go-dep-parser/pkg/nodejs/npm"
	"github.com/aquasecurity/go-dep-parser/pkg/nodejs/yarn"
	"github.com/aquasecurity/go-dep-parser/pkg/python/pipenv"
	"github.com/aquasecurity/go-dep-parser/pkg/ruby/bundler"
	psrTypes "github.com/aquasecurity/go-dep-parser/pkg/types"
)

type parser struct {
	Parse   func(r io.Reader) ([]psrTypes.Library, error)
	PkgType types.PkgType
}

var parserMap = map[string]parser{
	"Gemfile.lock":      {Parse: bundler.Parse, PkgType: types.PkgRubyGems},
	"go.sum":            {Parse: mod.Parse, PkgType: types.PkgGoModule},
	"Pipfile.lock":      {Parse: pipenv.Parse, PkgType: types.PkgPyPI},
	"yarn.lock":         {Parse: yarn.Parse, PkgType: types.PkgNPM},
	"package-lock.json": {Parse: npm.Parse, PkgType: types.PkgNPM},
}

func (x *usecase) SendScanRequest(req *model.ScanRepositoryRequest) error {
	x.scanQueue <- req
	return nil
}

func (x *usecase) InvokeScanThread() {
	go func() {
		if err := x.runScanThread(); err != nil {
			x.handleError(err)
		}
	}()
}

func (x *usecase) runScanThread() error {
	githubAppPEM, err := x.infra.Utils.ReadFile(x.config.GitHubAppPrivateKeyPath)
	if err != nil {
		return goerr.Wrap(err, "Failed to read github private key file")
	}

	detector := newVulnDetector(x.infra.GitHub, x.infra.NewTrivyDB, x.config.TrivyDBPath)

	for req := range x.scanQueue {
		ctx := context.Background()

		clients := &scanClients{
			DB:        x.infra.DB,
			GitHubApp: x.infra.NewGitHubApp(x.config.GitHubAppID, req.InstallID, githubAppPEM),
			Detector:  detector,
			Utils:     x.infra.Utils,
		}

		if err := scanRepository(ctx, req, clients); err != nil {
			x.handleError(goerr.Wrap(err).With("request", req))
		}
	}

	return nil
}

type scanClients struct {
	DB        db.Interface
	GitHubApp githubapp.Interface
	Detector  *vulnDetector
	Utils     *infra.Utils
}

func insertScanReport(ctx context.Context, client db.Interface, req *model.ScanRepositoryRequest, pkgs []*ent.PackageRecord, vulnList []*ent.Vulnerability, now time.Time) error {
	if err := client.PutVulnerabilities(ctx, vulnList); err != nil {
		return err
	}

	addedPkgs, err := client.PutPackages(ctx, pkgs)
	if err != nil {
		return err
	}

	repo, err := client.CreateRepo(ctx, &ent.Repository{
		Owner: req.Owner,
		Name:  req.RepoName,
	})
	if err != nil {
		return err
	}

	report := &ent.Scan{
		Branch:      req.Branch,
		CommitID:    req.CommitID,
		RequestedAt: req.RequestedAt,
		ScannedAt:   now.Unix(),
	}
	if _, err := client.PutScan(ctx, report, repo, addedPkgs); err != nil {
		return err
	}

	return nil
}

func scanRepository(ctx context.Context, req *model.ScanRepositoryRequest, clients *scanClients) error {
	var checkID int64
	if req.IsPullRequest {
		id, err := clients.GitHubApp.CreateCheckRun(&req.GitHubRepo, req.CommitID)
		if err != nil {
			return err
		}
		checkID = id
	}

	if err := clients.Detector.RefreshDB(); err != nil {
		return err
	}

	pkgs, err := detectPackages(req, clients)
	if err != nil {
		return err
	}

	scannedAt := clients.Utils.Now()
	vulnList, err := annotateVulnerability(clients.Detector, pkgs, scannedAt.Unix())
	if err != nil {
		return err
	}

	// Retrieve latest scan report to compare with current one before inserting
	latest, err := clients.DB.GetLatestScan(ctx, req.GitHubBranch)
	if err != nil {
		return goerr.Wrap(err)
	}

	if err := insertScanReport(ctx, clients.DB, req, pkgs, vulnList, scannedAt); err != nil {
		return err
	}

	var oldPkgs []*ent.PackageRecord
	if latest != nil {
		oldPkgs = latest.Edges.Packages
	}

	sourcePkgMap := map[string][]*ent.PackageRecord{}
	var newPkgs []*ent.PackageRecord
	for i := range pkgs {
		sourcePkgMap[pkgs[i].Source] = append(sourcePkgMap[pkgs[i].Source], pkgs[i])
		newPkgs = append(newPkgs, pkgs[i])
	}

	addPkgs, modPkgs, delPkgs := diffPackageList(oldPkgs, newPkgs)

	if req.IsPullRequest {
		// TODO: feedback
		logger.Info().
			Interface("checkID", checkID).
			Interface("added", addPkgs).
			Interface("modified", modPkgs).
			Interface("deleted", delPkgs).
			Msg("DO Feedback")
	}

	return nil
}

func detectPackages(req *model.ScanRepositoryRequest, clients *scanClients) ([]*ent.PackageRecord, error) {
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

func annotateVulnerability(dt *vulnDetector, pkgs []*ent.PackageRecord, seenAt int64) ([]*ent.Vulnerability, error) {
	detectedVulnMap := map[string]*ent.Vulnerability{}
	sourcePkgMap := map[string][]*ent.PackageRecord{}

	for i := range pkgs {
		sourcePkgMap[pkgs[i].Source] = append(sourcePkgMap[pkgs[i].Source], pkgs[i])

		vulns, err := dt.Detect(pkgs[i].Type, pkgs[i].Name, pkgs[i].Version)
		if err != nil {
			return nil, err
		}

		vulnMap := map[string]struct{}{}
		for _, vuln := range vulns {
			vulnMap[vuln.VulnID] = struct{}{}
			if vuln.Detail.LastModifiedDate != nil {
				vuln.LastModifiedAt = vuln.Detail.LastModifiedDate.Unix()
			}

			detectedVulnMap[vuln.VulnID] = &ent.Vulnerability{
				ID:             vuln.VulnID,
				FirstSeenAt:    seenAt,
				LastModifiedAt: seenAt,
				Title:          vuln.Detail.Title,
				Description:    vuln.Detail.Description,
				Severity:       vuln.Detail.Severity,
				CweID:          vuln.Detail.CweIDs,
				References:     vuln.Detail.References,
			}
		}
		var vulnIDs []string
		for vulnID := range vulnMap {
			vulnIDs = append(vulnIDs, vulnID)
		}
		pkgs[i].VulnIds = vulnIDs
	}

	var vulnList []*ent.Vulnerability
	for _, v := range detectedVulnMap {
		vulnList = append(vulnList, v)
	}

	return vulnList, nil
}

func diffPackageList(oldPkgs, newPkgs []*ent.PackageRecord) (addPkgs, modPkgs, delPkgs []*ent.PackageRecord) {
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

func mapPackages(pkgs []*ent.PackageRecord) map[string]*ent.PackageRecord {
	resp := make(map[string]*ent.PackageRecord)
	for _, pkg := range pkgs {
		key := fmt.Sprintf("%s|%s|%s", pkg.Source, pkg.Name, pkg.Version)
		resp[key] = pkg
	}
	return resp
}

func matchVulnerabilities(a, b *ent.PackageRecord) bool {
	copyVulnList := func(p *ent.PackageRecord) []string {
		v := make([]string, len(p.VulnIds))
		for i := range p.VulnIds {
			v[i] = p.VulnIds[i]
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
