package usecase

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
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

	var baseReport *model.ScanReport
	if req.Options.PullReqBranch != "" {
		branch, err := x.LookupBranch(&model.GitHubBranch{
			GitHubRepo: report.Target.GitHubRepo,
			Branch:     req.Options.PullReqBranch,
		})
		if err != nil {
			return err
		}
		r, err := x.LookupScanReport(branch.ReportSummary.ReportID)
		if err != nil {
			return err
		}
		baseReport = r
	}

	app, err := x.buildGitHubApp(req.InstallID)
	if err != nil {
		return err
	}

	if err := feedbackPullRequest(app, &req.Options, report, baseReport, x.config.FrontendURL); err != nil {
		return err
	}
	if err := feedbackCheckRun(app, &req.Options, report, baseReport, x.config.FrontendURL); err != nil {
		return err
	}
	return nil
}

func feedbackPullRequest(app interfaces.GitHubApp, feedback *model.FeedbackOptions, newReport, oldReport *model.ScanReport, frontendURL string) error {
	if feedback.PullReqID == nil {
		return nil
	}

	body := buildFeedbackComment(newReport, oldReport)

	logger.With("req", feedback).With("report", newReport).Info("Creating a PR comment")

	if err := app.CreateIssueComment(&newReport.Target.GitHubRepo, *feedback.PullReqID, body); err != nil {
		return err
	}

	return nil
}

func feedbackCheckRun(app interfaces.GitHubApp, feedback *model.FeedbackOptions, newReport, oldReport *model.ScanReport, frontendURL string) error {
	if feedback.CheckID == nil {
		return nil
	}

	logger.With("req", feedback).With("report", newReport).Info("Creating a PR comment")

	changes := diffReport(newReport, oldReport)

	// Default messages
	conclusion := "neutral"
	title := fmt.Sprintf("❗ %d vulnerabilities detected", len(changes.Unfixed)+len(changes.News))
	summary := fmt.Sprintf("New %d and remained %d vulnerabilities found", len(changes.News), len(changes.Unfixed))
	body := buildFeedbackComment(newReport, oldReport)

	if len(changes.Unfixed) == 0 && len(changes.News) == 0 {
		conclusion = "success"
		title = "🎉  No vulnerability detected"
		summary = "OK"
	}

	opt := &github.UpdateCheckRunOptions{
		Status:      github.String("completed"),
		CompletedAt: &github.Timestamp{Time: time.Unix(newReport.ScannedAt, 0)},
		Conclusion:  &conclusion,
		DetailsURL:  github.String(frontendURL + "/#/scan/report/" + newReport.ReportID),
		Output: &github.CheckRunOutput{
			Title:   &title,
			Summary: &summary,
			Text:    &body,
		},
	}

	if err := app.UpdateCheckRun(&newReport.Target.GitHubRepo, *feedback.CheckID, opt); err != nil {
		return err
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

type changeResult struct {
	News    []*vulnRecord
	Fixed   []*vulnRecord
	Unfixed []*vulnRecord
}

func diffReport(newReport, oldReport *model.ScanReport) (res changeResult) {
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
			res.Unfixed = append(res.Unfixed, n)
		}
		return
	}

	// Compare with previous report
	for _, n := range newMap {
		if _, ok := oldMap[n.key()]; ok {
			res.Unfixed = append(res.Unfixed, n)
		} else {
			res.Fixed = append(res.Fixed, n)
		}
	}
	for _, o := range oldMap {
		if _, ok := newMap[o.key()]; !ok {
			res.News = append(res.News, o)
		}
	}
	return
}

func buildFeedbackComment(report, base *model.ScanReport) string {
	var body string
	const listSize = 5

	changes := diffReport(report, base)
	if len(changes.News) == 0 && len(changes.Unfixed) == 0 {
		body += "🎉 **No vulnerable packages**\n\n"
	}

	// New vulnerabilities
	if len(changes.News) > 0 {
		body += "### 🚨 New vulnerabilities\n"
		for i := 0; i < len(changes.News) && i < listSize; i++ {
			v := changes.News[i]
			body += fmt.Sprintf("- %s: `%s` %s in %s\n",
				v.VulnID, v.PkgName, v.PkgVersion, v.Source)
		}
		if len(changes.News) > listSize {
			body += fmt.Sprintf("... and more %d packages\n\n", len(changes.News)-listSize)
		}
		body += "\n"
	}

	// Fixed vulnerabilities
	if len(changes.Fixed) > 0 {
		body += "### ✅ Fixed vulnerabilities\n"
		for i := 0; i < len(changes.Fixed) && i < listSize; i++ {
			v := changes.Fixed[i]
			body += fmt.Sprintf("- %s: `%s` %s in %s\n",
				v.VulnID, v.PkgName, v.PkgVersion, v.Source)
		}
		if len(changes.Fixed) > listSize {
			body += fmt.Sprintf("... and more %d packages\n\n", len(changes.Fixed)-listSize)
		}
		body += "\n"
	}

	if len(changes.Unfixed) > 0 {
		remainCount := map[string]int{}
		for _, vuln := range changes.Unfixed {
			remainCount[vuln.Source] = remainCount[vuln.Source] + 1
		}

		body += "### ⚠️ Unfixed vulnerable packages\n"
		for src, count := range remainCount {
			body += fmt.Sprintf("- %d packages in %s\n", count, src)
		}
	}

	return body
}
