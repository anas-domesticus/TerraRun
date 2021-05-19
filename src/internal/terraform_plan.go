package internal

import (
	"encoding/json"
	"fmt"
)

func GetTerraformPlan() Command {
	cmd := NewTerraformCommand()
	cmd.Parameters = append(cmd.Parameters, []Parameter{
		&SimpleParameter{Value: "plan"},
		&SimpleParameter{Value: "-out=plan.tfplan"},
		&SimpleParameter{Value: "-input=false"},
	}...)
	return cmd
}

func GetTerraformShowPlan() Command {
	cmd := NewTerraformCommand()
	cmd.Parameters = append(cmd.Parameters, []Parameter{
		&SimpleParameter{Value: "show"},
		&SimpleParameter{Value: "-json"},
		&SimpleParameter{Value: "plan.tfplan"},
	}...)
	return cmd
}

type PlanOutput struct {
}

func PlanWasSuccessful(output ExecuteOutput) bool {
	if output.Error != nil {
		return false
	}
	val := ValidateOutput{}
	err := json.Unmarshal(output.StdOut, &val)
	if err != nil {
		return false
	}
	return val.Valid
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
