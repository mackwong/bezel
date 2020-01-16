package cmd

import (
	"gotest.tools/assert"
	"testing"
)

func TestParseCmd(t *testing.T) {
	sc, err := loadSubConfig("../tests/sub-edge-config-master-10.4.72.1.yaml")
	assert.NilError(t, err)
	err = parseAllTemplates("../tests/templates", "/tmp", sc)
	assert.NilError(t, err)
}
