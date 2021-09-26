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

func TestNode(t *testing.T) {
	_, db := trivydb.NewMock()
	db.AdvisoryMap["GitHub Security Advisory Npm"] = map[string][]*model.AdvisoryData{
		"somepkg": {
			&model.AdvisoryData{
				VulnID: "CVE-1234-5678",
				Data:   []byte(`{"PatchedVersions":["1.2.0"],"VulnerableVersions":["\u003c 1.2.0"]}`),
			},
		},
	}
	db.AdvisoryMap["nodejs-security-wg"] = map[string][]*model.AdvisoryData{}

	db.VulnerabilityMap["CVE-1234-5678"] = &types.Vulnerability{
		Title: "blue",
	}

	dt := detector.New(db)

	t.Run("1.1.9 is vulnerable", func(t *testing.T) {
		results, err := dt.Detect(octovy.PkgNPM, "somepkg", "1.1.9")
		require.NoError(t, err)
		assert.Equal(t, 1, len(results))
	})

	t.Run("1.2.0 is not vulnerable", func(t *testing.T) {
		results, err := dt.Detect(octovy.PkgNPM, "somepkg", "1.2.0")
		require.NoError(t, err)
		assert.Nil(t, results)
	})

	t.Run("1.2.1 is not also vulnerable", func(t *testing.T) {
		results, err := dt.Detect(octovy.PkgNPM, "somepkg", "1.2.1")
		require.NoError(t, err)
		assert.Nil(t, results)
	})
}
