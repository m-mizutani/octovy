package usecase

import (
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra"
)

func (x *Usecase) InjectInfra(inject *infra.Interfaces) {
	x.infra = inject
}
func (x *Usecase) DisableInvokeThread() {
	x.disableInvokeThread = true
}

func SetErrorHandler(uc *Usecase, handler func(error)) {
	uc.testErrorHandler = handler
}

func NewUsecase(cfg *model.Config) *Usecase {
	return New(cfg)
}

func RunScanThread(uc *Usecase) error {
	return uc.runScanThread()
}

func CloseScanQueue(uc *Usecase) {
	close(uc.scanQueue)
}

type PostGitHubCommentInput postGitHubCommentInput

func PostGitHubComment(input *PostGitHubCommentInput) error {
	d := postGitHubCommentInput(*input)
	return postGitHubComment(&d)
}
