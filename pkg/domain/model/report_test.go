package model_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReport(t *testing.T) {
	t.Run("", func(t *testing.T) {
		oldPkgs := []*ent.PackageRecord{
			{
				Source: "x",
				Name:   "blue",
				Edges: ent.PackageRecordEdges{
					Vulnerabilities: []*ent.Vulnerability{
						{ID: "0001"},
						{ID: "0002"},
					},
				},
			},
			{
				Source: "x",
				Name:   "red",
				Edges: ent.PackageRecordEdges{
					Vulnerabilities: []*ent.Vulnerability{
						{ID: "0001"},
					},
				},
			},
		}
		newPkgs := []*ent.PackageRecord{
			{
				Source: "x",
				Name:   "blue",
				Edges: ent.PackageRecordEdges{
					Vulnerabilities: []*ent.Vulnerability{
						{ID: "0001"},
						{ID: "0002"},
					},
				},
			},
			{
				Source: "x",
				Name:   "orange",
				Edges: ent.PackageRecordEdges{
					Vulnerabilities: []*ent.Vulnerability{
						{ID: "0001"},
					},
				},
			},
		}
		changes := model.DiffVulnRecords(oldPkgs, newPkgs)
		db := model.NewVulnStatusDB([]*ent.VulnStatus{}, 1000)
		report := model.MakeReport(changes, db)
		assert.NotNil(t, report)
		{
			require.Len(t, report.Sources["x"].Added, 1)
			r := report.Sources["x"].Added[0]
			assert.Equal(t, "orange", r.Pkg.Name)
			assert.Equal(t, "x", r.Pkg.Source)
			assert.Equal(t, "0001", r.Vuln.ID)
		}
		{
			require.Len(t, report.Sources["x"].Deleted, 1)
			r := report.Sources["x"].Deleted[0]
			assert.Equal(t, "red", r.Pkg.Name)
			assert.Equal(t, "x", r.Pkg.Source)
			assert.Equal(t, "0001", r.Vuln.ID)
		}
		{
			require.Len(t, report.Sources["x"].Remained, 2)
			r := report.Sources["x"].Remained[0]
			assert.Equal(t, "blue", r.Pkg.Name)
		}
	})
}
