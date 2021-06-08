package db

import (
	"fmt"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

func vulnResponsePK(repo *model.GitHubRepo) string {
	return fmt.Sprintf("response:%s/%s", repo.Owner, repo.RepoName)
}
func vulnResponseSK(pkgType model.PkgType, pkgName string, vulnID string) string {
	return fmt.Sprintf("%s|%s|%s", pkgType, pkgName, vulnID)
}

func (x *DynamoClient) PutVulnResponse(resp *model.VulnResponse) error {
	if err := resp.IsValid(); err != nil {
		return err
	}

	record := &dynamoRecord{
		PK:  vulnResponsePK(&resp.GitHubRepo),
		SK:  vulnResponseSK(resp.PkgType, resp.PkgName, resp.VulnID),
		Doc: resp,
	}
	if resp.Duration > 0 {
		record.ExpiresAt = model.Int64(resp.CreatedAt + resp.Duration)
	}

	if err := x.table.Put(record).Run(); err != nil {
		return err
	}

	return nil
}

func (x *DynamoClient) GetVulnResponses(repo *model.GitHubRepo, now int64) ([]*model.VulnResponse, error) {
	var records []*dynamoRecord
	pk := vulnResponsePK(repo)
	if err := x.table.Get("pk", pk).Filter("attribute_not_exists(expires_at) OR ? < expires_at", now).All(&records); err != nil {
		return nil, goerr.Wrap(err).With("repo", repo)
	}

	vulnResponses := make([]*model.VulnResponse, len(records))
	for i := range records {
		if err := records[i].Unmarshal(&vulnResponses[i]); err != nil {
			return nil, err
		}
	}
	return vulnResponses, nil
}

func (x *DynamoClient) DeleteVulnResponse(resp *model.VulnResponse) error {
	if err := resp.IsValid(); err != nil {
		return err
	}

	pk := vulnResponsePK(&resp.GitHubRepo)
	sk := vulnResponseSK(resp.PkgType, resp.PkgName, resp.VulnID)

	if err := x.table.Delete("pk", pk).Range("sk", sk).Run(); err != nil {
		return err
	}

	return nil
}
