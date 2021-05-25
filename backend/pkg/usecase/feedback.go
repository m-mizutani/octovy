package usecase

import (
	"math"
	"time"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

func (x *Default) FeedbackScanResult(req *model.FeedbackRequest) error {
	logger.With("req", req).Info("Recv request")

	const (
		waitFactor = 1.2
		maxRetry   = 10
	)

	var report *model.ScanReport
	for i := 0; i < maxRetry; i++ {
		r, err := x.LookupScanReport(req.ReportID)
		if err != nil {
			return err
		}

		if report = r; report != nil {
			break
		}
		w := math.Pow(waitFactor, float64(i))
		time.Sleep(time.Millisecond * time.Duration(w*1000))
	}
	if report == nil {
		return goerr.New("Report is not found").With("req", req)
	}

	var (
		err        error
		appID      int64
		privateKey []byte
	)

	secrets, err := x.svc.GetSecrets()
	if err != nil {
		return err
	}
	if appID, err = secrets.GetGitHubAppID(); err != nil {
		return err
	}
	if privateKey, err = secrets.GithubAppPEM(); err != nil {
		return err
	}

	app := x.svc.Infra.NewGitHubApp(appID, req.InstallID, privateKey, x.config.GitHubEndpoint)

	if req.Options.PullReqID != nil {
		logger.With("req", req).With("report", report).Info("Creating a PR comment")
		if err := app.CreateIssueComment(&report.Target.GitHubRepo, *req.Options.PullReqID, report.ReportID); err != nil {
			return err
		}
	}

	return nil
}
