package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindAllStacks(t *testing.T) {
	tests := []struct {
		Name    string
		Input   string
		WantErr bool
		WantOut []TerraformStack
	}{
		{
			"empty_string",
			"",
			true,
			nil,
		},
		{
			"not_a_directory",
			"somewhere_made_up",
			true,
			nil,
		},
		{
			"empty dir",
			"testdata/empty",
			false,
			nil,
		},
		{
			"single_dir",
			"testdata/non_tf_dir/valid_subdir",
			false,
			[]TerraformStack{{
				"testdata/non_tf_dir/valid_subdir",
			}},
		},
		{
			"multiple_dirs",
			"testdata",
			false,
			[]TerraformStack{{
				"testdata/non_tf_dir/valid_subdir",
			}, {
				"testdata/valid_stack",
			}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			stacks, err := FindAllStacks(tc.Input)
			if tc.WantErr {
				assert.Error(t, err)
			}
			assert.Equal(t, tc.WantOut, stacks)
		})
	}
}

func TestIsTerraformStack(t *testing.T) {
	assert.False(t, IsTerraformStack("testdata/non_tf_dir"))
	assert.True(t, IsTerraformStack("testdata/valid_stack"))
}

func TestShouldRunForEnv(t *testing.T) {
	stack := TerraformStack{Path: "testdata/valid_stack"}
	assert.False(t, stack.ShouldRunForEnv(Environment{Name: "prod"}))
	assert.True(t, stack.ShouldRunForEnv(Environment{Name: "dev"}))
}

func TestForAllStacks(t *testing.T) {
	testFunc := func(cfg Config, stack TerraformStack) (ExecuteOutput, error) {
		return ExecuteOutput{}, nil
	}
	tests := []struct {
		Name    string
		Input   string
		WantErr bool
		WantLen int
	}{
		{
			"empty_string",
			"",
			true,
			0,
		},
		{
			"not_a_directory",
			"somewhere_made_up",
			true,
			0,
		},
		{
			"empty dir",
			"testdata/empty",
			false,
			0,
		},
		{
			"single_dir",
			"testdata/non_tf_dir/valid_subdir",
			false,
			1,
		},
		{
			"multiple_dirs",
			"testdata",
			false,
			2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			output, err := ForAllStacks(Config{BaseDir: tc.Input}, testFunc)
			if tc.WantErr {
				assert.Error(t, err)
			}
			assert.Equal(t, tc.WantLen, len(output))
		})
	}
}
