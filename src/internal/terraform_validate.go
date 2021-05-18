package internal

import "encoding/json"

func GetTerraformValidate() Command {
	cmd := NewTerraformCommand()
	cmd.Parameters = append(cmd.Parameters, []Parameter{
		&SimpleParameter{Value: "validate"},
		&SimpleParameter{Value: "-json"},
	}...)
	return cmd
}

type ValidateOutput struct {
	Valid        bool `json:"valid"`
	ErrorCount   int  `json:"error_count"`
	WarningCount int  `json:"warning_count"`
}

func ValidateWasSuccessful(output ExecuteOutput) bool {
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

func ValidateStack(config Config, stack TerraformStack) (ExecuteOutput, error) {
	initCmd := GetTerraformInit()
	validateCmd := GetTerraformValidate()
	output := initCmd.Execute(config, stack)
	if output.Error != nil {
		return ExecuteOutput{}, output.Error
	}
	output = validateCmd.Execute(config, stack)
	return output, nil
}
