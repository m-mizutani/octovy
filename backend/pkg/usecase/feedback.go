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

type feedbackProps struct {
	DB          interfaces.DBClient
	App         interfaces.GitHubApp
	NewReport   *model.ScanReport
	OldReport   *model.ScanReport
	Options     *model.FeedbackOptions
	FrontendURL string
	CheckFail   bool
}

func (x *Default) FeedbackScanResult(req *model.FeedbackRequest) error {
	logger.With("req", req).Info("Recv request")

	newReport, err := getScanReport(x.svc.DB(), req.ReportID)
	if err != nil {
		return goerr.Wrap(err).With("req", req)
	}

	oldReport, err := getOldReport(x.svc.DB(), &newReport.Target.GitHubRepo, req.Options.PullReqBranch)
	if err != nil {
		return err
	}

	app, err := x.buildGitHubApp(req.InstallID)
	if err != nil {
		return err
	}

	props := feedbackProps{
		DB:          x.svc.DB(),
		App:         app,
		NewReport:   newReport,
		OldReport:   oldReport,
		Options:     &req.Options,
		FrontendURL: x.config.FrontendBaseURL(),
		CheckFail:   x.config.ShouldFailIfVuln(),
	}

	if err := feedbackPullRequest(props); err != nil {
		return err
	}
	if err := feedbackCheckRun(props); err != nil {
		return err
	}
	return nil
}

func getScanReport(db interfaces.DBClient, reportID string) (*model.ScanReport, error) {
	const (
		waitFactor = 1.2
		maxRetry   = 10
	)

	var report *model.ScanReport
	for i := 0; i < maxRetry; i++ {
		r, err := db.LookupScanReport(reportID)
		if err != nil {
			return nil, err
		}

		if report = r; report != nil {
			break
		}
		w := math.Pow(waitFactor, float64(i))
		time.Sleep(time.Millisecond * time.Duration(w*1000))
	}
	if report == nil {
		return nil, goerr.New("Report is not found")
	}

	return report, nil
}

func getOldReport(db interfaces.DBClient, repo *model.GitHubRepo, branch string) (*model.ScanReport, error) {
	// Destination branch of merge
	if branch != "" {
		branch, err := db.LookupBranch(&model.GitHubBranch{
			GitHubRepo: *repo,
			Branch:     branch,
		})
		if err != nil {
			return nil, err
		}
		if branch != nil && branch.ReportSummary.ReportID != "" {
			r, err := db.LookupScanReport(branch.ReportSummary.ReportID)
			if err != nil {
				return nil, err
			}

			return r, nil
		}
	}

	return nil, nil
}

func feedbackPullRequest(props feedbackProps) error {
	if props.Options.PullReqID == nil {
		return nil
	}

	var pullReqReport *model.ScanReport
	lastScanLogs, err := props.DB.FindScanLogsByBranch(&props.NewReport.Target.GitHubBranch, 2)
	if err != nil {
		return err
	}
	if len(lastScanLogs) == 2 {
		lastScanReport, err := props.DB.LookupScanReport(lastScanLogs[1].Summary.ReportID)
		if err != nil {
			return err
		}
		pullReqReport = lastScanReport
	}
	if pullReqReport == nil {
		pullReqReport = props.OldReport
	}

	body := buildFeedbackComment(props.NewReport, pullReqReport, props.FrontendURL, false)
	if body == "" {
		return nil
	}

	logger.With("props", props).With("report", props.NewReport).Info("Creating a PR comment")

	if err := props.App.CreateIssueComment(&props.NewReport.Target.GitHubRepo, *props.Options.PullReqID, body); err != nil {
		return err
	}

	return nil
}

func feedbackCheckRun(props feedbackProps) error {
	if props.Options.CheckID == nil {
		return nil
	}

	logger.With("req", props.Options).With("report", props.NewReport).Info("Creating a PR comment")

	changes := diffReport(props.NewReport, props.OldReport)

	// Default messages
	conclusion := "neutral"
	title := fmt.Sprintf("❗ %d vulnerabilities detected", len(changes.Unfixed)+len(changes.News))
	summary := fmt.Sprintf("New %d and remained %d vulnerabilities found", len(changes.News), len(changes.Unfixed))
	body := buildFeedbackComment(props.NewReport, props.OldReport, props.FrontendURL, true)

	if len(changes.Unfixed) == 0 && len(changes.News) == 0 {
		conclusion = "success"
		title = "🎉  No vulnerability detected"
		summary = "OK"
	} else if props.CheckFail {
		conclusion = "failure"
	}

	opt := &github.UpdateCheckRunOptions{
		Name:        "Octovy: package vulnerability check",
		Status:      github.String("completed"),
		CompletedAt: &github.Timestamp{Time: time.Unix(props.NewReport.ScannedAt, 0)},
		Conclusion:  &conclusion,
		DetailsURL:  github.String(props.FrontendURL + "/#/scan/report/" + props.NewReport.ReportID),
		Output: &github.CheckRunOutput{
			Title:   &title,
			Summary: &summary,
			Text:    &body,
		},
	}

	if err := props.App.UpdateCheckRun(&props.NewReport.Target.GitHubRepo, *props.Options.CheckID, opt); err != nil {
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
			res.News = append(res.News, n)
		}
	}
	for _, o := range oldMap {
		if _, ok := newMap[o.key()]; !ok {
			res.Fixed = append(res.Fixed, o)
		}
	}
	return
}

func feedbackCommentVulnRecord(v *vulnRecord, url string) string {
	return fmt.Sprintf("- [%s](%s/#/vuln/%s): `%s` %s in [%s](%s)\n",
		v.VulnID, url, v.VulnID, v.PkgName, v.PkgVersion, v.Source, v.Source)
}

func buildFeedbackComment(report, base *model.ScanReport, frontendURL string, showUnfix bool) string {
	var body string
	const listSize = 5

	changes := diffReport(report, base)
	if len(changes.News) == 0 && len(changes.Unfixed) == 0 {
		body += "🎉 **No vulnerable packages**\n\n"
	}
	if len(changes.News) == 0 && len(changes.Fixed) == 0 && !showUnfix {
		return ""
	}

	// New vulnerabilities
	if len(changes.News) > 0 {
		body += "### 🚨 New vulnerabilities\n"
		for i := 0; i < len(changes.News) && i < listSize; i++ {
			v := changes.News[i]
			body += feedbackCommentVulnRecord(v, frontendURL)
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
			body += feedbackCommentVulnRecord(v, frontendURL)
		}
		if len(changes.Fixed) > listSize {
			body += fmt.Sprintf("... and more %d packages\n\n", len(changes.Fixed)-listSize)
		}
		body += "\n"
	}

	if showUnfix && len(changes.Unfixed) > 0 {
		remainCount := map[string]int{}
		for _, vuln := range changes.Unfixed {
			remainCount[vuln.Source] = remainCount[vuln.Source] + 1
		}

		body += "### ⚠️ Unfixed vulnerable packages\n"
		for src, count := range remainCount {
			body += fmt.Sprintf("- %d packages in %s\n", count, src)
		}
	}

	body += fmt.Sprintf("\nSee [report](%s/#/scan/report/%s) for more detail\n", frontendURL, report.ReportID)

	return body
}
