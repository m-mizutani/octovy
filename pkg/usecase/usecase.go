package usecase

import (
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra"
)

type UseCase interface {
	ScanGitHubRepo(ctx *model.Context, input *ScanGitHubRepoInput) error
}

type useCase struct {
	clients *infra.Clients
}

func New(clients *infra.Clients) UseCase {
	return &useCase{
		clients: clients,
	}
}
