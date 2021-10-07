package internal

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestTFApply(t *testing.T) {
	CleanTestDir()
	tests := []struct {
		Name          string
		Path          string
		WantErr       bool
		WantOutputErr bool
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
			false,
			true,
		},
		{
			"non_tf_dir",
			"./testdata/non_tf_dir",
			false,
			true,
		},
		{
			"valid_stack",
			"testdata/valid_stack",
			false,
			false,
		},
		{
			"bad_stack",
			"testdata/invalid_stack",
			false,
			true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := PlanStack(Config{BaseDir: "./"}, TerraformStack{Path: tc.Path})
			output, err := ApplyStack(Config{BaseDir: "./"}, TerraformStack{Path: tc.Path})

			if tc.WantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tc.WantOutputErr {
				assert.Error(t, output.Error)
			} else {
				assert.NoError(t, output.Error)
			}
			// Check that a state file was output
			if !tc.WantErr && !tc.WantOutputErr {
				_, err := os.Stat(filepath.Join(tc.Path, "terraform.tfstate"))
				assert.NoError(t, err)
				AssertEnvVarsCorrect(t, output.Command.EnvVars, true)
			}
		})
	}
	CleanTestDir()
}

func TestApplyWasSuccessful(t *testing.T) {
	tests := []struct {
		Name  string
		Input ExecuteOutput
		Valid bool
	}{
		{
			"success",
			ExecuteOutput{
				Command: &Command{},
				StdOut:  []byte(`Some stuff`),
				StdErr:  nil,
				Error:   nil,
			},
			true,
		},
		{
			"non-zero exit code",
			ExecuteOutput{
				Command: &Command{},
				StdOut:  nil,
				StdErr:  nil,
				Error:   errors.New("some error"),
			},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Valid, ApplyWasSuccessful(tc.Input))
		})
	}
	CleanTestDir()
}
