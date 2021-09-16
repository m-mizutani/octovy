package db_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/stretchr/testify/require"
)

func setupDB(t *testing.T) db.Interface {
	client := db.NewDBMock(t)
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning failed to close DB: %+v", err)
		}
	})

	// Set default vulnerabilities, CVE-2001-0001 and CVE-2001-0002
	vulnSet := []*ent.Vulnerability{
		{
			ID:             "CVE-2001-0001",
			FirstSeenAt:    1000,
			LastModifiedAt: 1200,
			Title:          "blue",
			Description:    "5",
		},
		{
			ID:             "CVE-2001-0002",
			FirstSeenAt:    2000,
			LastModifiedAt: 1345,
			Title:          "orange",
		},
	}
	ctx := context.Background()
	require.NoError(t, client.PutVulnerabilities(ctx, vulnSet))

	return client
}
