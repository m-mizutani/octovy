package model_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/stretchr/testify/assert"
)

func TestVulnStatusDB(t *testing.T) {
	t.Run("snooze disabled matched entry", func(t *testing.T) {
		db := model.NewVulnStatusDB([]*ent.VulnStatus{
			{
				Status:    types.StatusSnoozed,
				Source:    "x",
				PkgName:   "blue",
				VulnID:    "CVE-2000",
				ExpiresAt: 2000,
			},
		}, 1000)

		cases := []struct {
			title    string
			val      *model.VulnRecord
			expected bool
		}{
			{
				title: "disable matched entry",
				val: &model.VulnRecord{
					Pkg: &ent.PackageRecord{
						Source: "x",
						Name:   "blue",
					},
					Vuln: &ent.Vulnerability{
						ID: "CVE-2000",
					},
				},
				expected: false,
			},
			{
				title: "not affected by different source",
				val: &model.VulnRecord{
					Pkg: &ent.PackageRecord{
						Source: "y",
						Name:   "blue",
					},
					Vuln: &ent.Vulnerability{
						ID: "CVE-2000",
					},
				},
				expected: true,
			},
			{
				title: "not affected by different package name",
				val: &model.VulnRecord{
					Pkg: &ent.PackageRecord{
						Source: "x",
						Name:   "orange",
					},
					Vuln: &ent.Vulnerability{
						ID: "CVE-2000",
					},
				},
				expected: true,
			},
			{
				title: "not affected by different vuln",
				val: &model.VulnRecord{
					Pkg: &ent.PackageRecord{
						Source: "x",
						Name:   "blue",
					},
					Vuln: &ent.Vulnerability{
						ID: "CVE-1000",
					},
				},
				expected: true,
			},
		}

		for _, c := range cases {
			t.Run(c.title, func(t *testing.T) {
				assert.Equal(t, c.expected, db.IsQualified(c.val))

			})
		}
	})

	t.Run("snooze does not work when expired", func(t *testing.T) {
		db := model.NewVulnStatusDB([]*ent.VulnStatus{
			{
				Status:    types.StatusSnoozed,
				Source:    "x",
				PkgName:   "blue",
				VulnID:    "CVE-2000",
				ExpiresAt: 2000,
			},
		}, 2001)
		assert.True(t, db.IsQualified(&model.VulnRecord{
			Pkg: &ent.PackageRecord{
				Source: "x",
				Name:   "blue",
			},
			Vuln: &ent.Vulnerability{
				ID: "CVE-2000",
			},
		}))
	})

	t.Run("none does not disable", func(t *testing.T) {
		db := model.NewVulnStatusDB([]*ent.VulnStatus{
			{
				Status:    types.StatusNone,
				Source:    "x",
				PkgName:   "blue",
				VulnID:    "CVE-2000",
				ExpiresAt: 2000,
			},
		}, 1000)
		assert.True(t, db.IsQualified(&model.VulnRecord{
			Pkg: &ent.PackageRecord{
				Source: "x",
				Name:   "blue",
			},
			Vuln: &ent.Vulnerability{
				ID: "CVE-2000",
			},
		}))
	})

	t.Run("mitigated/notaffected disable even if having expires", func(t *testing.T) {
		rec := &model.VulnRecord{
			Pkg: &ent.PackageRecord{
				Source: "x",
				Name:   "blue",
			},
			Vuln: &ent.Vulnerability{
				ID: "CVE-2000",
			},
		}

		t.Run("mitigated", func(t *testing.T) {
			db := model.NewVulnStatusDB([]*ent.VulnStatus{
				{
					Status:    types.StatusMitigated,
					Source:    "x",
					PkgName:   "blue",
					VulnID:    "CVE-2000",
					ExpiresAt: 1,
				},
			}, 1000)
			assert.False(t, db.IsQualified(rec))
		})

		t.Run("mitigated", func(t *testing.T) {
			db := model.NewVulnStatusDB([]*ent.VulnStatus{
				{
					Status:    types.StatusMitigated,
					Source:    "x",
					PkgName:   "blue",
					VulnID:    "CVE-2000",
					ExpiresAt: 1,
				},
			}, 1000)
			assert.False(t, db.IsQualified(rec))
		})

	})
}
