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

func ownerPK() string {
	return "list:owner"
}
func ownerSK(owner string) string {
	return owner
}

func (x *DynamoClient) InsertRepo(repo *model.Repository) (bool, error) {
	ownerRecord := &dynamoRecord{
		PK:  ownerPK(),
		SK:  ownerSK(repo.Owner),
		Doc: model.Owner{Name: repo.Owner},
	}
	if err := x.table.Put(ownerRecord).Run(); err != nil {
		if isConditionalCheckErr(err) {
			return false, nil
		}
		return false, goerr.Wrap(err).With("ownerRecord", ownerRecord)
	}

	repoRecord := &dynamoRecord{
		PK:  repositoryPK(),
		SK:  repositorySK(repo.Owner, repo.RepoName),
		Doc: repo,
	}
	put := x.table.Put(repoRecord).If("attribute_not_exists(pk) AND attribute_not_exists(sk)")
	if err := put.Run(); err != nil {
		if isConditionalCheckErr(err) {
			return false, nil
		}
		return false, goerr.Wrap(err).With("repoRecord", repoRecord)
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

func (x *DynamoClient) UpdateBranchIfDefault(repo *model.GitHubRepo, branch *model.Branch) error {
	pk := repositoryPK()
	sk := repositorySK(repo.Owner, repo.RepoName)
	q := x.table.Update("pk", pk).Range("sk", sk).
		Set("doc.'Branch'", branch).
		If("doc.'DefaultBranch' = ?", branch.Branch).
		If("doc.'Branch'.'LastScannedAt' < ?", branch.LastScannedAt)
	if err := q.Run(); err != nil {
		if isConditionalCheckErr(err) {
			return nil
		}
		return goerr.Wrap(err).With("repo", repo).With("branch", branch)
	}
	return nil
}

func (x *DynamoClient) SetRepoDefaultBranchName(repo *model.GitHubRepo, branch string) error {
	pk := repositoryPK()
	sk := repositorySK(repo.Owner, repo.RepoName)
	update := x.table.Update("pk", pk).Range("sk", sk).
		Set("doc.'DefaultBranch'", branch)
	if err := update.Run(); err != nil {
		return goerr.Wrap(err)
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

func (x *DynamoClient) FindOwners() ([]*model.Owner, error) {
	var records []*dynamoRecord
	pk := ownerPK()
	if err := x.table.Get("pk", pk).All(&records); err != nil {
		if !isNotFoundErr(err) {
			return nil, goerr.Wrap(err)
		}
	}

	owners := make([]*model.Owner, len(records))
	for i := range records {
		if err := records[i].Unmarshal(&owners[i]); err != nil {
			return nil, err
		}
	}
	return owners, nil
}
