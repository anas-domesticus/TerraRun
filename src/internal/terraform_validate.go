package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

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
	fmt.Printf("Validating %s...\n", stack.Path)
	envVars := append(os.Environ(), fmt.Sprintf("TF_PLUGIN_CACHE_DIR=%s", config.TFPluginCacheDir))
	initCmd := GetTerraformInit()
	initCmd.Parameters = append(initCmd.Parameters, &SimpleParameter{Value: "-backend=false"})
	initCmd.EnvVars = envVars
	output, err := initCmd.Execute(config, stack)
	if err != nil {
		// TODO: Add detail to error
		return output, err
	}
	validateCmd := GetTerraformValidate()
	validateCmd.EnvVars = envVars
	output, err = validateCmd.Execute(config, stack)
	return output, err
}
