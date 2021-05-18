package internal

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTFValidate(t *testing.T) {

	tests := []struct {
		Name    string
		Path    string
		WantErr bool
		Valid   bool
	}{
		{
			"empty_string",
			"",
			false,
			true,
		},
		{
			"not_a_directory",
			"somewhere_made_up",
			true,
			false,
		},
		{
			"non_tf_dir",
			"testdata/non_tf_dir",
			false,
			true,
		},
		{
			"valid_stack",
			"testdata//valid_stack",
			false,
			true,
		},
		{
			"bad_stack",
			"testdata/invalid_stack",
			true,
			false,
		},
	}

	initCmd := GetTerraformInit()
	validateCmd := GetTerraformValidate()

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			// init stage
			initCmd.Parameters = append(initCmd.Parameters, &SimpleParameter{Value: "-backend=false"})
			output := initCmd.Execute(Config{BaseDir: "./"}, TerraformStack{Path: tc.Path})

			// & validate
			valOutput := ValidateOutput{}
			output = validateCmd.Execute(Config{BaseDir: "./"}, TerraformStack{Path: tc.Path})
			if tc.WantErr {
				assert.Error(t, output.Error)
			} else {
				assert.NoError(t, output.Error)
			}
			_ = json.Unmarshal(output.StdOut, &valOutput)
			assert.Equal(t, tc.Valid, valOutput.Valid)
		})
	}
}
