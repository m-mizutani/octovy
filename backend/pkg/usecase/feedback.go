package usecase

import (
	"fmt"
	"math"
	"strings"
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
	var (
		err        error
		appID      int64
		privateKey []byte
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

	branch, err := x.LookupBranch(&model.GitHubBranch{
		GitHubRepo: report.Target.GitHubRepo,
		Branch:     req.Options.PullReqBranch,
	})
	if err != nil {
		return err
	}
	baseReport, err := x.LookupScanReport(branch.ReportSummary.ReportID)
	if err != nil {
		return err
	}

	body := buildFeedbackComment(report, baseReport)

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

		if err := app.CreateIssueComment(&report.Target.GitHubRepo, *req.Options.PullReqID, body); err != nil {
			return err
		}
	}
	return nil
}

type vulnRecord struct {
	Source     string
	VulnID     string
	PkgName    string
	PkgVersion string
}

func (x *vulnRecord) key() string {
	return strings.Join([]string{x.Source, x.VulnID, x.PkgName, x.PkgVersion}, "|")
}

func diffReport(newReport, oldReport *model.ScanReport) (news, fixed, remains []*vulnRecord) {

	reportToMap := func(report *model.ScanReport) map[string]*vulnRecord {
		if report == nil {
			return nil
		}

		m := map[string]*vulnRecord{}
		for _, src := range report.Sources {
			for _, pkg := range src.Packages {
				for _, vulnID := range pkg.Vulnerabilities {
					r := &vulnRecord{
						Source:     src.Source,
						VulnID:     vulnID,
						PkgName:    pkg.Name,
						PkgVersion: pkg.Version,
					}
					m[r.key()] = r
				}
			}
		}
		return m
	}

	newMap := reportToMap(newReport)
	oldMap := reportToMap(oldReport)

	// If no previous report
	if oldMap == nil {
		for _, n := range newMap {
			remains = append(remains, n)
		}
		return
	}

	// Compare with previous report
	for _, n := range newMap {
		if _, ok := oldMap[n.key()]; ok {
			remains = append(remains, n)
		} else {
			fixed = append(fixed, n)
		}
	}
	for _, o := range oldMap {
		if _, ok := newMap[o.key()]; !ok {
			news = append(news, o)
		}
	}

	return
}

func buildFeedbackComment(report, base *model.ScanReport) string {
	var body string
	const listSize = 5

	newVuln, fixedVuln, remainedVuln := diffReport(report, base)

	// New vulnerabilities
	if len(newVuln) > 0 {
		body += "### New vulnerabilities\n"
		for i := 0; i < len(newVuln) && i < listSize; i++ {
			body += fmt.Sprintf("- %s: `%s` %s in %s\n",
				newVuln[i].VulnID, newVuln[i].PkgName,
				newVuln[i].PkgVersion, newVuln[i].Source)
		}
		if len(newVuln) > listSize {
			body += fmt.Sprintf("... and more %d packages\n\n", len(newVuln)-listSize)
		}
		body += "\n"
	}

	// Fixed vulnerabilities
	if len(fixedVuln) > 0 {
		body += "### Fixed vulnerabilities\n"
		for i := 0; i < len(fixedVuln) && i < listSize; i++ {
			body += fmt.Sprintf("- %s: `%s` %s in %s\n",
				fixedVuln[i].VulnID, fixedVuln[i].PkgName,
				fixedVuln[i].PkgVersion, fixedVuln[i].Source)
		}
		if len(fixedVuln) > listSize {
			body += fmt.Sprintf("... and more %d packages\n\n", len(fixedVuln)-listSize)
		}
		body += "\n"
	}

	if len(remainedVuln) > 0 {
		remainCount := map[string]int{}
		for _, vuln := range remainedVuln {
			remainCount[vuln.Source] = remainCount[vuln.Source] + 1
		}

		body += `### Remained vulnerable packages\n`
		for src, count := range remainCount {
			body += fmt.Sprintf("- %d packages in %s\n", count, src)
		}
	}

	return body
}
