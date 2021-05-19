package internal

import (
	"fmt"
)

func GetTerraformApply() Command {
	cmd := NewTerraformCommand()
	cmd.Parameters = append(cmd.Parameters, []Parameter{
		&SimpleParameter{Value: "apply"},
		&SimpleParameter{Value: "plan.tfplan"},
		&SimpleParameter{Value: "-input=false"},
	}...)
	return cmd
}

type ApplyOutput struct {
}

func ApplyWasSuccessful(output ExecuteOutput) bool {
	return output.Error == nil
}

func ApplyStack(config Config, stack TerraformStack) (ExecuteOutput, error) {
	fmt.Printf("Applying %s...", stack.Path)
	output, err := InitStack(config, stack)
	if err != nil {
		// TODO: Add detail to error
		return output, err
	}
	applyCmd := GetTerraformApply()
	output, err = applyCmd.Execute(config, stack)
	return output, err
}
