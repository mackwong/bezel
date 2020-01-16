package cmd

import (
	"fmt"
	"gitlab.bj.sensetime.com/diamond/bezel/pkg/model"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
	"io/ioutil"
	"testing"
)

func TestGetGlobalConfigByConfig(t *testing.T) {
	gc, err := GetGlobalConfigByConfig("../tests/no-master-ip.yaml")
	assert.NilError(t, err)

	o, _ := yaml.Marshal(gc)
	fmt.Println(string(o))

	c, err := ioutil.ReadFile("../tests/gc.yaml")
	assert.NilError(t, err)
	target := model.GlobalConfig{}
	err = yaml.Unmarshal(c, &target)
	assert.NilError(t, err)
	assert.DeepEqual(t, *gc, target)
}

func TestSplitFromGlobalConfig(t *testing.T) {
	c, err := ioutil.ReadFile("../tests/edge-config.yaml")
	assert.NilError(t, err)
	target := model.GlobalConfig{}
	err = yaml.Unmarshal(c, &target)
	assert.NilError(t, err)
	err = SplitFromGlobalConfig(&target, "/tmp")
	assert.NilError(t, err)
}
