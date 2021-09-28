package usecase

import (
	"context"
	"io"
	"path/filepath"
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
			x.HandleError(err)
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
		logger.Debug().Interface("req", req).Msg("Recv scan request")
		ctx := context.Background()

		clients := &scanClients{
			DB:        x.infra.DB,
			GitHubApp: x.infra.NewGitHubApp(x.config.GitHubAppID, req.InstallID, githubAppPEM),
			Detector:  detector,
			Utils:     x.infra.Utils,

			FrontendURL: x.config.FrontendURL,
		}

		if err := scanRepository(ctx, req, clients); err != nil {
			x.HandleError(goerr.Wrap(err).With("request", req))
		}
	}

	return nil
}

type scanClients struct {
	DB        db.Interface
	GitHubApp githubapp.Interface
	Detector  *vulnDetector
	Utils     *infra.Utils

	FrontendURL string
}

func insertScanReport(ctx context.Context, client db.Interface, req *model.ScanRepositoryRequest, pkgs []*ent.PackageRecord, vulnList []*ent.Vulnerability, now time.Time) (*ent.Scan, error) {
	if err := client.PutVulnerabilities(ctx, vulnList); err != nil {
		return nil, err
	}

	addedPkgs, err := client.PutPackages(ctx, pkgs)
	if err != nil {
		return nil, err
	}

	repo, err := client.CreateRepo(ctx, &ent.Repository{
		Owner: req.Owner,
		Name:  req.RepoName,
	})
	if err != nil {
		return nil, err
	}

	report := &ent.Scan{
		Branch:      req.Branch,
		CommitID:    req.CommitID,
		RequestedAt: req.RequestedAt,
		ScannedAt:   now.Unix(),
	}
	added, err := client.PutScan(ctx, report, repo, addedPkgs)
	if err != nil {
		return nil, err
	}

	// Re-get scan result to retrieve all related records
	got, err := client.GetScan(ctx, added.ID)
	if err != nil {
		return nil, err
	}

	return got, nil
}

func scanRepository(ctx context.Context, req *model.ScanRepositoryRequest, clients *scanClients) error {
	var latest *ent.Scan
	// TargetBranch should be destination branch of PR (PR opened event) OR
	// PR branch (synchronized event).
	if req.TargetBranch != "" {
		// Retrieve latest scan report to compare with current one before inserting
		scan, err := clients.DB.GetLatestScan(ctx, model.GitHubBranch{
			GitHubRepo: req.GitHubRepo,
			Branch:     req.TargetBranch,
		})
		if err != nil {
			return goerr.Wrap(err)
		}
		latest = scan
		logger.Debug().Interface("scanID", scan.ID).Msg("Got latest scan")
	}

	/*
		// Disabled check run temporary
		check := newCheckRun(clients.GitHubApp)
		if err := check.create(&req.GitHubRepo, req.CommitID); err != nil {
			return err
		}
	*/
	if err := clients.Detector.RefreshDB(); err != nil {
		return err
	}

	newPkgs, err := crawlPackages(req, clients)
	if err != nil {
		return err
	}

	scannedAt := clients.Utils.Now()
	vulnList, err := annotateVulnerability(clients.Detector, newPkgs, scannedAt.Unix())
	if err != nil {
		return err
	}

	newScan, err := insertScanReport(ctx, clients.DB, req, newPkgs, vulnList, scannedAt)
	if err != nil {
		return err
	}
	logger.Debug().Str("scanID", newScan.ID).Msg("inserted scan report")

	if req.PullReqNumber != nil {
		status, err := clients.DB.GetVulnStatus(ctx, &req.GitHubRepo)
		if err != nil {
			return err
		}

		oldPkgs := []*ent.PackageRecord{}
		if latest != nil {
			oldPkgs = latest.Edges.Packages
		}

		changes := model.DiffVulnRecords(oldPkgs, newScan.Edges.Packages)
		db := model.NewVulnStatusDB(status, scannedAt.Unix())
		input := &postGitHubCommentInput{
			App:           clients.GitHubApp,
			Target:        &req.ScanTarget,
			Scan:          newScan,
			FrontendURL:   clients.FrontendURL,
			PullReqNumber: req.PullReqNumber,
			Report:        model.MakeReport(changes, db),
			GitHubEvent:   req.PullReqAction,
		}

		if err := postGitHubComment(input); err != nil {
			return err
		}
	}

	/*
		if err := check.complete(newScan.ID, changes, clients.FrontendURL); err != nil {
			return err
		}
	*/
	return nil
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
