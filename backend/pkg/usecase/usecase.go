package usecase

import (
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
	"github.com/m-mizutani/octovy/backend/pkg/service"
)

var logger = golambda.Logger

type Default struct {
	config *model.Config
	svc    *service.Service
}

func New(cfg *model.Config) interfaces.Usecases {
	return &Default{
		config: cfg,
		svc:    service.New(cfg),
	}
}
