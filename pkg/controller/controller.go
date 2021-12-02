package controller

import (
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/usecase"
)

type Controller struct {
	Config  *model.Config
	usecase *usecase.Usecase
}

func New() *Controller {
	var cfg model.Config
	return &Controller{
		Config:  &cfg,
		usecase: usecase.New(&cfg),
	}
}
