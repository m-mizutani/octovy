package trivy_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/infra/trivy"

	trivy_model "github.com/m-mizutani/octovy/pkg/domain/model/trivy"
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
	ctx := context.Background()
	gt.NoError(t, client.Run(ctx, []string{
		"fs",
		target,
		"-f", "json",
		"-o", tmp.Name(),
		"--list-all-pkgs",
	}))

	var report trivy_model.Report
	body := gt.R1(os.ReadFile(tmp.Name())).NoError(t)
	gt.NoError(t, json.Unmarshal(body, &report))
	gt.V(t, report.SchemaVersion).Equal(2)
	gt.A(t, report.Results).Longer(0).Any(func(v trivy_model.Result) bool {
		if v.Target == "go.mod" {
			gt.A(t, v.Packages).Any(func(v trivy_model.Package) bool {
				return v.Name == "github.com/m-mizutani/goerr"
			})
		}

		return v.Target == "go.mod"
	})
}
