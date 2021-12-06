package db_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	ctx := model.NewContext()

	t.Run("got inserted scan report by ID", func(t *testing.T) {
		client := setupDB(t)

		pkgSet := []*ent.PackageRecord{
			{
				Type:    "gomod",
				Source:  "go.mod",
				Name:    "xxx",
				Version: "v0.1.1",
				VulnIds: []string{"CVE-2001-0001", "CVE-2001-0002"},
			},
		}

		scan := &ent.Scan{
			CommitID:    "1234567",
			Branch:      "main",
			RequestedAt: 100,
			ScannedAt:   200,
			CheckID:     999,
		}

		addedPkg, err := client.PutPackages(ctx, pkgSet)
		require.NoError(t, err)

		repo, err := client.CreateRepo(ctx, &ent.Repository{
			Owner:     "blue",
			Name:      "five",
			InstallID: 1,
		})
		require.NoError(t, err)

		addedscan, err := client.PutScan(ctx, scan, repo, addedPkg)
		require.NoError(t, err)

		got, err := client.GetScan(ctx, addedscan.ID)
		require.NoError(t, err)
		assert.Equal(t, got.CheckID, scan.CheckID)
		require.Len(t, got.Edges.Packages, 1)
	})

	t.Run("lookup latest scan report by branch", func(t *testing.T) {
		client := setupDB(t)

		repo, err := client.CreateRepo(ctx, &ent.Repository{
			Owner:     "blue",
			Name:      "five",
			InstallID: 1,
		})
		require.NoError(t, err)
		repo_another, err := client.CreateRepo(ctx, &ent.Repository{
			Owner:     "orange",
			Name:      "doll",
			InstallID: 1,
		})
		require.NoError(t, err)

		pkgSet1, err := client.PutPackages(ctx, []*ent.PackageRecord{
			{
				Type:    "gomod",
				Source:  "go.mod",
				Name:    "x",
				Version: "v0.1.1",
				VulnIds: []string{"CVE-2001-0001", "CVE-2001-0002"},
			},
		})
		require.NoError(t, err)
		pkgSet2, err := client.PutPackages(ctx, []*ent.PackageRecord{
			{
				Type:    "gomod",
				Source:  "go.mod",
				Name:    "y",
				Version: "v0.1.1",
				VulnIds: []string{"CVE-2001-0001", "CVE-2001-0002"},
			},
		})
		require.NoError(t, err)

		scan1 := &ent.Scan{
			CommitID:    "aaa",
			Branch:      "main",
			RequestedAt: 100,
			ScannedAt:   200,
			CheckID:     999,
		}
		scan2 := &ent.Scan{
			CommitID:    "bbb",
			Branch:      "main",
			RequestedAt: 100,
			ScannedAt:   199,
			CheckID:     999,
		}
		scan3 := &ent.Scan{
			CommitID:    "ccc",
			Branch:      "other-branch",
			RequestedAt: 100,
			ScannedAt:   200,
			CheckID:     999,
		}
		scan4 := &ent.Scan{
			CommitID:    "ddd",
			Branch:      "main",
			RequestedAt: 100,
			ScannedAt:   200,
			CheckID:     999,
		}

		// target
		_, err = client.PutScan(ctx, scan1, repo, pkgSet1)
		require.NoError(t, err)
		// same repo/branch, but old
		_, err = client.PutScan(ctx, scan2, repo, pkgSet2)
		require.NoError(t, err)
		// same repo, but other branch
		_, err = client.PutScan(ctx, scan3, repo, pkgSet2)
		require.NoError(t, err)
		// same branch, but other repo
		_, err = client.PutScan(ctx, scan4, repo_another, pkgSet2)
		require.NoError(t, err)

		latest, err := client.GetLatestScan(ctx, model.GitHubBranch{
			GitHubRepo: model.GitHubRepo{
				Owner: "blue",
				Name:  "five",
			},
			Branch: "main",
		})
		require.NoError(t, err)
		require.NotNil(t, latest)
		assert.Equal(t, "aaa", latest.CommitID)
		require.Len(t, latest.Edges.Repository, 1)
		assert.Equal(t, "five", latest.Edges.Repository[0].Name)
		require.Len(t, latest.Edges.Packages, 1)
		assert.Equal(t, "x", latest.Edges.Packages[0].Name)
	})

	t.Run("save default branch latest", func(t *testing.T) {
		client := setupDB(t)
		defaultBranch := "main"
		repo, err := client.CreateRepo(ctx, &ent.Repository{
			Owner:         "blue",
			Name:          "five",
			InstallID:     1,
			DefaultBranch: &defaultBranch,
		})
		require.NoError(t, err)

		s1, err := client.PutScan(ctx, &ent.Scan{
			Branch:      defaultBranch,
			CommitID:    "aaa",
			RequestedAt: 100,
			ScannedAt:   200,
			CheckID:     999,
		}, repo, nil)
		require.NoError(t, err)

		_, err = client.PutScan(ctx, &ent.Scan{
			Branch:      defaultBranch,
			CommitID:    "345",
			RequestedAt: 100,
			ScannedAt:   100, // older than s1
			CheckID:     999,
		}, repo, nil)
		require.NoError(t, err)

		_, err = client.PutScan(ctx, &ent.Scan{
			Branch:      "not-default-branch",
			CommitID:    "a12",
			RequestedAt: 100,
			ScannedAt:   300, // newer than s1, but not default
			CheckID:     999,
		}, repo, nil)
		require.NoError(t, err)

		scans, err := client.GetLatestScans(ctx)
		require.NoError(t, err)
		require.Len(t, scans, 1)
		assert.Equal(t, s1.ID, scans[0].ID)
	})
}
