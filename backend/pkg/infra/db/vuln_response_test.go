package db_test

import (
	"testing"

	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVulnResponse(t *testing.T) {
	t.Run("Insert and find vulnResponse", func(t *testing.T) {
		client := newTestTable(t)
		values := []*model.VulnResponse{
			{
				GitHubRepo: model.GitHubRepo{
					Owner:    "clock",
					RepoName: "tower",
				},
				PkgType:   model.PkgNPM,
				PkgName:   "blue",
				VulnID:    "v1",
				Type:      model.RespNever,
				CreatedAt: 1000,
			},
			{
				GitHubRepo: model.GitHubRepo{
					Owner:    "clock",
					RepoName: "tower",
				},
				PkgType:   model.PkgNPM,
				PkgName:   "orange",
				VulnID:    "v2",
				Type:      model.RespNever,
				CreatedAt: 1000,
			},
			{
				GitHubRepo: model.GitHubRepo{
					Owner:    "clock",
					RepoName: "tower",
				},
				PkgType:   model.PkgNPM,
				PkgName:   "black",
				VulnID:    "v3",
				Type:      model.RespNever,
				CreatedAt: 1000,
				Duration:  200,
			},
			{
				GitHubRepo: model.GitHubRepo{
					Owner:    "clock",
					RepoName: "bridge",
				},
				PkgType:   model.PkgNPM,
				PkgName:   "red",
				VulnID:    "v3",
				Type:      model.RespNever,
				CreatedAt: 1000,
			},
			{
				GitHubRepo: model.GitHubRepo{
					Owner:    "garden",
					RepoName: "tower",
				},
				PkgType:   model.PkgNPM,
				PkgName:   "white",
				VulnID:    "v2",
				Type:      model.RespNever,
				CreatedAt: 1000,
			},
		}

		for _, v := range values {
			require.NoError(t, client.PutVulnResponse(v))
		}

		r1, err := client.GetVulnResponses(&model.GitHubRepo{
			Owner:    "clock",
			RepoName: "tower",
		}, 1199)
		require.NoError(t, err)
		require.Len(t, r1, 3)
		assert.Contains(t, r1, values[0])
		assert.Contains(t, r1, values[1])
		assert.Contains(t, r1, values[2])

		r2, err := client.GetVulnResponses(&model.GitHubRepo{
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
		oldResp := model.VulnResponse{
			GitHubRepo: model.GitHubRepo{
				Owner:    "clock",
				RepoName: "tower",
			},
			PkgType:   model.PkgNPM,
			PkgName:   "blue",
			VulnID:    "v1",
			Type:      model.RespNever,
			CreatedAt: 1000,
		}
		newResp := model.VulnResponse{
			GitHubRepo: model.GitHubRepo{
				Owner:    "clock",
				RepoName: "tower",
			},
			PkgType:   model.PkgNPM,
			PkgName:   "blue",
			VulnID:    "v1",
			Type:      model.RespSnooze,
			CreatedAt: 1000,
		}
		newRespWithOtherVuln := model.VulnResponse{
			GitHubRepo: model.GitHubRepo{
				Owner:    "clock",
				RepoName: "tower",
			},
			PkgType:   model.PkgNPM,
			PkgName:   "blue",
			VulnID:    "v2",
			Type:      model.RespSnooze,
			CreatedAt: 1000,
		}
		require.NoError(t, client.PutVulnResponse(&oldResp))
		require.NoError(t, client.PutVulnResponse(&newResp))
		r1, err := client.GetVulnResponses(&model.GitHubRepo{
			Owner:    "clock",
			RepoName: "tower",
		}, 1000)
		require.NoError(t, err)
		require.Len(t, r1, 1)
		assert.Equal(t, newResp, *r1[0])

		require.NoError(t, client.PutVulnResponse(&newRespWithOtherVuln))
		r2, err := client.GetVulnResponses(&model.GitHubRepo{
			Owner:    "clock",
			RepoName: "tower",
		}, 1000)
		require.NoError(t, err)
		require.Len(t, r2, 2)
	})

	t.Run("Remove and not found", func(t *testing.T) {
		client := newTestTable(t)
		newVulnResp := func() *model.VulnResponse {
			return &model.VulnResponse{
				GitHubRepo: model.GitHubRepo{
					Owner:    "clock",
					RepoName: "tower",
				},
				PkgType:   model.PkgNPM,
				PkgName:   "blue",
				VulnID:    "v1",
				Type:      model.RespNever,
				CreatedAt: 1000,
			}

		}
		v := newVulnResp()
		require.NoError(t, client.PutVulnResponse(v))

		// Not affect
		t.Run("Diff Owner", func(t *testing.T) {
			d := newVulnResp()
			d.Owner = "garden"
			require.NoError(t, client.DeleteVulnResponse(d))
			r, err := client.GetVulnResponses(&v.GitHubRepo, 1000)
			require.NoError(t, err)
			assert.Len(t, r, 1)
		})

		t.Run("Diff RepoName", func(t *testing.T) {
			d := newVulnResp()
			d.RepoName = "bridge"
			require.NoError(t, client.DeleteVulnResponse(d))
			r, err := client.GetVulnResponses(&v.GitHubRepo, 1000)
			require.NoError(t, err)
			assert.Len(t, r, 1)
		})

		t.Run("Diff PkgType", func(t *testing.T) {
			d := newVulnResp()
			d.PkgType = model.PkgGoModule
			require.NoError(t, client.DeleteVulnResponse(d))
			r, err := client.GetVulnResponses(&v.GitHubRepo, 1000)
			require.NoError(t, err)
			assert.Len(t, r, 1)
		})

		t.Run("Diff PkgName", func(t *testing.T) {
			d := newVulnResp()
			d.PkgName = "orange"
			require.NoError(t, client.DeleteVulnResponse(d))
			r, err := client.GetVulnResponses(&v.GitHubRepo, 1000)
			require.NoError(t, err)
			assert.Len(t, r, 1)
		})

		t.Run("Diff VulnID", func(t *testing.T) {
			d := newVulnResp()
			d.VulnID = "v0"
			require.NoError(t, client.DeleteVulnResponse(d))
			r, err := client.GetVulnResponses(&v.GitHubRepo, 1000)
			require.NoError(t, err)
			assert.Len(t, r, 1)
		})

		// Affect
		t.Run("Removable with same key", func(t *testing.T) {
			d := newVulnResp()
			require.NoError(t, client.DeleteVulnResponse(d))
			r, err := client.GetVulnResponses(&v.GitHubRepo, 1000)
			require.NoError(t, err)
			assert.Len(t, r, 0)
		})
	})
}
