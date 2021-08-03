package usecase

import (
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func (x *Default) RegisterRepository(repo *model.Repository) error {
	inserted, err := x.PutNewRepository(repo)
	if err != nil {
		return err
	}

	if !inserted {
		if err := x.UpdateRepositoryDefaultBranch(&repo.GitHubRepo, repo.DefaultBranch); err != nil {
			return err
		}
	}
	return nil
}

func (x *Default) PutNewRepository(repo *model.Repository) (bool, error) {
	return x.svc.DB().InsertRepo(repo)
}

func (x *Default) UpdateRepositoryDefaultBranch(repo *model.GitHubRepo, branch string) error {
	return x.svc.DB().SetRepoDefaultBranchName(repo, branch)
}

func (x *Default) FindOwners() ([]*model.Owner, error) {
	return x.svc.DB().FindOwners()
}

func (x *Default) FindRepos() ([]*model.Repository, error) {
	return x.svc.DB().FindRepo()
}

func (x *Default) FindReposByOwner(owner string) ([]*model.Repository, error) {
	return x.svc.DB().FindRepoByOwner(owner)
}

func (x *Default) FindReposByFullName(owner string, name string) (*model.Repository, error) {
	repos, err := x.svc.DB().FindRepoByOwner(owner)
	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		if repo.RepoName == name {
			return repo, nil
		}
	}
	return nil, nil
}

func (x *Default) LookupBranch(branch *model.GitHubBranch) (*model.Branch, error) {
	return x.svc.DB().LookupBranch(branch)
}

func (x *Default) FindPkgs(pkgType model.PkgType, name string) ([]*model.PackageRecord, error) {
	return x.svc.DB().FindPackageRecordsByName(pkgType, name)
}

func (x *Default) FindPkgsByRepo(branch *model.GitHubBranch) ([]*model.PackageRecord, error) {
	return x.svc.DB().FindPackageRecordsByBranch(branch)
}
