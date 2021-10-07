package internal

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindAllStacks(t *testing.T) {
	tests := []struct {
		Name    string
		Input   string
		WantErr bool
		WantOut map[int]TerraformStack
	}{
		{
			"empty_string",
			"",
			true,
			map[int]TerraformStack(nil),
		},
		{
			"not_a_directory",
			"somewhere_made_up",
			true,
			map[int]TerraformStack(nil),
		},
		{
			"empty dir",
			"testdata/empty",
			false,
			map[int]TerraformStack{},
		},
		{
			"single_dir",
			"testdata/non_tf_dir/valid_subdir",
			false,
			map[int]TerraformStack{0: {
				"testdata/non_tf_dir/valid_subdir",
				StackConfig{},
			}},
		},
		{
			"multiple_dirs",
			"testdata",
			false,
			map[int]TerraformStack{0: {
				"testdata/invalid_stack",
				StackConfig{},
			}, 1: {
				"testdata/non_tf_dir/valid_subdir",
				StackConfig{},
			}, 2: {
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

func TestFilterStacksForEnv(t *testing.T) {
	tests := []struct {
		Name    string
		Input   map[int]TerraformStack
		WantOut map[int]TerraformStack
		Environment
	}{
		{
			"empty",
			map[int]TerraformStack{},
			map[int]TerraformStack{},
			Environment{},
		},
		{
			"no_env",
			map[int]TerraformStack{0: {
				Path:   "testdata/non_tf_dir/valid_subdir",
				config: StackConfig{},
			}},
			map[int]TerraformStack{0: {
				Path:   "testdata/non_tf_dir/valid_subdir",
				config: StackConfig{},
			}},
			Environment{},
		},
		{
			"wrong_env",
			map[int]TerraformStack{0: {
				Path:   "testdata/non_tf_dir/valid_subdir",
				config: StackConfig{},
			}},
			map[int]TerraformStack{},
			Environment{Name: "foo"},
		},
		{
			"right_env",
			map[int]TerraformStack{0: {
				Path:   "testdata/non_tf_dir/valid_subdir",
				config: StackConfig{},
			}},
			map[int]TerraformStack{0: {
				Path:   "testdata/non_tf_dir/valid_subdir",
				config: StackConfig{},
			}},
			Environment{Name: "dev"},
		},
		{
			"right_env_multiple",
			map[int]TerraformStack{0: {
				Path:   "testdata/non_tf_dir/valid_subdir",
				config: StackConfig{},
			}, 1: {
				Path:   "testdata/valid_stack",
				config: StackConfig{},
			}},
			map[int]TerraformStack{0: {
				Path:   "testdata/non_tf_dir/valid_subdir",
				config: StackConfig{},
			}, 1: {
				Path:   "testdata/valid_stack",
				config: StackConfig{},
			}},
			Environment{Name: "dev"},
		},
		{
			"right_env_only_one",
			map[int]TerraformStack{0: {
				Path:   "testdata/non_tf_dir/valid_subdir",
				config: StackConfig{},
			}, 1: {
				Path:   "testdata/valid_stack",
				config: StackConfig{},
			}},
			map[int]TerraformStack{0: {
				Path:   "testdata/non_tf_dir/valid_subdir",
				config: StackConfig{},
			}},
			Environment{Name: "test"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			stacks := FilterStacksForEnv(tc.Input, tc.Environment)
			assert.Equal(t, tc.WantOut, stacks)
		})
	}
}

func TestForAllStacks(t *testing.T) {
	testFunc := func(cfg Config, stack TerraformStack) (ExecuteOutput, error) {
		return ExecuteOutput{
			Stack: stack,
		}, nil
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
		{
			"dependencies_read_test",
			"testdata/valid_stack",
			true,
			0,
			Environment{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			output, err := ForAllStacks(Config{BaseDir: tc.Input, Env: tc.Environment}, testFunc)
			if tc.WantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.WantLen, len(output))
		})
	}
}

func TestPathToKey(t *testing.T) {
	tests := []struct {
		Name      string
		InputMap  map[int]TerraformStack
		InputPath Dependency
		Expected  int
	}{
		{
			"empty",
			map[int]TerraformStack{},
			Dependency("path"),
			-1,
		},
		{
			"path_not_present",
			map[int]TerraformStack{0: {
				"testdata/invalid_stack",
				StackConfig{},
			}, 1: {
				"testdata/non_tf_dir/valid_subdir",
				StackConfig{},
			}, 2: {
				"testdata/valid_stack",
				StackConfig{},
			}},
			Dependency("path"),
			-1,
		},
		{
			"path_not_present",
			map[int]TerraformStack{0: {
				"testdata/invalid_stack",
				StackConfig{},
			}, 1: {
				"testdata/non_tf_dir/valid_subdir",
				StackConfig{},
			}, 2: {
				"testdata/valid_stack",
				StackConfig{},
			}},
			Dependency("testdata/valid_stack"),
			2,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			key := pathToKey(tc.InputMap, tc.InputPath)
			assert.Equal(t, tc.Expected, key)
		})
	}
}

func TestResolveDependencies(t *testing.T) {
	tests := []struct {
		Name      string
		InputMap  map[int]TerraformStack
		InputKey  int
		Expected  []int
		WantError bool
	}{
		{
			"completely_empty",
			map[int]TerraformStack{},
			0,
			[]int{},
			false,
		},
		{
			"empty_dependency",
			map[int]TerraformStack{0: {
				"testdata/invalid_stack",
				StackConfig{},
			}, 1: {
				"testdata/non_tf_dir/valid_subdir",
				StackConfig{},
			}, 2: {
				"testdata/valid_stack",
				StackConfig{},
			},
			},
			2,
			[]int{},
			false,
		},
		{
			"simple_dependency",
			map[int]TerraformStack{0: {
				"testdata/invalid_stack",
				StackConfig{},
			}, 1: {
				"testdata/non_tf_dir/valid_subdir",
				StackConfig{},
			}, 2: {
				"testdata/valid_stack",
				StackConfig{
					Depends: []Dependency{
						Dependency("testdata/non_tf_dir/valid_subdir"),
					},
				},
			}},
			2,
			[]int{1},
			false,
		},
		{
			"simple_multiple_dependencies",
			map[int]TerraformStack{0: {
				"testdata/invalid_stack",
				StackConfig{},
			}, 1: {
				"testdata/non_tf_dir/valid_subdir",
				StackConfig{},
			}, 2: {
				"testdata/valid_stack",
				StackConfig{
					Depends: []Dependency{
						Dependency("testdata/non_tf_dir/valid_subdir"),
						Dependency("testdata/invalid_stack"),
					},
				},
			}},
			2,
			[]int{1, 0},
			false,
		},
		{
			"recursive_dependencies",
			map[int]TerraformStack{0: {
				"testdata/invalid_stack",
				StackConfig{},
			}, 1: {
				"testdata/non_tf_dir/valid_subdir",
				StackConfig{
					Depends: []Dependency{
						Dependency("testdata/invalid_stack"),
					},
				},
			}, 2: {
				"testdata/valid_stack",
				StackConfig{
					Depends: []Dependency{
						Dependency("testdata/non_tf_dir/valid_subdir"),
					},
				},
			}},
			2,
			[]int{1, 0},
			false,
		},
		{
			"dependency_loop",
			map[int]TerraformStack{0: {
				"testdata/invalid_stack",
				StackConfig{},
			}, 1: {
				"testdata/non_tf_dir/valid_subdir",
				StackConfig{
					Depends: []Dependency{
						Dependency("testdata/valid_stack"),
					},
				},
			}, 2: {
				"testdata/valid_stack",
				StackConfig{
					Depends: []Dependency{
						Dependency("testdata/non_tf_dir/valid_subdir"),
					},
				},
			}},
			2,
			[]int{},
			true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			deps, err := resolveDependencies(tc.InputMap, tc.InputKey, nil)
			assert.Equal(t, err != nil, tc.WantError)
			assert.Equal(t, tc.Expected, deps)
		})
	}
}

func TestCheckDependencies(t *testing.T) {
	tests := []struct {
		Name      string
		InputMap  map[int]TerraformStack
		Expected  []int
		WantError bool
	}{
		{
			"completely_empty",
			map[int]TerraformStack{},
			[]int{},
			false,
		},
		{
			"dependency_loop",
			map[int]TerraformStack{0: {
				"testdata/invalid_stack",
				StackConfig{},
			}, 1: {
				"testdata/non_tf_dir/valid_subdir",
				StackConfig{
					Depends: []Dependency{
						Dependency("testdata/valid_stack"),
					},
				},
			}, 2: {
				"testdata/valid_stack",
				StackConfig{
					Depends: []Dependency{
						Dependency("testdata/non_tf_dir/valid_subdir"),
					},
				},
			}},
			[]int{},
			true,
		},
		{
			"dependency_missing",
			map[int]TerraformStack{0: {
				"testdata/invalid_stack",
				StackConfig{},
			}, 1: {
				"testdata/non_tf_dir/valid_subdir",
				StackConfig{},
			}, 2: {
				"testdata/valid_stack",
				StackConfig{
					Depends: []Dependency{
						Dependency("not_a_path"),
					},
				},
			}},
			[]int{},
			true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			err := checkDependencies(tc.InputMap)
			assert.Equal(t, err != nil, tc.WantError)
		})
	}
}

func TestDependenciesResolved(t *testing.T) {
	tests := []struct {
		Name      string
		InputDeps []Dependency
		InputMap  map[int]ExecuteOutput
		Expected  bool
	}{
		{
			"completely_empty",
			[]Dependency{},
			map[int]ExecuteOutput{},
			true,
		},
		{
			"no_dependencies",
			[]Dependency{},
			map[int]ExecuteOutput{
				0: {TerraformStack{
					"some_path",
					StackConfig{},
				},
					&Command{},
					[]byte{},
					[]byte{},
					nil,
				},
			},
			true,
		},
		{
			"resolved_dependencies",
			[]Dependency{
				Dependency("some_path"),
			},
			map[int]ExecuteOutput{
				0: {TerraformStack{
					"some_path",
					StackConfig{},
				},
					&Command{},
					[]byte{},
					[]byte{},
					nil,
				},
			},
			true,
		},
		{
			"failed_dependencies",
			[]Dependency{
				Dependency("some_path"),
			},
			map[int]ExecuteOutput{
				0: {TerraformStack{
					"some_path",
					StackConfig{},
				},
					&Command{},
					[]byte{},
					[]byte{},
					errors.New("dummy error"),
				},
			},
			false,
		},
		{
			"unresolveable_dependencies",
			[]Dependency{
				Dependency("some_path"),
			},
			map[int]ExecuteOutput{
				0: {TerraformStack{
					"some_other_path",
					StackConfig{},
				},
					&Command{},
					[]byte{},
					[]byte{},
					nil,
				},
			},
			false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, dependenciesFulfilled(tc.InputDeps, tc.InputMap), tc.Expected)
		})
	}
}
