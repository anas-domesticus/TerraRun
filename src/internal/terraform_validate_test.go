package internal

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func CleanTestDir() {
	remove := []string{
		"testdata/valid_stack/plan.tfplan",
		"testdata/valid_stack/.terraform",
		"testdata/valid_stack/.terraform.lock.hcl",
		"testdata/valid_stack/terraform.tfstate",
		"testdata/non_tf_dir/valid_subdir/plan.tfplan",
		"testdata/non_tf_dir/valid_subdir/.terraform",
		"testdata/non_tf_dir/valid_subdir/.terraform.lock.hcl",
		"testdata/non_tf_dir/terraform.tfstate",
	}
	for _, v := range remove {
		os.RemoveAll(v)
	}
}

func TestTFValidate(t *testing.T) {
	CleanTestDir()
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
			false,
			true,
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
			if !tc.WantErr && !tc.WantOutputErr {
				AssertEnvVarsCorrect(t, output.Command.EnvVars, false)
			}
			_ = json.Unmarshal(output.StdOut, &valOutput)
			assert.Equal(t, tc.Valid, valOutput.Valid)
		})
	}
	CleanTestDir()
}

func TestValidateWasSuccessful(t *testing.T) {
	tests := []struct {
		Name  string
		Input ExecuteOutput
		Valid bool
	}{
		{
			"success",
			ExecuteOutput{
				Stack:   TerraformStack{},
				Command: &Command{},
				StdOut:  []byte(`{"valid":true,"error_count":0,"warning_count":0,"diagnostics":[]}`),
				StdErr:  nil,
				Error:   nil,
			},
			true,
		},
		{
			"failure",
			ExecuteOutput{
				Stack:   TerraformStack{},
				Command: &Command{},
				StdOut:  []byte(`{"valid":false,"error_count":1,"warning_count":0,"diagnostics":[]}`),
				StdErr:  nil,
				Error:   nil,
			},
			false,
		},
		{
			"non-zero exit code",
			ExecuteOutput{
				Stack:   TerraformStack{},
				Command: &Command{},
				StdOut:  nil,
				StdErr:  nil,
				Error:   errors.New("some error"),
			},
			false,
		},
		{
			"wrong json",
			ExecuteOutput{
				Stack:   TerraformStack{},
				Command: &Command{},
				StdOut:  []byte(`{"totally_unrelated_json":"yes"}`),
				StdErr:  nil,
				Error:   nil,
			},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Valid, ValidateWasSuccessful(tc.Input))
		})
	}
}
