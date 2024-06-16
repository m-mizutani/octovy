package model_test

import (
	_ "embed"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

//go:embed testdata/config/ignore.cue
var testConfigIgnoreCue []byte

func TestIgnoreConfig(t *testing.T) {
	cfg, err := model.BuildConfig(testConfigIgnoreCue)
	gt.NoError(t, err)
	gt.A(t, cfg.IgnoreTargets).Length(2).
		At(0, func(t testing.TB, v model.IgnoreTarget) {
			gt.Equal(t, v.File, "test.data")
			gt.A(t, v.Vulns).Length(1).At(0, func(t testing.TB, v model.IgnoreVuln) {
				gt.Equal(t, v.ID, "CVE-2017-9999")
				gt.Equal(t, v.Description, "This is test data")
				gt.Equal(t, v.ExpiresAt.Year(), 2018)
			})
		}).
		At(1, func(t testing.TB, v model.IgnoreTarget) {
			gt.Equal(t, v.File, "test2.data")
			gt.A(t, v.Vulns).Length(2).
				At(0, func(t testing.TB, v model.IgnoreVuln) {
					gt.Equal(t, v.ID, "CVE-2017-11423")
					gt.Equal(t, v.ExpiresAt.Year(), 2022)
				}).
				At(1, func(t testing.TB, v model.IgnoreVuln) {
					gt.Equal(t, v.ID, "CVE-2023-11423")
					gt.Equal(t, v.ExpiresAt.Year(), 2023)
				})
		})
}
