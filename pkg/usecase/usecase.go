package usecase

import (
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra"
)

type useCase struct {
	tableID types.BQTableID
	clients *infra.Clients
}

func New(clients *infra.Clients, options ...Option) interfaces.UseCase {
	uc := &useCase{
		tableID: "scans",
		clients: clients,
	}

	for _, opt := range options {
		opt(uc)
	}

	return uc
}

type Option func(*useCase)

func WithBigQueryTableID(tableID types.BQTableID) Option {
	return func(x *useCase) {
		x.tableID = tableID
	}
}
