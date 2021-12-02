package controller

import (
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

type Controller struct {
	Config *model.Config
}

func New() *Controller {
	var cfg model.Config
	return &Controller{
		Config: &cfg,
	}
}
