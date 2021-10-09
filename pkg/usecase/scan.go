package usecase

import (
	"context"
	"time"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/infra/trivy"
)

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
	for req := range x.scanQueue {
		logger.Debug().Interface("req", req).Msg("Recv scan request")
		ctx := context.Background()

		clients := &scanClients{
			DB:          x.infra.DB,
			GitHubApp:   x.infra.NewGitHubApp(x.config.GitHubAppID, req.InstallID, []byte(x.config.GitHubAppPrivateKey)),
			Utils:       x.infra.Utils,
			Trivy:       x.infra.Trivy,
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
	Trivy     trivy.Interface
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
			logger.Debug().Interface("scanID", scan.ID).Interface("branch", branch).Msg("Got latest scan")
		} else {
			logger.Debug().Interface("branch", branch).Msg("No previous scan")
		}
	}

	/*
		// Disabled check run temporary
		check := newCheckRun(clients.GitHubApp)
		if err := check.create(&req.GitHubRepo, req.CommitID); err != nil {
			return err
		}
	*/

	codes, err := setupGitHubCodes(ctx, req, clients.GitHubApp)
	if codes != nil {
		defer codes.RemoveAll()
	}
	if err != nil {
		return err
	}

	report, err := clients.Trivy.Scan(codes.Path)
	if err != nil {
		return err
	}

	scannedAt := clients.Utils.Now()
	newPkgs, vulnList := model.TrivyReportToEnt(report, scannedAt)

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
