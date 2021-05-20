package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			_, _ = PlanStack(Config{BaseDir: "./"}, TerraformStack{Path: tc.Path})
			json, err := getPlanJSON(Config{BaseDir: "./"}, TerraformStack{Path: tc.Path})
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
