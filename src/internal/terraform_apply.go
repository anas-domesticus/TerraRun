package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetTerraformApply() Command {
	cmd := NewTerraformCommand()
	cmd.Parameters = append(cmd.Parameters, []Parameter{
		&SimpleParameter{Value: "apply"},
		&SimpleParameter{Value: "plan.tfplan"},
	}...)
	return cmd
}

type ApplyOutput struct {
}

func ApplyWasSuccessful(output ExecuteOutput) bool {
	return output.Error == nil
}

func ApplyStack(config Config, stack TerraformStack) (ExecuteOutput, error) {
	if !PlanPresent(stack) {
		return ExecuteOutput{}, fmt.Errorf("plan file missing for: %s", stack.Path)
	}
	fmt.Printf("Applying %s...\n", stack.Path)
	output, err := InitStack(config, stack)
	if err != nil {
		// TODO: Add detail to error
		return output, err
	}
	applyCmd := GetTerraformApply()
	output, err = applyCmd.Execute(config, stack)
	return output, err
}

func PlanPresent(stack TerraformStack) bool {
	_, err := os.Stat(filepath.Join(stack.Path, "plan.tfplan"))
	return err == nil
}
