package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func NewTerraformCommand() Command {
	return Command{
		Binary:  "terraform",
		EnvVars: append(os.Environ(), "TF_IN_AUTOMATION=true"),
	}
}

func BuildTFStackEnvVars(stack TerraformStack) []string {
	toReturn := append([]string{}, fmt.Sprintf("TF_VAR_terrarun_stack_name=%s", filepath.Base(stack.Path)))
	return toReturn
}
