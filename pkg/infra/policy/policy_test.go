package policy_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/policy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckRule(t *testing.T) {
	inv := model.NewScanReport(&ent.Scan{
		Edges: ent.ScanEdges{
			Packages: []*ent.PackageRecord{
				{
					Type:    "gomod",
					Source:  "go.sum",
					Name:    "github.com/dgrijalva/jwt-go",
					Version: "3.2.0+incompatible",
					VulnIds: []string{"CVE-2020-26160"},
					Edges: ent.PackageRecordEdges{
						Vulnerabilities: []*ent.Vulnerability{
							{
								ID: "CVE-2020-26160",
								Edges: ent.VulnerabilityEdges{
									CustomSeverity: &ent.Severity{
										Label: "high",
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil, 0)

	ctx := model.NewContext()
	t.Run("always success", func(t *testing.T) {
		check, err := policy.NewCheck(`package octovy.check
conclusion = "success"
`)
		require.NoError(t, err)

		result, err := check.Result(ctx, inv)
		require.NoError(t, err)
		assert.Equal(t, "success", result.Conclusion)
		assert.Empty(t, result.Messages)
	})

	t.Run("failure if severity is high", func(t *testing.T) {
		check, err := policy.NewCheck(`package octovy.check
default conclusion = "success"
conclusion = "failure" {
    vuln := input.sources[_].packages[_].vulnerabilities[_]
    vuln.custom_severity == "high"
}
`)
		require.NoError(t, err)

		result, err := check.Result(ctx, inv)
		require.NoError(t, err)
		assert.Equal(t, "failure", result.Conclusion)
		assert.Empty(t, result.Messages)
	})

	t.Run("message is set by msg", func(t *testing.T) {
		check, err := policy.NewCheck(`package octovy.check
		conclusion = "success"
		messages = {"blue"}
`)
		require.NoError(t, err)

		result, err := check.Result(ctx, inv)
		require.NoError(t, err)
		assert.Equal(t, "success", result.Conclusion)
		assert.Contains(t, result.Messages, "blue")
	})

	t.Run("err if invalid rego", func(t *testing.T) {
		_, err := policy.NewCheck(`package octovy.check
		default result = "success"
		conclusion = "failure" {
			vuln := input.sources[_].packages[_].vulnerabilities[_]
			vuln.custom_severity == "high"
		`) // missing tail bracket
		require.Error(t, err)
	})

	t.Run("err if missing result field", func(t *testing.T) {
		check, err := policy.NewCheck(`package octovy.check
		xxx = "failure" {
			vuln := input.sources[_].packages[_].vulnerabilities[_]
			vuln.custom_severity == "high"
		}
		`)
		require.NoError(t, err)

		result, err := check.Result(ctx, inv)
		require.Error(t, err)
		assert.ErrorIs(t, err, model.ErrInvalidPolicyResult)
		assert.Nil(t, result)
	})

	t.Run("err if package name is not octovy.check", func(t *testing.T) {
		check, err := policy.NewCheck(`package octovy.cheeeeeeeeeeeeeeeeeeek
		default result = "success"
		conclusion = "failure" {
			vuln := input.sources[_].packages[_].vulnerabilities[_]
			vuln.custom_severity == "high"
		}
		`)
		require.NoError(t, err)

		result, err := check.Result(ctx, inv)
		require.ErrorIs(t, err, model.ErrInvalidPolicyResult)
		assert.Nil(t, result)
	})

	t.Run("err if missing package", func(t *testing.T) {
		_, err := policy.NewCheck(`
		default conclusion = "success"
		conclusion = "failure" {
			vuln := input.sources[_].packages[_].vulnerabilities[_]
			vuln.custom_severity == "high"
		}
		`)
		require.Error(t, err)
	})

	t.Run("err if result is not string", func(t *testing.T) {
		check, err := policy.NewCheck(`package octovy.check
		conclusion = 0
`)
		require.NoError(t, err)

		result, err := check.Result(ctx, inv)
		require.Error(t, err)
		assert.ErrorIs(t, err, model.ErrInvalidPolicyResult)
		assert.Nil(t, result)
	})
}
