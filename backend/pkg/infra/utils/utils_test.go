package utils_test

import (
	"testing"

	"github.com/m-mizutani/octovy/backend/pkg/infra/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	s1 := utils.GenerateToken(32)
	s2 := utils.GenerateToken(32)
	assert.NotEqual(t, s1, s2)
	t.Log("Example: ", s1)
}
