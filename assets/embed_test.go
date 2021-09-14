package assets_test

import (
	"os"
	"testing"

	"github.com/m-mizutani/octovy/assets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssets(t *testing.T) {
	t.Run("index.html loadable", func(t *testing.T) {
		raw, err := assets.Assets().ReadFile("dist/index.html")
		require.NoError(t, err)
		assert.Contains(t, string(raw), "<html>")
	})

	t.Run("bundle.js loadable", func(t *testing.T) {
		if _, ok := os.LookupEnv("GITHUB_WORKFLOW"); ok {
			t.Skip("bundle.js is not generated in GitHub Actions")
		}
		raw, err := assets.Assets().ReadFile("dist/bundle.js")
		require.NoError(t, err)
		assert.Contains(t, string(raw), "__webpack_modules__")
	})

	t.Run("files in out of ./dist/ directory can not be loaded", func(t *testing.T) {
		_, err := assets.Assets().ReadFile("assets.go")
		require.Error(t, err)
	})
}
