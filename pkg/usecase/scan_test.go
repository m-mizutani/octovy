package usecase_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanProcedure(t *testing.T) {
	uc, mock := setupUsecase(t)
	injectGitHubMock(t, mock)

	uc.SendScanRequest(&model.ScanRepositoryRequest{
		InstallID: 1,
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "five",
				},
				Branch: "main",
			},
			CommitID:    "1234567",
			UpdatedAt:   2000,
			RequestedAt: 2100,
		},
	})
	usecase.CloseScanQueue(uc)

	require.NoError(t, uc.Init())
	require.NoError(t, usecase.RunScanThread(uc))

	ctx := context.Background()
	report, err := mock.DB.GetLatestScan(ctx, model.GitHubBranch{
		GitHubRepo: model.GitHubRepo{
			Owner:    "blue",
			RepoName: "five",
		},
		Branch: "main",
	})
	require.NoError(t, err)
	assert.Equal(t, "1234567", report.CommitID)
}
