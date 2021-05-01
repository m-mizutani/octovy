package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStepDownDirectory(t *testing.T) {
	assert.Equal(t, "blue/Gemfile.lock", stepDownDirectory("root/blue/Gemfile.lock"))
	assert.Equal(t, "blue/Gemfile.lock", stepDownDirectory("./root/blue/Gemfile.lock"))
	assert.Equal(t, "blue/green/Gemfile.lock", stepDownDirectory("/root/blue/green/Gemfile.lock"))
}
