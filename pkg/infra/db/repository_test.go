package db_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRepositories(t *testing.T) {
	str := func(s string) *string { return &s }
	ctx := context.Background()

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

		r1, err := client.CreateRepo(ctx, &ent.Repository{
			Owner:         "blue",
			Name:          "five",
			InstallID:     1,
			DefaultBranch: str("main"),
		})
		require.NoError(t, err)
		r2, err := client.CreateRepo(ctx, &ent.Repository{
			Owner:         "orange",
			Name:          "puppet",
			InstallID:     1,
			DefaultBranch: str("main"),
		})
		require.NoError(t, err)

		_, err = client.PutScan(ctx, scan, r1, addedPkg)
		require.NoError(t, err)
		_, err = client.PutScan(ctx, scan, r2, addedPkg)
		require.NoError(t, err)

		got, err := client.GetRepositories(ctx)
		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.NotNil(t, got[0].Edges.Latest)
		assert.NotNil(t, got[1].Edges.Latest)
	})

}
