package usecase

import (
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

func (x *Default) SendScanRequest(req *model.ScanRepositoryRequest) error {
	return x.svc.SendScanRequest(req)
}
