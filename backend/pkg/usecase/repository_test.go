package usecase_test

import (
	"testing"

	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
	"github.com/m-mizutani/octovy/backend/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPutNewRepository(t *testing.T) interfaces.Usecases {
	dbClient := newTestTable(t)

	cfg := model.NewConfig()
	cfg.TableName = dbClient.TableName()

	uc := usecase.New(cfg)
	usecase.InjectDBClient(uc, dbClient)
	return uc
}

func TestPutNewRepository(t *testing.T) {
	t.Run("put repositories", func(t *testing.T) {
		uc := setupPutNewRepository(t)
		repo1 := &model.Repository{
			GitHubRepo: model.GitHubRepo{
				Owner:    "five",
				RepoName: "blue",
			},
			URL:           "https://github-enterprise.example.com/five/blue",
			DefaultBranch: "main",
			Branch: model.Branch{
				ReportSummary: model.ScanReportSummary{PkgTypes: []model.PkgType{}},
			},
		}

		repo2 := &model.Repository{
			GitHubRepo: model.GitHubRepo{
				Owner:    "five",
				RepoName: "orange",
			},
			URL:           "https://github-enterprise.example.com/five/orange",
			DefaultBranch: "main",
			Branch: model.Branch{
				ReportSummary: model.ScanReportSummary{PkgTypes: []model.PkgType{}},
			},
		}

		repo3 := &model.Repository{
			GitHubRepo: model.GitHubRepo{
				Owner:    "three",
				RepoName: "heaven",
			},
			URL:           "https://github-enterprise.example.com/three/heaven",
			DefaultBranch: "main",
			Branch: model.Branch{
				ReportSummary: model.ScanReportSummary{PkgTypes: []model.PkgType{}},
			},
		}

		testInsert := func(t *testing.T, repo *model.Repository) {
			inserted, err := uc.PutNewRepository(repo)
			require.NoError(t, err)
			assert.True(t, inserted)
		}
		testInsert(t, repo1)
		testInsert(t, repo2)
		testInsert(t, repo3)

		db := usecase.EjectDBClient(uc)
		result1, err := db.FindRepo()
		require.NoError(t, err)
		require.Equal(t, 3, len(result1))

		// Find all repository
		var r1 *model.Repository
		for i := range result1 {
			if result1[i].RepoName == "blue" {
				r1 = result1[i]
			}
		}
		require.NotNil(t, r1)
		assert.Equal(t, "five", r1.Owner)
		assert.Equal(t, "blue", r1.RepoName)
		assert.Equal(t, "https://github-enterprise.example.com/five/blue", r1.URL)
		assert.Equal(t, "main", r1.DefaultBranch)

		// Find "five" owner repository
		result2, err := db.FindRepoByOwner("five")
		require.NoError(t, err)
		assert.Contains(t, result2, repo1)
		assert.Contains(t, result2, repo2)
		assert.NotContains(t, result2, repo3)
	})

	t.Run("Change branches and default branch", func(t *testing.T) {
		uc := setupPutNewRepository(t)
		repo1 := &model.Repository{
			GitHubRepo: model.GitHubRepo{
				Owner:    "five",
				RepoName: "blue",
			},
			URL:           "https://github-enterprise.example.com/five/blue",
			DefaultBranch: "main",
		}

		repo2 := &model.Repository{
			GitHubRepo:    repo1.GitHubRepo,
			URL:           repo1.URL,
			DefaultBranch: "y",
		}

		require.NoError(t, uc.RegisterRepository(repo1))
		require.NoError(t, uc.RegisterRepository(repo2))

		result1, err := usecase.EjectDBClient(uc).FindRepo()
		require.NoError(t, err)
		require.Equal(t, 1, len(result1))
		assert.Equal(t, "y", result1[0].DefaultBranch)
	})
}
