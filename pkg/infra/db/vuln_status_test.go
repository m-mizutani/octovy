package db_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVulnResponse(t *testing.T) {
	t.Run("Insert and find vulnResponse", func(t *testing.T) {
		client := newTestTable(t)
		values := []*model.VulnStatus{
			{
				GitHubRepo: model.GitHubRepo{
					Owner:    "clock",
					RepoName: "tower",
				},
				VulnPackageKey: model.VulnPackageKey{
					Source:  "package-lock.json",
					PkgType: model.PkgNPM,
					PkgName: "blue",
					VulnID:  "v1",
				},
				Status:    model.StatusNone,
				CreatedAt: 100,
			},
			{
				GitHubRepo: model.GitHubRepo{
					Owner:    "clock",
					RepoName: "tower",
				},
				VulnPackageKey: model.VulnPackageKey{
					Source:  "package-lock.json",
					PkgType: model.PkgNPM,
					PkgName: "orange",
					VulnID:  "v2",
				},
				Status:    model.StatusNone,
				CreatedAt: 100,
			},
			{
				GitHubRepo: model.GitHubRepo{
					Owner:    "clock",
					RepoName: "tower",
				},
				VulnPackageKey: model.VulnPackageKey{
					Source:  "package-lock.json",
					PkgType: model.PkgNPM,
					PkgName: "black",
					VulnID:  "v3",
				},
				Status:    model.StatusNone,
				ExpiresAt: 1200,
				CreatedAt: 100,
			},
			{
				GitHubRepo: model.GitHubRepo{
					Owner:    "clock",
					RepoName: "bridge",
				},
				VulnPackageKey: model.VulnPackageKey{
					Source:  "package-lock.json",
					PkgType: model.PkgNPM,
					PkgName: "red",
					VulnID:  "v3",
				},
				Status:    model.StatusNone,
				CreatedAt: 100,
			},
			{
				GitHubRepo: model.GitHubRepo{
					Owner:    "garden",
					RepoName: "tower",
				},
				VulnPackageKey: model.VulnPackageKey{
					Source:  "package-lock.json",
					PkgType: model.PkgNPM,
					PkgName: "white",
					VulnID:  "v2",
				},
				Status:    model.StatusNone,
				CreatedAt: 100,
			},
		}

		for _, v := range values {
			require.NoError(t, client.PutVulnStatus(v))
		}

		r1, err := client.GetVulnStatus(&model.GitHubRepo{
			Owner:    "clock",
			RepoName: "tower",
		}, 1199)
		require.NoError(t, err)
		require.Len(t, r1, 3)
		assert.Contains(t, r1, values[0])
		assert.Contains(t, r1, values[1])
		assert.Contains(t, r1, values[2])

		r2, err := client.GetVulnStatus(&model.GitHubRepo{
			Owner:    "clock",
			RepoName: "tower",
		}, 1200)
		require.NoError(t, err)
		require.Len(t, r2, 2)
		assert.Contains(t, r2, values[0])
		assert.Contains(t, r2, values[1])
	})

	t.Run("Overwrite vulnResponse", func(t *testing.T) {
		client := newTestTable(t)
		oldResp := model.VulnStatus{
			GitHubRepo: model.GitHubRepo{
				Owner:    "clock",
				RepoName: "tower",
			},
			VulnPackageKey: model.VulnPackageKey{
				Source:  "package-lock.json",
				PkgType: model.PkgNPM,
				PkgName: "blue",
				VulnID:  "v1",
			},
			Status:    model.StatusNone,
			CreatedAt: 100,
		}
		newResp := model.VulnStatus{
			GitHubRepo: model.GitHubRepo{
				Owner:    "clock",
				RepoName: "tower",
			},
			VulnPackageKey: model.VulnPackageKey{
				Source:  "package-lock.json",
				PkgType: model.PkgNPM,
				PkgName: "blue",
				VulnID:  "v1",
			},
			Status:    model.StatusSnoozed,
			CreatedAt: 101,
			ExpiresAt: 1100,
		}
		newRespWithOtherVuln := model.VulnStatus{
			GitHubRepo: model.GitHubRepo{
				Owner:    "clock",
				RepoName: "tower",
			},
			VulnPackageKey: model.VulnPackageKey{
				Source:  "package-lock.json",
				PkgType: model.PkgNPM,
				PkgName: "blue",
				VulnID:  "v2",
			},
			Status:    model.StatusSnoozed,
			CreatedAt: 103,
			ExpiresAt: 1100,
		}
		require.NoError(t, client.PutVulnStatus(&oldResp))
		require.NoError(t, client.PutVulnStatus(&newResp))
		r1, err := client.GetVulnStatus(&model.GitHubRepo{
			Owner:    "clock",
			RepoName: "tower",
		}, 1000)
		require.NoError(t, err)
		require.Len(t, r1, 1)
		assert.Equal(t, newResp, *r1[0])

		logs, err := client.GetVulnStatusLogs(&newResp.GitHubRepo, &newResp.VulnPackageKey)
		require.NoError(t, err)
		require.Len(t, logs, 2)
		assert.Equal(t, logs[0], &oldResp)
		assert.Equal(t, logs[1], &newResp)

		require.NoError(t, client.PutVulnStatus(&newRespWithOtherVuln))
		r2, err := client.GetVulnStatus(&model.GitHubRepo{
			Owner:    "clock",
			RepoName: "tower",
		}, 1000)
		require.NoError(t, err)
		require.Len(t, r2, 2)
	})
}
