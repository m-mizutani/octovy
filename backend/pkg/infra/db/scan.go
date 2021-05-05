package db

import (
	"fmt"
	"time"

	"github.com/guregu/dynamo"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/model"
)

const scanResultTimeKey = "2006-01-02T15:04:05"

func scanResultPK(branch *model.GitHubBranch) string {
	return fmt.Sprintf("scan:%s/%s@%s", branch.Owner, branch.RepoName, branch.Branch)
}
func scanResultSK(result *model.ScanResult) string {
	return fmt.Sprintf("%s/%s", time.Unix(result.ScannedAt, 0).Format(scanResultTimeKey), result.Target.CommitID)
}
func scanResultPK2(repo *model.GitHubRepo) string {
	return fmt.Sprintf("scan:%s/%s", repo.Owner, repo.RepoName)
}
func scanResultSK2(commitID string, scannedAt int64) string {
	return scanResultSK2Prefix(commitID) + time.Unix(scannedAt, 0).Format(scanResultTimeKey)
}
func scanResultSK2Prefix(commitID string) string {
	return commitID + "/"
}

func (x *DynamoClient) InsertScanResult(result *model.ScanResult) error {
	record := &dynamoRecord{
		PK:  scanResultPK(&result.Target.GitHubBranch),
		SK:  scanResultSK(result),
		PK2: scanResultPK2(&result.Target.GitHubRepo),
		SK2: scanResultSK2(result.Target.CommitID, result.ScannedAt),
		Doc: result,
	}

	if err := x.table.Put(record).Run(); err != nil {
		return goerr.Wrap(err).With("record", record)
	}

	return nil
}

func (x *DynamoClient) FindLatestScanResults(branch *model.GitHubBranch, n int) ([]*model.ScanResult, error) {
	var records []*dynamoRecord
	pk := scanResultPK(branch)
	if err := x.table.Get("pk", pk).Limit(int64(n)).Order(dynamo.Descending).All(&records); err != nil {
		return nil, goerr.Wrap(err).With("pk", pk)
	}

	scanResults := make([]*model.ScanResult, len(records))
	for i := range records {
		if err := records[i].Unmarshal(&scanResults[i]); err != nil {
			return nil, err
		}
	}
	return scanResults, nil
}

func (x *DynamoClient) FindScanResult(commit *model.GitHubCommit) (*model.ScanResult, error) {
	pk2 := scanResultPK2(&commit.GitHubRepo)
	sk2Prefix := scanResultSK2Prefix(commit.CommitID)

	var records []*dynamoRecord
	q := x.table.Get("pk2", pk2).Index(dynamoGSIName2nd).Range("sk2", dynamo.BeginsWith, sk2Prefix).Order(dynamo.Descending)
	if err := q.All(&records); err != nil {
		if !isNotFoundErr(err) {
			return nil, goerr.Wrap(err).With("pk2", pk2).With("sk2Prefix", sk2Prefix)
		}
		return nil, nil
	}

	var result model.ScanResult
	if err := records[0].Unmarshal(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
