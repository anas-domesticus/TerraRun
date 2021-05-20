package internal

import (
	"encoding/json"
	"github.com/hashicorp/terraform-json"
)

func GetTerraformShowPlan() Command {
	cmd := NewTerraformCommand()
	cmd.Parameters = append(cmd.Parameters, []Parameter{
		&SimpleParameter{Value: "show"},
		&SimpleParameter{Value: "-json"},
		&SimpleParameter{Value: GetPlanFilename()},
	}...)
	return cmd
}

func getPlanJSON(config Config, stack TerraformStack) ([]byte, error) {
	showCmd := GetTerraformShowPlan()
	output, err := showCmd.Execute(config, stack)
	if err != nil {
		return nil, err
	}
	return output.StdOut, output.Error
}

func GetShowOutput(config Config, stack TerraformStack) (ShowOutput, error) {
	data, err := getPlanJSON(config, stack)
	if err != nil {
		return ShowOutput{}, err
	}
	output := tfjson.Plan{}
	err = json.Unmarshal(data, &output)
	return ShowOutput{
		Stack: stack,
		Plan:  output,
	}, err
}

type ShowOutput struct {
	Stack TerraformStack
	Plan  tfjson.Plan
}
