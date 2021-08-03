package usecase

import (
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/service"
)

var logger = golambda.Logger

type Default struct {
	config    *model.Config
	svc       *service.Service
	scanQueue chan *model.ScanRepositoryRequest
}

func New(cfg *model.Config) interfaces.Usecases {
	return &Default{
		config:    cfg,
		svc:       service.New(cfg),
		scanQueue: make(chan *model.ScanRepositoryRequest, 256),
	}
}
