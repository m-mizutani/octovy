package usecase

import (
	"context"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

func (x *usecase) RegisterRepository(ctx context.Context, repo *ent.Repository) (*ent.Repository, error) {
	if !x.initialized {
		panic("usecase is not initialized")
	}

	return x.infra.DB.CreateRepo(ctx, repo)
}

func (x *usecase) UpdateVulnStatus(ctx context.Context, req *model.UpdateVulnStatusRequest) error {
	panic("not implemented") // TODO: Implement
}

func (x *usecase) LookupScanReport(ctx context.Context, scanID string) (*ent.Scan, error) {
	return x.infra.DB.GetScan(context.Background(), scanID)
}
