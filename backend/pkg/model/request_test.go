package model_test

import (
	"testing"

	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestScanRepositoryRequest(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		assert.NoError(t, (&model.ScanRepositoryRequest{
			ScanTarget: model.ScanTarget{
				GitHubBranch: model.GitHubBranch{
					GitHubRepo: model.GitHubRepo{
						Owner:    "five",
						RepoName: "blue",
					},
					Branch: "master",
				},
				Ref:       "beefcafe",
				UpdatedAt: 0,
			},
			InstallID: 1,
		}).IsValid())
	})

	t.Run("Invalid case", func(t *testing.T) {
		t.Run("No Owner", func(t *testing.T) {
			assert.ErrorIs(t, (&model.ScanRepositoryRequest{
				ScanTarget: model.ScanTarget{
					GitHubBranch: model.GitHubBranch{
						GitHubRepo: model.GitHubRepo{
							RepoName: "blue",
						},
						Branch: "master",
					},
					Ref:       "beefcafe",
					UpdatedAt: 0,
				},
				InstallID: 1,
			}).IsValid(), model.ErrInvalidScanRequest)
		})

		t.Run("No RepoName", func(t *testing.T) {
			assert.ErrorIs(t, (&model.ScanRepositoryRequest{
				ScanTarget: model.ScanTarget{
					GitHubBranch: model.GitHubBranch{
						GitHubRepo: model.GitHubRepo{
							Owner: "five",
						},
						Branch: "master",
					},
					Ref:       "beefcafe",
					UpdatedAt: 0,
				},
				InstallID: 1,
			}).IsValid(), model.ErrInvalidScanRequest)
		})

		t.Run("No Branch", func(t *testing.T) {
			assert.ErrorIs(t, (&model.ScanRepositoryRequest{
				ScanTarget: model.ScanTarget{
					GitHubBranch: model.GitHubBranch{
						GitHubRepo: model.GitHubRepo{
							Owner:    "five",
							RepoName: "blue",
						},
					},
					Ref:       "beefcafe",
					UpdatedAt: 0,
				},
				InstallID: 1,
			}).IsValid(), model.ErrInvalidScanRequest)
		})

		t.Run("No CommitID", func(t *testing.T) {
			assert.ErrorIs(t, (&model.ScanRepositoryRequest{
				ScanTarget: model.ScanTarget{
					GitHubBranch: model.GitHubBranch{
						GitHubRepo: model.GitHubRepo{
							Owner:    "five",
							RepoName: "blue",
						},
						Branch: "master",
					},
					UpdatedAt: 0,
				},
				InstallID: 1,
			}).IsValid(), model.ErrInvalidScanRequest)
		})

		t.Run("No InstallID", func(t *testing.T) {
			assert.ErrorIs(t, (&model.ScanRepositoryRequest{
				ScanTarget: model.ScanTarget{
					GitHubBranch: model.GitHubBranch{
						GitHubRepo: model.GitHubRepo{
							Owner:    "five",
							RepoName: "blue",
						},
						Branch: "master",
					},
					Ref:       "beefcafe",
					UpdatedAt: 0,
				},
				InstallID: 0,
			}).IsValid(), model.ErrInvalidScanRequest)
		})
	})
}
