package scanner

import (
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
)

var logger = golambda.Logger

type Scanner struct {
	usecase interfaces.Usecases
}

func New(uc interfaces.Usecases) *Scanner {
	return &Scanner{
		usecase: uc,
	}
}

func (x *Scanner) Run() error {
	for {
		req := x.usecase.RecvScanRequest()
		logger.With("req", req).Info("Recv scan request")
		if err := x.usecase.ScanRepository(req); err != nil {
			golambda.EmitError(err)
			logger.With("err", err).Error("Error in ScanRepository")
		}
	}
}
