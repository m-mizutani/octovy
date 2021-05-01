package usecase_test

import (
	"testing"

	"github.com/m-mizutani/octovy/backend/pkg/infra"
	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/m-mizutani/octovy/backend/pkg/service"
	"github.com/m-mizutani/octovy/backend/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPutNewRepository(t *testing.T) *service.Service {
	dbClient := newTestTable(t)

	cfg := service.NewConfig()
	cfg.TableName = dbClient.TableName()
	svc := service.New(cfg)
	svc.NewDB = func(region, tableName string) (infra.DBClient, error) {
		return dbClient, nil
	}
	return svc
}

func TestPutNewRepository(t *testing.T) {
	t.Run("put repositories", func(t *testing.T) {
		svc := setupPutNewRepository(t)
		repo1 := &model.Repository{
			GitHubRepo: model.GitHubRepo{
				Owner:    "five",
				RepoName: "blue",
			},
			URL:           "https://github-enterprise.example.com/five/blue",
			Branches:      []string{"master"},
			DefaultBranch: "main",
		}

		repo2 := &model.Repository{
			GitHubRepo: model.GitHubRepo{
				Owner:    "five",
				RepoName: "orange",
			},
			URL:           "https://github-enterprise.example.com/five/orange",
			Branches:      []string{"master"},
			DefaultBranch: "main",
		}

		repo3 := &model.Repository{
			GitHubRepo: model.GitHubRepo{
				Owner:    "three",
				RepoName: "heaven",
			},
			URL:           "https://github-enterprise.example.com/three/heaven",
			Branches:      []string{"master"},
			DefaultBranch: "main",
		}

		testInsert := func(t *testing.T, repo *model.Repository) {
			inserted, err := usecase.New().PutNewRepository(svc, repo)
			require.NoError(t, err)
			assert.True(t, inserted)
		}
		testInsert(t, repo1)
		testInsert(t, repo2)
		testInsert(t, repo3)

		result1, err := svc.DB().FindRepo()
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
		assert.Contains(t, r1.Branches, "master")
		assert.Equal(t, "main", r1.DefaultBranch)

		// Find "five" owner repository
		result2, err := svc.DB().FindRepoByOwner("five")
		require.NoError(t, err)
		assert.Contains(t, result2, repo1)
		assert.Contains(t, result2, repo2)
		assert.NotContains(t, result2, repo3)
	})

	t.Run("Change branches and default branch", func(t *testing.T) {
		svc := setupPutNewRepository(t)
		repo1 := &model.Repository{
			GitHubRepo: model.GitHubRepo{
				Owner:    "five",
				RepoName: "blue",
			},
			URL:           "https://github-enterprise.example.com/five/blue",
			Branches:      []string{"master"},
			DefaultBranch: "main",
		}

		repo2 := &model.Repository{
			GitHubRepo:    repo1.GitHubRepo,
			URL:           repo1.URL,
			Branches:      []string{"x"},
			DefaultBranch: "y",
		}

		uc := usecase.New()
		require.NoError(t, uc.RegisterRepository(svc, repo1))
		require.NoError(t, uc.RegisterRepository(svc, repo2))

		result1, err := svc.DB().FindRepo()
		require.NoError(t, err)
		require.Equal(t, 1, len(result1))
		assert.Equal(t, []string{"master"}, result1[0].Branches)
		assert.Equal(t, "y", result1[0].DefaultBranch)
	})
}
