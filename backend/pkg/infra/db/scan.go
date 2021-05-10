package db

import (
	"fmt"
	"time"

	"github.com/guregu/dynamo"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/model"
)

const scanLogTimeKey = "2006-01-02T15:04:05"

func scanLogPK(branch *model.GitHubBranch) string {
	return fmt.Sprintf("scan_log:%s/%s@%s", branch.Owner, branch.RepoName, branch.Branch)
}
func scanLogSK(log *model.ScanLog) string {
	return fmt.Sprintf("%s/%s", time.Unix(log.ScannedAt, 0).Format(scanLogTimeKey), log.Target.CommitID)
}
func scanLogPK2(repo *model.GitHubRepo) string {
	return fmt.Sprintf("scan_log:%s/%s", repo.Owner, repo.RepoName)
}
func scanLogSK2(commitID string, scannedAt int64) string {
	return scanLogSK2Prefix(commitID) + time.Unix(scannedAt, 0).Format(scanLogTimeKey)
}
func scanLogSK2Prefix(commitID string) string {
	return commitID + "/"
}

func scanReportPK(reportID string) string {
	return `scan_report:` + reportID
}

func scanReportSK() string {
	return "*"
}

func (x *DynamoClient) InsertScanReport(report *model.ScanReport) error {
	if err := report.IsValid(); err != nil {
		return err
	}

	tx := x.db.WriteTx()

	scanLog := report.ToLog()
	logRecord := &dynamoRecord{
		PK:  scanLogPK(&scanLog.Target.GitHubBranch),
		SK:  scanLogSK(scanLog),
		PK2: scanLogPK2(&scanLog.Target.GitHubRepo),
		SK2: scanLogSK2(scanLog.Target.CommitID, scanLog.ScannedAt),
		Doc: scanLog,
	}
	tx = tx.Put(x.table.Put(logRecord))

	reportRecord := &dynamoRecord{
		PK:  scanReportPK(report.ReportID),
		SK:  scanReportSK(),
		Doc: report,
	}
	tx = tx.Put(x.table.Put(reportRecord))

	if err := tx.Run(); err != nil {
		return goerr.Wrap(err).With("log", logRecord).With("report", reportRecord)
	}

	return nil
}

func (x *DynamoClient) LookupScanReport(reportID string) (*model.ScanReport, error) {
	pk := scanReportPK(reportID)
	sk := scanReportSK()

	var record *dynamoRecord

	if err := x.table.Get("pk", pk).Range("sk", dynamo.Equal, sk).One(&record); err != nil {
		if !isNotFoundErr(err) {
			return nil, goerr.Wrap(err).With("ReportID", reportID)
		}
		return nil, nil
	}

	var report model.ScanReport
	if err := record.Unmarshal(&report); err != nil {
		return nil, err
	}
	return &report, nil
}

func (x *DynamoClient) FindScanLogsByBranch(branch *model.GitHubBranch, n int) ([]*model.ScanLog, error) {
	var records []*dynamoRecord
	pk := scanLogPK(branch)
	if err := x.table.Get("pk", pk).Limit(int64(n)).Order(dynamo.Descending).All(&records); err != nil {
		return nil, goerr.Wrap(err).With("pk", pk)
	}

	scanLogs := make([]*model.ScanLog, len(records))
	for i := range records {
		if err := records[i].Unmarshal(&scanLogs[i]); err != nil {
			return nil, err
		}
	}
	return scanLogs, nil
}

func (x *DynamoClient) FindScanLogsByCommit(commit *model.GitHubCommit, n int) ([]*model.ScanLog, error) {
	var records []*dynamoRecord
	pk := scanLogPK2(&commit.GitHubRepo)
	sk := scanLogSK2Prefix(commit.CommitID)

	if err := x.table.Get("pk2", pk).Index(dynamoGSIName2nd).Range("sk2", dynamo.BeginsWith, sk).Limit(int64(n)).Order(dynamo.Descending).All(&records); err != nil {
		return nil, goerr.Wrap(err).With("pk", pk)
	}

	scanLogs := make([]*model.ScanLog, len(records))
	for i := range records {
		if err := records[i].Unmarshal(&scanLogs[i]); err != nil {
			return nil, err
		}
	}
	return scanLogs, nil
}
