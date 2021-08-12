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
		WantOut []*TerraformStack
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
			[]*TerraformStack{{
				"testdata/non_tf_dir/valid_subdir",
				StackConfig{},
			}},
		},
		{
			"multiple_dirs",
			"testdata",
			false,
			[]*TerraformStack{{
				"testdata/invalid_stack",
				StackConfig{},
			}, {
				"testdata/non_tf_dir/valid_subdir",
				StackConfig{},
			}, {
				"testdata/valid_stack",
				StackConfig{
					Depends: []Dependency{
						Dependency("testdata/non_tf_dir/valid_subdir"),
					},
				},
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

func TestGetStackConfig(t *testing.T) {
	tests := []struct {
		Name    string
		Path    string
		WantErr bool
		WantOut StackConfig
	}{
		{
			"empty_string",
			"",
			true,
			StackConfig{},
		},
		{
			"not_a_directory",
			"somewhere_made_up",
			true,
			StackConfig{},
		},
		{
			"empty dir",
			"testdata/empty",
			false,
			StackConfig{},
		},
		{
			"no_config",
			"testdata/non_tf_dir/valid_subdir",
			false,
			StackConfig{},
		},
		{
			"valid_config",
			"testdata/valid_stack",
			false,
			StackConfig{
				Depends: []Dependency{
					Dependency("testdata/non_tf_dir/valid_subdir"),
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := getStackConfig(tc.Path)
			if tc.WantErr {
				assert.Error(t, err)
			}
			assert.Equal(t, tc.WantOut, result)
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
	assert.True(t, stack.ShouldRunForEnv(Environment{}))
}

func TestForAllStacks(t *testing.T) {
	testFunc := func(cfg Config, stack *TerraformStack) (ExecuteOutput, error) {
		return ExecuteOutput{}, nil
	}
	tests := []struct {
		Name    string
		Input   string
		WantErr bool
		WantLen int
		Environment
	}{
		{
			"empty_string",
			"",
			true,
			0,
			Environment{},
		},
		{
			"not_a_directory",
			"somewhere_made_up",
			true,
			0,
			Environment{},
		},
		{
			"empty dir",
			"testdata/empty",
			false,
			0,
			Environment{},
		},
		{
			"single_dir",
			"testdata/non_tf_dir/valid_subdir",
			false,
			1,
			Environment{},
		},
		{
			"multiple_dirs",
			"testdata",
			false,
			3,
			Environment{},
		},
		{
			"env_filter_test",
			"testdata",
			false,
			1,
			Environment{Name: "test"},
		},
		{
			"env_filter_dev",
			"testdata",
			false,
			2,
			Environment{Name: "dev"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			output, err := ForAllStacks(Config{BaseDir: tc.Input, Env: tc.Environment}, testFunc)
			if tc.WantErr {
				assert.Error(t, err)
			}
			assert.Equal(t, tc.WantLen, len(output))
		})
	}
}
