package usecase

import (
	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/m-mizutani/octovy/backend/pkg/service"
)

func (x *Default) RegisterRepository(svc *service.Service, repo *model.Repository) error {
	inserted, err := x.PutNewRepository(svc, repo)
	if err != nil {
		return err
	}

	if !inserted {
		if err := x.UpdateRepositoryDefaultBranch(svc, &repo.GitHubRepo, repo.DefaultBranch); err != nil {
			return err
		}
	}
	return nil
}

func (x *Default) PutNewRepository(svc *service.Service, repo *model.Repository) (bool, error) {
	return svc.DB().InsertRepo(repo)
}

func (x *Default) UpdateRepositoryDefaultBranch(svc *service.Service, repo *model.GitHubRepo, branch string) error {
	return svc.DB().SetRepoDefaultBranch(repo, branch)
}

func (x *Default) UpdateRepositoryBranches(svc *service.Service, repo *model.GitHubRepo, branches []string) error {
	return svc.DB().SetRepoBranches(repo, branches)
}

func (x *Default) FindRepos(svc *service.Service) ([]*model.Repository, error) {
	return svc.DB().FindRepo()
}

func (x *Default) FindReposByOwner(svc *service.Service, owner string) ([]*model.Repository, error) {
	return svc.DB().FindRepoByOwner(owner)
}

func (x *Default) FindReposByFullName(svc *service.Service, owner string, name string) (*model.Repository, error) {
	repos, err := svc.DB().FindRepoByOwner(owner)
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

func (x *Default) FindPkgs(svc *service.Service, pkgType model.PkgType, name string) ([]*model.Package, error) {
	return svc.DB().FindPackagesByName(pkgType, name)
}

func (x *Default) FindPkgsByRepo(svc *service.Service, branch *model.GitHubBranch) ([]*model.Package, error) {
	return svc.DB().FindPackagesByBranch(branch)
}
