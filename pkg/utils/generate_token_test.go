package utils_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	t1 := utils.GenerateToken(24)
	t2 := utils.GenerateToken(24)

	assert.Len(t, t1, 24)
	assert.Len(t, t2, 24)
	assert.NotEqual(t, t1, t2)
}
