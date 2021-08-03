package usecase

import (
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func (x *Default) SendScanRequest(req *model.ScanRepositoryRequest) error {
	x.scanQueue <- req
	return nil
}

func (x *Default) RecvScanRequest() *model.ScanRepositoryRequest {
	return <-x.scanQueue
}
