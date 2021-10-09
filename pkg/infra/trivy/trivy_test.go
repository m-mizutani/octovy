package trivy_test

import (
	"os"
	"testing"

	"github.com/m-mizutani/octovy/pkg/infra/trivy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrivy(t *testing.T) {
	if _, ok := os.LookupEnv("WITH_TRIVY_COMMAND"); !ok {
		t.Skip()
	}

	cmd := trivy.New(trivy.DefaultName)
	result, err := cmd.Scan("./testdata")
	require.NoError(t, err)
	require.Len(t, result.Results, 1)
	assert.Equal(t, "Gemfile.lock", result.Results[0].Target)
	assert.Len(t, result.Results[0].Packages, 42)

}
