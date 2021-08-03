package model_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestConfigRules(t *testing.T) {
	t.Run("PR comment trigger", func(t *testing.T) {
		cfg := model.Config{
			RulePullReqCommentTriggers: "opened|synchronize|reopened",
		}
		assert.True(t, cfg.ShouldCommentPR("opened"))
		assert.True(t, cfg.ShouldCommentPR("synchronize"))
		assert.True(t, cfg.ShouldCommentPR("reopened"))
		assert.False(t, cfg.ShouldCommentPR("ready_for_review"))
		assert.False(t, cfg.ShouldCommentPR(""))
	})

	t.Run("PR comment trigger", func(t *testing.T) {
		cfg1 := model.Config{
			RuleFailCheckIfVuln: "YES",
		}
		assert.True(t, cfg1.ShouldFailIfVuln())

		cfg2 := model.Config{
			RuleFailCheckIfVuln: "",
		}
		assert.False(t, cfg2.ShouldFailIfVuln())
	})
}
