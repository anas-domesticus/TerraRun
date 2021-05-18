package internal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExecute(t *testing.T) {

	tests := []struct {
		Name          string
		Binary        string
		Params        []Parameter
		WantErr       bool
		WantOutputErr bool
		WantStdOut    string
		WantStdErr    string
	}{
		{
			"missing_binary",
			"",
			nil,
			true,
			false,
			"",
			"",
		},
		{
			"non_zero_exit_code",
			"/bin/false",
			nil,
			false,
			true,
			"",
			"",
		},
		{
			"success",
			"/bin/true",
			nil,
			false,
			false,
			"",
			"",
		},
		{
			"output_test",
			"/bin/echo",
			[]Parameter{&SimpleParameter{Value: "foo"}},
			false,
			false,
			"foo\n",
			"",
		},
		{
			"env_param_injection",
			"/bin/echo",
			[]Parameter{&ParameterWithPlaceholders{Value: "{{Environment}}"}},
			false,
			false,
			"test\n",
			"",
		},
		{
			"stack_param_injection",
			"/bin/echo",
			[]Parameter{&ParameterWithPlaceholders{Value: "{{StackRelPath}}"}},
			false,
			false,
			"testdata/valid_stack\n",
			"",
		},
	}
	cfg := Config{Env: Environment{Name: "test"}}
	stack := TerraformStack{Path: "testdata/valid_stack"}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := Command{Parameters: tc.Params, Binary: tc.Binary}
			output, err := cmd.Execute(cfg, stack, os.Environ())
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
			assert.Equal(t, tc.WantStdOut, string(output.StdOut))
			assert.Equal(t, tc.WantStdErr, string(output.StdErr))
		})
	}
}
