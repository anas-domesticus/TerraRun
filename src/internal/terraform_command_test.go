package internal

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func AssertEnvVarsCorrect(t *testing.T, envVars []string, stackVars bool) {
	assert.True(t, containsEnvVar(envVars, "TF_IN_AUTOMATION"))
	if stackVars {
		assert.True(t, containsEnvVar(envVars, "TF_VAR_terrarun_stack_name"))
	}
}

func containsEnvVar(s []string, e string) bool {
	for _, a := range s {
		if strings.HasPrefix(a, e) {
			return true
		}
	}
	return false
}
