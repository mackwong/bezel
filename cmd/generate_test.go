package cmd

import (
	"gotest.tools/assert"
	"testing"
)

func TestGenerate(t *testing.T) {
	err := generate("/tmp/test-generate.yaml")
	assert.NilError(t, err)
}
