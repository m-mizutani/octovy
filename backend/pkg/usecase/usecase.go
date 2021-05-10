package usecase

import (
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/m-mizutani/octovy/backend/pkg/service"
)

var logger = golambda.Logger

type Usecases interface {
	ScanRepository(svc *service.Service, req *model.ScanRepositoryRequest) error
	SendScanRequest(svc *service.Service, req *model.ScanRepositoryRequest) error

	RegisterRepository(svc *service.Service, repo *model.Repository) error
	PutNewRepository(svc *service.Service, repo *model.Repository) (bool, error)
	UpdateRepositoryDefaultBranch(svc *service.Service, repo *model.GitHubRepo, branch string) error

	FindRepos(svc *service.Service) ([]*model.Repository, error)
	FindReposByOwner(svc *service.Service, owner string) ([]*model.Repository, error)
	FindReposByFullName(svc *service.Service, owner, name string) (*model.Repository, error)
	FindPkgs(svc *service.Service, pkgType model.PkgType, name string) ([]*model.PackageRecord, error)
	FindPkgsByRepo(svc *service.Service, branch *model.GitHubBranch) ([]*model.PackageRecord, error)
}

type Default struct{}

func New() Usecases {
	return &Default{}
}
