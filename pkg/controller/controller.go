package controller

import (
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/usecase"
)

type Controller struct {
	Config  *model.Config
	Usecase interfaces.Usecases
}

func New() *Controller {
	ctrl := &Controller{
		Config: model.NewConfig(),
	}

	ctrl.Usecase = usecase.New(ctrl.Config)
	return ctrl
}
