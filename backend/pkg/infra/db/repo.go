package db

import (
	"github.com/guregu/dynamo"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/model"
)

func repositoryPK() string {
	return "list:repository"
}
func repositorySK(owner, name string) string {
	return owner + "/" + name
}

func (x *DynamoClient) InsertRepo(repo *model.Repository) (bool, error) {
	record := &dynamoRecord{
		PK:  repositoryPK(),
		SK:  repositorySK(repo.Owner, repo.RepoName),
		Doc: repo,
	}
	put := x.table.Put(record).If("attribute_not_exists(pk) AND attribute_not_exists(sk)")
	if err := put.Run(); err != nil {
		if isConditionalCheckErr(err) {
			return false, nil
		}
		return false, goerr.Wrap(err).With("record", record)
	}

	return true, nil
}

func (x *DynamoClient) SetRepoBranches(repo *model.GitHubRepo, branches []string) error {
	pk := repositoryPK()
	sk := repositorySK(repo.Owner, repo.RepoName)
	update := x.table.Update("pk", pk).Range("sk", sk).
		Set("doc.'Branches'", branches)
	if err := update.Run(); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}

func (x *DynamoClient) SetRepoDefaultBranch(repo *model.GitHubRepo, branch string) error {
	pk := repositoryPK()
	sk := repositorySK(repo.Owner, repo.RepoName)
	update := x.table.Update("pk", pk).Range("sk", sk).
		Set("doc.'DefaultBranch'", branch)
	if err := update.Run(); err != nil {
		return goerr.Wrap(err)
	}

	setBranch := x.table.Update("pk", pk).Range("sk", sk).
		Set("doc.'Branches'", []string{branch}).
		If("doc.'Branches'.length = 0")
	if err := setBranch.Run(); err != nil && isConditionalCheckErr(err) {
		return goerr.Wrap(err).With("repo", repo)
	}

	return nil
}

func recordToRepo(records []*dynamoRecord) ([]*model.Repository, error) {
	repositories := make([]*model.Repository, len(records))
	for i := range records {
		if err := records[i].Unmarshal(&repositories[i]); err != nil {
			return nil, err
		}
	}
	return repositories, nil
}

func (x *DynamoClient) FindRepo() ([]*model.Repository, error) {
	var records []*dynamoRecord
	pk := repositoryPK()
	if err := x.table.Get("pk", pk).All(&records); err != nil {
		if !isNotFoundErr(err) {
			return nil, goerr.Wrap(err)
		}
	}

	return recordToRepo(records)
}

func (x *DynamoClient) FindRepoByOwner(owner string) ([]*model.Repository, error) {
	var records []*dynamoRecord
	pk := repositoryPK()
	sk := repositorySK(owner, "")
	if err := x.table.Get("pk", pk).Range("sk", dynamo.BeginsWith, sk).All(&records); err != nil {
		if !isNotFoundErr(err) {
			return nil, goerr.Wrap(err)
		}
	}

	return recordToRepo(records)
}

func (x *DynamoClient) FindRepoByFullName(owner, name string) (*model.Repository, error) {
	var record *dynamoRecord
	pk := repositoryPK()
	sk := repositorySK(owner, name)
	if err := x.table.Get("pk", pk).Range("sk", dynamo.Equal, sk).One(&record); err != nil {
		if !isNotFoundErr(err) {
			return nil, goerr.Wrap(err)
		}
	}

	var repo model.Repository
	if err := record.Unmarshal(&repo); err != nil {
		return nil, err
	}

	return &repo, nil
}
