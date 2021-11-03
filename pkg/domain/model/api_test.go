package model_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestRequestSeverity(t *testing.T) {
	t.Run("empty color is valid", func(t *testing.T) {
		req := model.RequestSeverity{
			Label: "a",
		}
		assert.NoError(t, req.IsValid())
		assert.NotEmpty(t, req.Color)
	})

	t.Run("invalid color", func(t *testing.T) {
		invalid := []string{
			"abcdef",
			"#ABCDEG",
			" #ABCDEF",
			"#ABCDEF ",
		}

		for _, c := range invalid {
			t.Run("test with "+c, func(t *testing.T) {
				req := model.RequestSeverity{
					Label: "a",
					Color: c,
				}
				assert.Error(t, req.IsValid())
			})
		}
	})

}
