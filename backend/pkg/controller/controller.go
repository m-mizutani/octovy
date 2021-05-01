package controller

import (
	"github.com/m-mizutani/octovy/backend/pkg/infra"
	"github.com/m-mizutani/octovy/backend/pkg/service"
	"github.com/m-mizutani/octovy/backend/pkg/usecase"
)

type Controller struct {
	Infra   *infra.Interfaces
	Config  *service.Config
	Usecase usecase.Usecases
}

func New() *Controller {
	ctrl := &Controller{
		Usecase: usecase.New(),
		Config:  service.NewConfig(),
	}
	return ctrl
}
