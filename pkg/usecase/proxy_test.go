package usecase_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRepositories(t *testing.T) {
	uc, mock := setupUsecase(t)
	injectGitHubMock(t, mock)
	ctx := context.Background()
	branch := "main"

	uc.SendScanRequest(&model.ScanRepositoryRequest{
		InstallID: 1,
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "blue",
					RepoName: "five",
				},
				Branch: branch,
			},
			CommitID:    "1234567",
			UpdatedAt:   2000,
			RequestedAt: 2100,
		},
	})
	usecase.CloseScanQueue(uc)

	require.NoError(t, uc.Init())

	t.Run("no result before scan", func(t *testing.T) {
		r0, err := uc.GetRepositories(ctx)
		require.NoError(t, err)
		require.Len(t, r0, 0)
	})

	_, err := mock.DB.CreateRepo(ctx, &ent.Repository{
		Owner:         "blue",
		Name:          "five",
		DefaultBranch: &branch,
	})
	require.NoError(t, err)
	require.NoError(t, usecase.RunScanThread(uc))

	t.Run("got result after scan", func(t *testing.T) {
		resp, err := uc.GetRepositories(ctx)
		require.NoError(t, err)
		require.Len(t, resp, 1)
		assert.Equal(t, "blue", resp[0].Owner)
		assert.Equal(t, "five", resp[0].Name)
		require.Len(t, resp[0].Edges.Main, 1)
		assert.Equal(t, "1234567", resp[0].Edges.Main[0].CommitID)
	})
}
