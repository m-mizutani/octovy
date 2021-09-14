package detector_test

import (
	"testing"

	"github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	octovy "github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/trivydb"
	"github.com/m-mizutani/octovy/pkg/usecase/detector"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoMod(t *testing.T) {
	_, db := trivydb.NewMock()
	db.AdvisoryMap["go::GitLab Advisory Database"] = map[string][]*model.AdvisoryData{
		"somepkg": {
			&model.AdvisoryData{
				VulnID: "CVE-1234-5678",
				Data:   []byte(`{"PatchedVersions":["v1.5.5"],"VulnerableVersions":["\u003e=v1.5.0 \u003c=v1.5.4"]}`),
			},
		},
	}

	vulnSet := []*model.Vulnerability{
		{
			VulnID: "CVE-1234-5678",
			Detail: types.Vulnerability{
				Title: "blue",
			},
		},
	}
	db.VulnerabilityMap["CVE-1234-5678"] = &vulnSet[0].Detail

	dt := detector.New(db)

	t.Run("1.5.1 is vulnerable", func(t *testing.T) {
		results, err := dt.Detect(octovy.PkgGoModule, "somepkg", "1.5.1")
		require.NoError(t, err)
		assert.Contains(t, results, vulnSet[0])
	})

	t.Run("1.5.5 is not vulnerable", func(t *testing.T) {
		results, err := dt.Detect(octovy.PkgGoModule, "somepkg", "1.5.5")
		require.NoError(t, err)
		assert.Nil(t, results)
	})

	t.Run("1.5.6 is not also vulnerable", func(t *testing.T) {
		results, err := dt.Detect(octovy.PkgGoModule, "somepkg", "1.5.6")
		require.NoError(t, err)
		assert.Nil(t, results)
	})
}
