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

func PlanWasSuccessful(output ExecuteOutput) bool {
	return output.Error == nil && PlanPresent(output.Stack)
}

func PlanStack(config Config, stack TerraformStack) (ExecuteOutput, error) {
	output, err := InitStack(config, stack)
	if err != nil || output.Error != nil {
		// TODO: Add detail to error
		return output, err
	}
	planCmd := GetTerraformPlan()
	if config.Env.Name != "" {
		planCmd.Parameters = append(planCmd.Parameters, &SimpleParameter{Value: fmt.Sprintf("-var-file=env-%s.tfvars", config.Env.Name)})
	}
	planCmd.EnvVars = append(planCmd.EnvVars, BuildTFStackEnvVars(stack)...)
	output, err = planCmd.Execute(config, stack)
	return output, err
}

func PlanPresent(stack TerraformStack) bool {
	_, err := os.Stat(filepath.Join(stack.Path, GetPlanFilename()))
	return err == nil
}
