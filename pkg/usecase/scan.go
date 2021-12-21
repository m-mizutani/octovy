package usecase

import (
	"time"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/infra/opa"
	"github.com/m-mizutani/octovy/pkg/infra/policy"
	"github.com/m-mizutani/octovy/pkg/infra/trivy"
)

func (x *Usecase) SendScanRequest(req *model.ScanRepositoryRequest) error {
	x.scanQueue <- req
	return nil
}

func (x *Usecase) Scan(ctx *model.Context, req *model.ScanRepositoryRequest) error {
	ctx.With("scan_req", req)
	ctx.Log().Debug("recv scan request")

	clients := &scanClients{
		DB:          x.infra.DB,
		GitHubApp:   x.infra.NewGitHubApp(req.InstallID),
		Utils:       x.infra.Utils,
		Trivy:       x.infra.Trivy,
		CheckPolicy: x.infra.CheckPolicy,
		OPAClient:   x.infra.OPAClient,
		FrontendURL: x.config.FrontendURL,
	}

	if err := scanRepository(ctx, req, clients); err != nil {
		return err
	}

	return nil
}

type scanClients struct {
	DB          db.Interface
	GitHubApp   githubapp.Interface
	Trivy       trivy.Interface
	Utils       *infra.Utils
	CheckPolicy policy.Check
	OPAClient   opa.Interface

	FrontendURL string
}

func insertScan(ctx *model.Context, client db.Interface, req *model.ScanTarget, pkgs []*ent.PackageRecord, vulnList []*ent.Vulnerability, now time.Time) (*ent.Scan, error) {
	if err := client.PutVulnerabilities(ctx, vulnList); err != nil {
		return nil, err
	}

	addedPkgs, err := client.PutPackages(ctx, pkgs)
	if err != nil {
		return nil, err
	}

	repo, err := client.CreateRepo(ctx, &ent.Repository{
		Owner: req.Owner,
		Name:  req.Name,
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

func scanRepository(ctx *model.Context, req *model.ScanRepositoryRequest, clients *scanClients) error {
	var latest *ent.Scan
	// TargetBranch should be destination branch of PR (PR opened event) OR
	// PR branch (synchronized event).
	if req.TargetBranch != "" {
		// Retrieve latest scan report to compare with current one before inserting
		branch := model.GitHubBranch{
			GitHubRepo: req.GitHubRepo,
			Branch:     req.TargetBranch,
		}
		scan, err := clients.DB.GetLatestScan(ctx, branch)
		if err != nil {
			return goerr.Wrap(err)
		}
		latest = scan
		if scan != nil {
			ctx.Log().With("scanID", scan.ID).With("branch", branch).Debug("Got latest scan")
		} else {
			ctx.Log().With("branch", branch).Debug("No previous scan")
		}
	}

	check := newCheckRun(clients.GitHubApp)
	if clients.CheckPolicy != nil || clients.OPAClient != nil {
		if err := check.create(ctx, &req.GitHubRepo, req.CommitID); err != nil {
			return err
		}

		// Nothing happend if check completed properly
		defer check.fallback(ctx)
	}

	codes, err := setupGitHubCodes(ctx, req, clients.GitHubApp)
	if codes != nil {
		defer codes.RemoveAll()
	}
	if err != nil {
		return err
	}

	trivyResult, err := clients.Trivy.Scan(codes.Path)
	if err != nil {
		return err
	}

	// TODO: Merge insert scan procedure with PushTrivyResult
	scannedAt := clients.Utils.Now()
	newPkgs, vulnList := model.TrivyReportToEnt(trivyResult, scannedAt)

	newScan, err := insertScan(ctx, clients.DB, &req.ScanTarget, newPkgs, vulnList, scannedAt)
	if err != nil {
		return err
	}
	ctx.Log().With("scanID", newScan.ID).Debug("inserted scan report")

	status, err := clients.DB.GetVulnStatus(ctx, &req.GitHubRepo)
	if err != nil {
		return err
	}

	oldPkgs := []*ent.PackageRecord{}
	if latest != nil {
		oldPkgs = latest.Edges.Packages
	}

	now := scannedAt.Unix()
	changes := model.DiffVulnRecords(oldPkgs, newScan.Edges.Packages)
	db := model.NewVulnStatusDB(status, now)
	report := model.MakeReport(newScan.ID, changes, db, clients.FrontendURL)

	if req.PullReqNumber != nil {
		input := &postGitHubCommentInput{
			App:           clients.GitHubApp,
			Target:        &req.ScanTarget,
			PullReqNumber: req.PullReqNumber,
			Report:        report,
			GitHubEvent:   req.PullReqAction,
		}

		if err := postGitHubComment(input); err != nil {
			return err
		}
	}

	scanReport := model.NewScanReport(newScan, status, now)
	var result *model.GitHubCheckResult

	if clients.OPAClient != nil {
		var r model.GitHubCheckResult
		if err := clients.OPAClient.Data(ctx, opa.Check, scanReport, &r); err != nil {
			return err
		}
		result = &r
	} else if clients.CheckPolicy != nil {
		r, err := clients.CheckPolicy.Result(ctx, scanReport)
		if err != nil {
			return err
		}
		result = r
	}

	if result != nil {
		if err := check.complete(ctx, newScan.ID, report, clients.FrontendURL, result); err != nil {
			return err
		}
	}

	return nil
}
