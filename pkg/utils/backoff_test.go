package utils_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackoff(t *testing.T) {
	old := utils.BackoffBaseWaitTime
	utils.BackoffBaseWaitTime = 0.1 // For testing
	defer func() { utils.BackoffBaseWaitTime = old }()

	t.Run("try 5 times", func(t *testing.T) {
		n := 0
		err := utils.Backoff(10, func() (bool, error) {
			n++
			if n < 5 {
				return false, nil
			}
			return true, nil
		})
		require.NoError(t, err)
		assert.Equal(t, 5, n)
	})

	t.Run("try exceeded", func(t *testing.T) {
		n := 0
		err := utils.Backoff(10, func() (bool, error) {
			n++
			return false, nil
		})
		assert.ErrorIs(t, err, utils.ErrBackoffLimitExceeded)
		require.Equal(t, 10, n)
	})
}
