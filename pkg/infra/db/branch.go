package db

import (
	"fmt"
	"time"

	"github.com/guregu/dynamo"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

const branchTimeKey = "2006-01-02T15:04:05"

func branchPK(repo *model.GitHubRepo) string {
	return fmt.Sprintf("branch:%s/%s", repo.Owner, repo.RepoName)
}
func branchSK(branchName string) string {
	return branchName
}
func branchPK2(repo *model.GitHubRepo) string {
	return fmt.Sprintf("branch:%s/%s", repo.Owner, repo.RepoName)
}
func branchSK2(branch *model.Branch) string {
	return fmt.Sprintf("%s/%s", time.Unix(branch.LastScannedAt, 0).Format(branchTimeKey), branch.Branch)
}

func (x *DynamoClient) UpdateBranch(branch *model.Branch) error {
	record := &dynamoRecord{
		PK:  branchPK(&branch.GitHubRepo),
		SK:  branchSK(branch.Branch),
		PK2: branchPK2(&branch.GitHubRepo),
		SK2: branchSK2(branch),
		Doc: branch,
	}

	q := x.table.Put(record).If("(attribute_not_exists(pk) AND attribute_not_exists(sk)) OR doc.'LastScannedAt' < ?", branch.LastScannedAt)
	if err := q.Run(); err != nil {
		if !isConditionalCheckErr(err) {
			return goerr.Wrap(err).With("branch", branch)
		}
	}

	return nil
}

func (x *DynamoClient) LookupBranch(branch *model.GitHubBranch) (*model.Branch, error) {
	pk := branchPK(&branch.GitHubRepo)
	sk := branchSK(branch.Branch)

	var record *dynamoRecord

	if err := x.table.Get("pk", pk).Range("sk", dynamo.Equal, sk).One(&record); err != nil {
		if !isNotFoundErr(err) {
			return nil, goerr.Wrap(err).With("branch", branch)
		}
		return nil, nil
	}

	var resp model.Branch
	if err := record.Unmarshal(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (x *DynamoClient) FindLatestScannedBranch(repo *model.GitHubRepo, n int) ([]*model.Branch, error) {
	pk2 := branchPK2(repo)
	var records []*dynamoRecord
	if err := x.table.Get("pk2", pk2).Index(dynamoGSIName2nd).Limit(int64(n)).Order(dynamo.Descending).All(&records); err != nil {
		return nil, goerr.Wrap(err).With("pk2", pk2)
	}

	branches := make([]*model.Branch, len(records))
	for i := range records {
		if err := records[i].Unmarshal(&branches[i]); err != nil {
			return nil, err
		}
	}
	return branches, nil
}
