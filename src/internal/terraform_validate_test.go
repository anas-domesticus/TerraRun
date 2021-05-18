package internal

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTFValidate(t *testing.T) {

	tests := []struct {
		Name          string
		Path          string
		WantErr       bool
		WantOutputErr bool
		Valid         bool
	}{
		{
			"empty_string",
			"",
			false,
			false,
			true,
		},
		{
			"not_a_directory",
			"somewhere_made_up",
			true,
			false,
			false,
		},
		{
			"non_tf_dir",
			"./testdata/non_tf_dir",
			false,
			false,
			true,
		},
		{
			"valid_stack",
			"testdata/valid_stack",
			false,
			false,
			true,
		},
		{
			"bad_stack",
			"testdata/invalid_stack",
			true,
			false,
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			output, err := ValidateStack(Config{BaseDir: "./"}, TerraformStack{Path: tc.Path})
			if tc.WantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			valOutput := ValidateOutput{}
			if tc.WantOutputErr {
				assert.Error(t, output.Error)
			} else {
				assert.NoError(t, output.Error)
			}
			_ = json.Unmarshal(output.StdOut, &valOutput)
			assert.Equal(t, tc.Valid, valOutput.Valid)
		})
	}
}
