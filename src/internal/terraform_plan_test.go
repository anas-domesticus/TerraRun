package internal

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTFPlan(t *testing.T) {

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
			true,
			false,
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
			output, err := PlanStack(Config{BaseDir: "./"}, TerraformStack{Path: tc.Path})
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
		})
	}
}

func TestPlanWasSuccessful(t *testing.T) {
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
				StdOut:  []byte(`Some stuff`),
				StdErr:  nil,
				Error:   nil,
			},
			true,
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
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Valid, PlanWasSuccessful(tc.Input))
		})
	}
}
