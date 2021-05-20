package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetPlanFilename() string {
	return "plan.tfplan"
}

func GetTerraformPlan() Command {
	cmd := NewTerraformCommand()
	cmd.Parameters = append(cmd.Parameters, []Parameter{
		&SimpleParameter{Value: "plan"},
		&SimpleParameter{Value: fmt.Sprintf("-out=%s", GetPlanFilename())},
		&SimpleParameter{Value: "-input=false"},
	}...)
	return cmd
}

type PlanOutput struct {
}

func PlanWasSuccessful(output ExecuteOutput) bool {
	return output.Error == nil && PlanPresent(output.Stack)
}

func PlanStack(config Config, stack TerraformStack) (ExecuteOutput, error) {
	fmt.Printf("Planning %s...\n", stack.Path)
	output, err := InitStack(config, stack)
	if err != nil {
		// TODO: Add detail to error
		return output, err
	}
	planCmd := GetTerraformPlan()
	output, err = planCmd.Execute(config, stack)
	return output, err
}

func PlanPresent(stack TerraformStack) bool {
	_, err := os.Stat(filepath.Join(stack.Path, GetPlanFilename()))
	return err == nil
}
