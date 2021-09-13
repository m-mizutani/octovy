package db_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	client := setupDB(t)
	ctx := context.Background()

	vulnSet := []*ent.Vulnerability{
		{
			ID:             "CVE-2001-1000",
			FirstSeenAt:    1000,
			LastModifiedAt: 1200,
			Title:          "blue",
			Description:    "5",
		},
		{
			ID:             "CVE-2002-1000",
			FirstSeenAt:    2000,
			LastModifiedAt: 1345,
			Title:          "orange",
		},
	}

	pkgSet := []*ent.PackageRecord{
		{
			Type:    types.PkgGoModule,
			Source:  "go.mod",
			Name:    "xxx",
			Version: "v0.1.1",
			VulnIds: []string{"CVE-2001-1000", "CVE-2002-1000"},
		},
	}

	scan := &ent.Scan{
		CommitID:    "1234567",
		RequestedAt: 100,
		ScannedAt:   200,
		CheckID:     999,
	}

	require.NoError(t, client.PutVulnerabilities(ctx, vulnSet))

	addedPkg, err := client.PutPackages(ctx, pkgSet)
	require.NoError(t, err)

	branch, err := client.GetBranch(ctx, &db.BranchKey{
		Owner:    "blue",
		RepoName: "five",
		Branch:   "main",
	})
	require.NoError(t, err)

	addedscan, err := client.PutScan(ctx, scan, branch, addedPkg)
	require.NoError(t, err)

	got, err := client.GetScan(ctx, addedscan.ID)
	require.NoError(t, err)
	assert.Equal(t, got.CheckID, scan.CheckID)
	require.Len(t, got.Edges.Packages, 1)
}
