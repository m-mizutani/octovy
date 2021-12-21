package usecase

import (
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func (x *Usecase) PushTrivyResult(ctx *model.Context, req *model.PushTrivyResultRequest) error {
	scannedAt := x.infra.Utils.Now()
	newPkgs, vulnList := model.TrivyReportToEnt(&req.Report, scannedAt)

	newScan, err := insertScan(ctx, x.infra.DB, &req.Target, newPkgs, vulnList, scannedAt)
	if err != nil {
		return err
	}
	ctx.Log().With("scanID", newScan.ID).Debug("inserted scan report")

	return nil
}
