package trivydb_test

import (
	"os"
	"testing"

	"github.com/m-mizutani/octovy/backend/pkg/infra"
	"github.com/m-mizutani/octovy/backend/pkg/infra/trivydb"
	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDB(t *testing.T) infra.TrivyDBClient {
	dbPath := os.Getenv("TRIVY_DB_PATH")
	if dbPath == "" {
		t.Skip("TRIVY_DB_PATH is not set")
	}

	db, err := trivydb.New(dbPath)
	require.NoError(t, err)

	return db
}

func TestDBAccessAdvisory(t *testing.T) {
	db := setupDB(t)

	t.Run("Get advisories", func(t *testing.T) {
		/*
			bson @ rubygems:
			- CVE-2015-4411: {"PatchedVersions":["\u003e= 3.0.4"]}
			- CVE-2015-4412: {"PatchedVersions":["~\u003e 1.12.3","\u003e= 3.0.4"]}
		*/
		adv, err := db.GetAdvisories("ruby-advisory-db", "bson")
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(adv), 2)

		cve2015_4411 := &model.AdvisoryData{
			VulnID: "CVE-2015-4411",
			Data:   []byte(`{"PatchedVersions":["\u003e= 3.0.4"]}`),
		}
		cve2015_4412 := &model.AdvisoryData{
			VulnID: "CVE-2015-4412",
			Data:   []byte(`{"PatchedVersions":["~\u003e 1.12.3","\u003e= 3.0.4"]}`),
		}

		assert.Contains(t, adv, cve2015_4411)
		assert.Contains(t, adv, cve2015_4412)
	})

	t.Run("Get error with invalid source name", func(t *testing.T) {
		_, err := db.GetAdvisories("hoge", "bson")
		require.Error(t, err)
	})

	t.Run("Get nil with non vulnerable package name", func(t *testing.T) {
		adv, err := db.GetAdvisories("ruby-advisory-db", "xyz")
		require.NoError(t, err)
		assert.Zero(t, len(adv))
	})
}

func TestDBAccessVulnerability(t *testing.T) {
	db := setupDB(t)

	t.Run("Get vulnerability", func(t *testing.T) {
		vuln, err := db.GetVulnerability("CVE-2015-4411")
		require.NoError(t, err)
		require.NotNil(t, vuln)
		assert.Equal(t, "rubygem-moped: Denial of Service with crafted ObjectId string (incomplete fix for CVE-2015-4410)", vuln.Title)
		assert.NotZero(t, len(vuln.References))
	})

	t.Run("Get no vulnerability with invalid CVE number", func(t *testing.T) {
		vuln, err := db.GetVulnerability("CVE-2015-XXXX")
		require.NoError(t, err)
		require.Nil(t, vuln)
	})
}
