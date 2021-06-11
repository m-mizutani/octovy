package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/dynamo"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

func vulnStatusPK(repo *model.GitHubRepo) string {
	return fmt.Sprintf("vuln_status:%s/%s", repo.Owner, repo.RepoName)
}
func vulnStatusSK(key *model.VulnPackageKey) string {
	return fmt.Sprintf("%s|%s|%s|%s", key.Source, key.PkgType, key.PkgName, key.VulnID)
}

const vulnStatusLogTimeKey = "2006-01-02T15:04:05"

func vulnStatusLogPK(repo *model.GitHubRepo) string {
	return fmt.Sprintf("vuln_status_log:%s/%s", repo.Owner, repo.RepoName)
}
func vulnStatusLogSK(key *model.VulnPackageKey, createdAt int64) string {
	ts := time.Unix(createdAt, 0)
	return fmt.Sprintf("%s%s|%s", vulnStatusLogSKPrefix(key),
		ts.Format(vulnStatusLogTimeKey), uuid.New().String())
}
func vulnStatusLogSKPrefix(key *model.VulnPackageKey) string {
	return vulnStatusSK(key) + "|"
}

func (x *DynamoClient) PutVulnStatus(status *model.VulnStatus) error {
	if err := status.IsValid(); err != nil {
		return err
	}

	tx := x.db.WriteTx()

	record := &dynamoRecord{
		PK:  vulnStatusPK(&status.GitHubRepo),
		SK:  vulnStatusSK(&status.VulnPackageKey),
		Doc: status,
	}
	if status.ExpiresAt > 0 {
		record.ExpiresAt = &status.ExpiresAt
	}
	tx = tx.Put(x.table.Put(record))

	logRecord := &dynamoRecord{
		PK:  vulnStatusLogPK(&status.GitHubRepo),
		SK:  vulnStatusLogSK(&status.VulnPackageKey, status.CreatedAt),
		Doc: status,
	}
	tx = tx.Put(x.table.Put(logRecord))

	if err := tx.Run(); err != nil {
		return goerr.Wrap(err).With("record", record).With("log", logRecord)
	}

	return nil
}

func (x *DynamoClient) GetVulnStatus(repo *model.GitHubRepo, now int64) ([]*model.VulnStatus, error) {
	var records []*dynamoRecord
	pk := vulnStatusPK(repo)
	if err := x.table.Get("pk", pk).Filter("attribute_not_exists(expires_at) OR ? < expires_at", now).All(&records); err != nil {
		return nil, goerr.Wrap(err).With("repo", repo)
	}

	vulnStatuss := make([]*model.VulnStatus, len(records))
	for i := range records {
		if err := records[i].Unmarshal(&vulnStatuss[i]); err != nil {
			return nil, err
		}
	}
	return vulnStatuss, nil
}

func (x *DynamoClient) GetVulnStatusLogs(repo *model.GitHubRepo, key *model.VulnPackageKey) ([]*model.VulnStatus, error) {
	pk := vulnStatusLogPK(repo)
	skPrefix := vulnStatusLogSKPrefix(key)

	var records []*dynamoRecord

	if err := x.table.Get("pk", pk).Range("sk", dynamo.BeginsWith, skPrefix).All(&records); err != nil {
		return nil, err
	}

	vulnStatuss := make([]*model.VulnStatus, len(records))
	for i := range records {
		if err := records[i].Unmarshal(&vulnStatuss[i]); err != nil {
			return nil, err
		}
	}
	return vulnStatuss, nil

}
