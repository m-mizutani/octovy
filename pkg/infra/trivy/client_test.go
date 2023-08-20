package trivy_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/infra/trivy"

	ftypes "github.com/aquasecurity/trivy/pkg/fanal/types"
	ttypes "github.com/aquasecurity/trivy/pkg/types"
)

func Test(t *testing.T) {
	path, ok := os.LookupEnv("TEST_TRIVY_PATH")
	if !ok {
		t.Skip("TEST_TRIVY_PATH is not set")
	}

	target := gt.R1(filepath.Abs("../../../")).NoError(t)
	t.Log(target)

	tmp := gt.R1(os.CreateTemp("", "trivy-scan-*.json")).NoError(t)
	gt.NoError(t, tmp.Close())

	client := trivy.New(path)
	gt.NoError(t, client.Run([]string{
		"fs",
		target,
		"-f", "json",
		"-o", tmp.Name(),
		"--list-all-pkgs",
	}))

	var report ttypes.Report
	body := gt.R1(os.ReadFile(tmp.Name())).NoError(t)
	gt.NoError(t, json.Unmarshal(body, &report))
	gt.V(t, report.SchemaVersion).Equal(2)
	gt.A(t, report.Results).Length(1).At(0, func(t testing.TB, v ttypes.Result) {
		gt.A(t, v.Packages).Any(func(t testing.TB, v ftypes.Package) bool {
			return v.Name == "github.com/m-mizutani/goerr"
		})
	})
}
