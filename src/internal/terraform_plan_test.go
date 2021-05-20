package internal

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestTFPlan(t *testing.T) {
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
			// Check that a plan file was output
			if !tc.WantErr && !tc.WantOutputErr {
				_, err := os.Stat(filepath.Join(tc.Path, "plan.tfplan"))
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlanWasSuccessful(t *testing.T) {
	testStack := TerraformStack{Path: "testdata/valid_stack"}
	f, _ := os.Create(filepath.Join(testStack.Path, "/plan.tfplan"))
	f.Write([]byte("some stuff"))
	f.Close()

	tests := []struct {
		Name  string
		Input ExecuteOutput
		Valid bool
	}{
		{
			"success",
			ExecuteOutput{
				Stack:   testStack,
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
				Stack:   testStack,
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
	CleanTestDir()
}

func TestGetPlanJSON(t *testing.T) {
	CleanTestDir()
	tests := []struct {
		Name    string
		Path    string
		WantErr bool
		BadPlan bool
	}{
		{
			"empty_string",
			"",
			true,
			false,
		},
		{
			"not_a_directory",
			"somewhere_made_up",
			true,
			false,
		},
		{
			"no_plan",
			"./testdata/non_tf_dir",
			true,
			false,
		},
		{
			"valid_stack",
			"testdata/valid_stack",
			false,
			false,
		},
		{
			"bad_plan",
			"testdata/valid_stack",
			true,
			true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			_, _ = PlanStack(Config{BaseDir: "./"}, TerraformStack{Path: tc.Path})
			if tc.BadPlan {
				os.Remove(filepath.Join(tc.Path, "plan.tfplan"))
				f, _ := os.Create(filepath.Join(tc.Path, "plan.tfplan"))
				f.Write([]byte("some_bad_data"))
				f.Close()
			}
			json, err := GetPlanJSON(Config{BaseDir: "./"}, TerraformStack{Path: tc.Path})
			if tc.WantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			// Check that a plan file was output
			if !tc.WantErr {
				assert.True(t, len(json) > 0)
			}
		})
	}
	CleanTestDir()
}
