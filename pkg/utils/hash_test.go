package utils_test

import (
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/utils"
)

func TestHashBranch(t *testing.T) {
	v := utils.HashBranch("test")
	gt.Equal(t, v, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08")
}
