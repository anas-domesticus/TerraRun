package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecute(t *testing.T) {

	tests := []struct {
		Name       string
		Binary     string
		Params     []Parameter
		WantErr    bool
		WantStdOut string
		WantStdErr string
	}{
		{
			"missing_binary",
			"",
			nil,
			true,
			"",
			"",
		},
		{
			"non_zero_exit_code",
			"/bin/false",
			nil,
			true,
			"",
			"",
		},
		{
			"success",
			"/bin/true",
			nil,
			false,
			"",
			"",
		},
		{
			"output_test",
			"/usr/bin/echo",
			[]Parameter{&SimpleParameter{Value: "foo"}},
			false,
			"foo\n",
			"",
		},
		{
			"env_param_injection",
			"/usr/bin/echo",
			[]Parameter{&ParameterWithPlaceholders{Value: "{{Environment}}"}},
			false,
			"test\n",
			"",
		},
		{
			"stack_param_injection",
			"/usr/bin/echo",
			[]Parameter{&ParameterWithPlaceholders{Value: "{{StackRelPath}}"}},
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
			output := cmd.Execute(&cfg, stack)
			if tc.WantErr {
				assert.Error(t, output.Error)
			} else {
				assert.NoError(t, output.Error)
			}
			assert.Equal(t, tc.WantStdOut, string(output.StdOut))
			assert.Equal(t, tc.WantStdErr, string(output.StdErr))
		})
	}
}
