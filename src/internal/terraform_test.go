package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTFValidate(t *testing.T) {
	initCmd := GetTerraformInit()
	initCmd.Parameters = append(initCmd.Parameters, &SimpleParameter{Value: "-backend=false"})
	output := initCmd.Execute(Config{BaseDir: "./"}, TerraformStack{Path: "testdata/valid_stack"})
	assert.NoError(t, output.Error)
}
