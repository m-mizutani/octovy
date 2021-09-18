package usecase

import (
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra"
)

type Usecase usecase

func (x *usecase) InjectInfra(inject *infra.Interfaces) {
	x.infra = inject
}

func SetErrorHandler(uc Interface, handler func(error)) {
	uc.(*usecase).testErrorHandler = handler
}

func NewUsecase(cfg *model.Config) *usecase {
	return New(cfg).(*usecase)
}

func RunScanThread(uc Interface) error {
	return uc.(*usecase).runScanThread()
}

func CloseScanQueue(uc Interface) {
	close(uc.(*usecase).scanQueue)
}
