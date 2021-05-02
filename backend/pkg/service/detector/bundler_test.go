package detector_test

import (
	"testing"

	"github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/m-mizutani/octovy/backend/pkg/infra/trivydb"
	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/m-mizutani/octovy/backend/pkg/service/detector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBundler(t *testing.T) {
	_, db := trivydb.NewMock()
	db.AdvisoryMap["ruby-advisory-db"] = map[string][]*model.AdvisoryData{
		"somepkg": {
			&model.AdvisoryData{
				VulnID: "CVE-1234-5678",
				Data:   []byte(`{"PatchedVersions":["\u003e= 3.2.5"]}`),
			},
		},
	}
	db.AdvisoryMap["GitHub Security Advisory Rubygems"] = map[string][]*model.AdvisoryData{
		"somepkg": {
			&model.AdvisoryData{
				VulnID: "CVE-2345-6789",
				Data:   []byte(`{"PatchedVersions":["3.3.5"],"VulnerableVersions":["\u003c 3.3.5"]}`),
			},
		},
	}
	vulnSet := []*types.Vulnerability{
		{
			Title: "blue",
		},
		{
			Title: "orange",
		},
		{
			Title: "red",
		},
	}
	db.VulnerabilityMap["CVE-1234-5678"] = vulnSet[0]
	db.VulnerabilityMap["CVE-2345-6789"] = vulnSet[1]
	db.VulnerabilityMap["CVE-3456-7890"] = vulnSet[2]

	dt := detector.New(db)

	t.Run("detect both with 3.2.4", func(t *testing.T) {
		results, err := dt.Detect(model.PkgBundler, "somepkg", "3.2.4")
		require.NoError(t, err)
		assert.Contains(t, results, vulnSet[0])
		assert.Contains(t, results, vulnSet[1])
		assert.NotContains(t, results, vulnSet[2])
	})

	t.Run("detect CVE-2345-6789 with 3.2.5", func(t *testing.T) {
		results, err := dt.Detect(model.PkgBundler, "somepkg", "3.2.5")
		require.NoError(t, err)
		assert.NotContains(t, results, vulnSet[0])
		assert.Contains(t, results, vulnSet[1])
		assert.NotContains(t, results, vulnSet[2])
	})

	t.Run("detect no vulnerability with 3.3.5", func(t *testing.T) {
		results, err := dt.Detect(model.PkgBundler, "somepkg", "3.3.5")
		require.NoError(t, err)
		assert.Nil(t, results)
	})
}
