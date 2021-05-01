package usecase

import (
	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/m-mizutani/octovy/backend/pkg/service"
)

func (x *Default) SendScanRequest(svc *service.Service, req *model.ScanRepositoryRequest) error {
	return svc.SendScanRequest(req)
}
