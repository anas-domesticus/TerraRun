package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReplace(t *testing.T) {
	P := Placeholder{Before: "foo", After: "bar"}

	tests := []struct {
		Name    string
		Input   string
		WantOut string
	}{
		{
			"empty_string",
			"",
			"",
		},
		{
			"no_placeholders",
			"some_string",
			"some_string",
		},
		{
			"one_placeholder",
			"{{foo}}_test",
			"bar_test",
		},
		{
			"more_than_one_placeholder",
			"{{foo}}_test_{{foo}}",
			"bar_test_bar",
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			out := P.Replace(tc.Input)
			assert.Equal(t, tc.WantOut, out)
		})
	}
}
